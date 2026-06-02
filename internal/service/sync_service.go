package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/repository"
	"github.com/tendo-mulira/tnotes-teams/internal/websocket"
)

// SyncService handles offline sync processing.
type SyncService struct {
	repos *repository.Repositories
	hub   *websocket.Hub
}

// NewSyncService creates a new SyncService.
func NewSyncService(repos *repository.Repositories, hub *websocket.Hub) *SyncService {
	return &SyncService{repos: repos, hub: hub}
}

// SyncOperation represents a single offline mutation to be synced.
type SyncOperation struct {
	ID         string                 `json:"id"`          // Client-side idempotency key
	Action     string                 `json:"action"`      // e.g., "card_created", "card_moved", "subtask_toggled"
	EntityType string                 `json:"entity_type"` // e.g., "card", "subtask", "comment"
	EntityID   string                 `json:"entity_id"`   // UUID of the entity
	Data       map[string]interface{} `json:"data"`        // Mutation payload
	Timestamp  time.Time              `json:"timestamp"`   // Client-side timestamp
	BoardID    string                 `json:"board_id"`    // Board context
}

// SyncBatchInput is the input for processing an offline batch.
type SyncBatchInput struct {
	Operations []SyncOperation `json:"operations"`
}

// SyncResult represents the result of syncing a single operation.
type SyncResult struct {
	OperationID string      `json:"operation_id"`
	Status      string      `json:"status"` // "synced", "conflict", "error"
	Message     string      `json:"message,omitempty"`
	ServerData  interface{} `json:"server_data,omitempty"`
}

// SyncBatchResponse is the response after processing a batch.
type SyncBatchResponse struct {
	SyncedCount   int          `json:"synced_count"`
	ConflictCount int          `json:"conflict_count"`
	ErrorCount    int          `json:"error_count"`
	Results       []SyncResult `json:"results"`
}

// ProcessBatch processes a batch of offline mutations.
// Uses last-write-wins with server timestamps for conflict resolution.
func (s *SyncService) ProcessBatch(ctx context.Context, input SyncBatchInput, userID uuid.UUID) (*SyncBatchResponse, error) {
	response := &SyncBatchResponse{
		Results: make([]SyncResult, 0, len(input.Operations)),
	}

	for _, op := range input.Operations {
		result := s.processOperation(ctx, op, userID)
		response.Results = append(response.Results, result)

		switch result.Status {
		case "synced":
			response.SyncedCount++
		case "conflict":
			response.ConflictCount++
		case "error":
			response.ErrorCount++
		}
	}

	return response, nil
}

// processOperation processes a single sync operation.
func (s *SyncService) processOperation(ctx context.Context, op SyncOperation, userID uuid.UUID) SyncResult {
	result := SyncResult{
		OperationID: op.ID,
	}

	switch op.Action {
	case "card_moved":
		result = s.syncCardMove(ctx, op, userID)
	case "card_updated":
		result = s.syncCardUpdate(ctx, op, userID)
	case "subtask_toggled":
		result = s.syncSubtaskToggle(ctx, op, userID)
	case "subtask_created":
		result = s.syncSubtaskCreate(ctx, op, userID)
	case "comment_added":
		result = s.syncCommentAdd(ctx, op, userID)
	case "note_created":
		result = s.syncNoteCreate(ctx, op, userID)
	case "card_created":
		result = s.syncCardCreate(ctx, op, userID)
	case "card_deleted":
		result = s.syncCardDelete(ctx, op, userID)
	default:
		result.Status = "error"
		result.Message = "unknown action: " + op.Action
	}

	result.OperationID = op.ID
	return result
}

func (s *SyncService) syncCardMove(ctx context.Context, op SyncOperation, userID uuid.UUID) SyncResult {
	cardID, err := uuid.Parse(op.EntityID)
	if err != nil {
		return SyncResult{Status: "error", Message: "invalid card ID"}
	}

	// Check if card still exists and get its current state
	card, err := s.repos.Card.GetByID(ctx, cardID)
	if err != nil {
		return SyncResult{Status: "error", Message: "card not found"}
	}

	// Last-write-wins: if the card was updated after this offline mutation, skip
	if card.UpdatedAt.After(op.Timestamp) {
		return SyncResult{
			Status:     "conflict",
			Message:    "card was modified after your offline change - server version kept",
			ServerData: card,
		}
	}

	newListIDStr, _ := op.Data["new_list_id"].(string)
	newListID, err := uuid.Parse(newListIDStr)
	if err != nil {
		return SyncResult{Status: "error", Message: "invalid new list ID"}
	}

	newPos, _ := op.Data["new_position"].(float64)
	_, err = s.repos.Card.Move(ctx, cardID, newListID, int(newPos))
	if err != nil {
		return SyncResult{Status: "error", Message: err.Error()}
	}

	return SyncResult{Status: "synced"}
}

func (s *SyncService) syncCardUpdate(ctx context.Context, op SyncOperation, userID uuid.UUID) SyncResult {
	cardID, err := uuid.Parse(op.EntityID)
	if err != nil {
		return SyncResult{Status: "error", Message: "invalid card ID"}
	}

	card, err := s.repos.Card.GetByID(ctx, cardID)
	if err != nil {
		return SyncResult{Status: "error", Message: "card not found"}
	}

	if card.UpdatedAt.After(op.Timestamp) {
		return SyncResult{
			Status:  "conflict",
			Message: "card was modified after your offline change - server version kept",
		}
	}

	var title, description *string
	if t, ok := op.Data["title"].(string); ok {
		title = &t
	}
	if d, ok := op.Data["description"].(string); ok {
		description = &d
	}

	_, err = s.repos.Card.Update(ctx, cardID, title, description, nil, nil)
	if err != nil {
		return SyncResult{Status: "error", Message: err.Error()}
	}

	return SyncResult{Status: "synced"}
}

func (s *SyncService) syncSubtaskToggle(ctx context.Context, op SyncOperation, userID uuid.UUID) SyncResult {
	subtaskID, err := uuid.Parse(op.EntityID)
	if err != nil {
		return SyncResult{Status: "error", Message: "invalid subtask ID"}
	}

	completed, _ := op.Data["is_completed"].(bool)
	_, err = s.repos.Subtask.SetCompleted(ctx, subtaskID, completed)
	if err != nil {
		return SyncResult{Status: "error", Message: err.Error()}
	}

	// Recalculate progress
	st, _ := s.repos.Subtask.GetByID(ctx, subtaskID)
	if st != nil {
		total, comp, _ := s.repos.Subtask.CountByCard(ctx, st.CardID)
		if total > 0 {
			progress := (comp * 100) / total
			s.repos.Card.UpdateProgress(ctx, st.CardID, progress)
		}
	}

	return SyncResult{Status: "synced"}
}

func (s *SyncService) syncCommentAdd(ctx context.Context, op SyncOperation, userID uuid.UUID) SyncResult {
	cardID, err := uuid.Parse(op.EntityID)
	if err != nil {
		return SyncResult{Status: "error", Message: "invalid card ID"}
	}

	content, _ := op.Data["content"].(string)
	if content == "" {
		return SyncResult{Status: "error", Message: "comment content required"}
	}

	_, err = s.repos.Comment.Create(ctx, content, cardID, userID)
	if err != nil {
		return SyncResult{Status: "error", Message: err.Error()}
	}

	return SyncResult{Status: "synced"}
}

func (s *SyncService) syncCardCreate(ctx context.Context, op SyncOperation, userID uuid.UUID) SyncResult {
	listIDStr, _ := op.Data["list_id"].(string)
	listID, err := uuid.Parse(listIDStr)
	if err != nil {
		return SyncResult{Status: "error", Message: "invalid list ID"}
	}

	title, _ := op.Data["title"].(string)
	description, _ := op.Data["description"].(string)

	pos, err := s.repos.Card.GetNextPosition(ctx, listID)
	if err != nil {
		return SyncResult{Status: "error", Message: err.Error()}
	}

	_, err = s.repos.Card.Create(ctx, title, description, listID, userID, pos, nil, nil)
	if err != nil {
		return SyncResult{Status: "error", Message: err.Error()}
	}

	return SyncResult{Status: "synced"}
}

func (s *SyncService) syncSubtaskCreate(ctx context.Context, op SyncOperation, userID uuid.UUID) SyncResult {
	cardIDStr, _ := op.Data["card_id"].(string)
	cardID, err := uuid.Parse(cardIDStr)
	if err != nil {
		return SyncResult{Status: "error", Message: "invalid card ID"}
	}

	title, _ := op.Data["title"].(string)
	pos, err := s.repos.Subtask.GetNextPosition(ctx, cardID)
	if err != nil {
		return SyncResult{Status: "error", Message: err.Error()}
	}

	_, err = s.repos.Subtask.Create(ctx, title, cardID, pos)
	if err != nil {
		return SyncResult{Status: "error", Message: err.Error()}
	}

	total, comp, _ := s.repos.Subtask.CountByCard(ctx, cardID)
	if total > 0 {
		s.repos.Card.UpdateProgress(ctx, cardID, (comp*100)/total)
	}

	return SyncResult{Status: "synced"}
}

func (s *SyncService) syncNoteCreate(ctx context.Context, op SyncOperation, userID uuid.UUID) SyncResult {
	cardIDStr, _ := op.Data["card_id"].(string)
	cardID, err := uuid.Parse(cardIDStr)
	if err != nil {
		return SyncResult{Status: "error", Message: "invalid card ID"}
	}

	content, _ := op.Data["content"].(string)
	if content == "" {
		return SyncResult{Status: "error", Message: "note content required"}
	}

	_, err = s.repos.QuickNote.Create(ctx, content, cardID, userID)
	if err != nil {
		return SyncResult{Status: "error", Message: err.Error()}
	}

	return SyncResult{Status: "synced"}
}

func (s *SyncService) syncCardDelete(ctx context.Context, op SyncOperation, userID uuid.UUID) SyncResult {
	cardID, err := uuid.Parse(op.EntityID)
	if err != nil {
		return SyncResult{Status: "error", Message: "invalid card ID"}
	}

	if err := s.repos.Card.Delete(ctx, cardID); err != nil {
		return SyncResult{Status: "error", Message: err.Error()}
	}

	return SyncResult{Status: "synced"}
}
