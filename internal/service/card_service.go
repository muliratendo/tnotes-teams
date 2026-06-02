package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/repository"
	"github.com/tendo-mulira/tnotes-teams/internal/websocket"
)

// CardService handles card business logic.
type CardService struct {
	repos *repository.Repositories
	hub   *websocket.Hub
}

// NewCardService creates a new CardService.
func NewCardService(repos *repository.Repositories, hub *websocket.Hub) *CardService {
	return &CardService{repos: repos, hub: hub}
}

// CreateCardInput is the input for creating a card.
type CreateCardInput struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Labels      []string   `json:"labels,omitempty"`
}

// UpdateCardInput is the input for updating a card.
type UpdateCardInput struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Labels      []string   `json:"labels,omitempty"`
}

// MoveCardInput is the input for moving a card.
type MoveCardInput struct {
	NewListID   uuid.UUID `json:"new_list_id"`
	NewPosition int       `json:"new_position"`
}

// CreateCard creates a new card in a list.
func (s *CardService) CreateCard(ctx context.Context, listID uuid.UUID, input CreateCardInput, userID uuid.UUID) (*CardDTO, error) {
	if input.Title == "" {
		return nil, errors.New("card title is required")
	}

	// Get list to find the board for activity logging
	list, err := s.repos.List.GetByID(ctx, listID)
	if err != nil {
		return nil, errors.New("list not found")
	}

	// Get next position
	pos, err := s.repos.Card.GetNextPosition(ctx, listID)
	if err != nil {
		return nil, err
	}

	card, err := s.repos.Card.Create(ctx, input.Title, input.Description, listID, userID, pos, input.DueDate, input.Labels)
	if err != nil {
		return nil, err
	}

	// Log activity
	s.repos.ActivityLog.Log(ctx, list.BoardID, userID, "card_created", "card", &card.ID, nil, map[string]interface{}{
		"title":   card.Title,
		"list_id": listID.String(),
	})

	return toCardDTO(card), nil
}

// GetCardWithDetails returns a card with all associated data.
func (s *CardService) GetCardWithDetails(ctx context.Context, id uuid.UUID) (*CardDTO, error) {
	card, err := s.repos.Card.GetWithDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	return toCardDTO(card), nil
}

// UpdateCard updates a card's fields.
func (s *CardService) UpdateCard(ctx context.Context, id uuid.UUID, input UpdateCardInput, userID uuid.UUID) (*CardDTO, error) {
	card, err := s.repos.Card.Update(ctx, id, input.Title, input.Description, input.DueDate, input.Labels)
	if err != nil {
		return nil, err
	}

	// Get board ID for activity log
	fullCard, _ := s.repos.Card.GetWithDetails(ctx, id)
	if fullCard != nil && fullCard.Edges.List != nil {
		s.repos.ActivityLog.Log(ctx, fullCard.Edges.List.BoardID, userID, "card_updated", "card", &id, nil, nil)
	}

	return toCardDTO(card), nil
}

// MoveCard moves a card to a new list and/or position.
func (s *CardService) MoveCard(ctx context.Context, id uuid.UUID, input MoveCardInput, userID uuid.UUID) (*CardDTO, error) {
	card, err := s.repos.Card.Move(ctx, id, input.NewListID, input.NewPosition)
	if err != nil {
		return nil, err
	}

	// Broadcast card move via WebSocket
	// (will be handled by the WebSocket event system)

	return toCardDTO(card), nil
}

// DeleteCard deletes a card.
func (s *CardService) DeleteCard(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	return s.repos.Card.Delete(ctx, id)
}

// CreateSubtask adds a subtask to a card.
type CreateSubtaskInput struct {
	Title string `json:"title"`
}

func (s *CardService) CreateSubtask(ctx context.Context, cardID uuid.UUID, input CreateSubtaskInput, userID uuid.UUID) (*SubtaskDTO, error) {
	if input.Title == "" {
		return nil, errors.New("subtask title is required")
	}

	pos, err := s.repos.Subtask.GetNextPosition(ctx, cardID)
	if err != nil {
		return nil, err
	}

	st, err := s.repos.Subtask.Create(ctx, input.Title, cardID, pos)
	if err != nil {
		return nil, err
	}

	// Recalculate card progress
	s.recalculateProgress(ctx, cardID)

	return &SubtaskDTO{
		ID:          st.ID,
		Title:       st.Title,
		IsCompleted: st.IsCompleted,
		Position:    st.Position,
		CardID:      st.CardID,
	}, nil
}

// ToggleSubtask toggles a subtask's completion status.
func (s *CardService) ToggleSubtask(ctx context.Context, subtaskID uuid.UUID, userID uuid.UUID) (*SubtaskDTO, error) {
	st, err := s.repos.Subtask.Toggle(ctx, subtaskID)
	if err != nil {
		return nil, err
	}

	// Recalculate card progress
	s.recalculateProgress(ctx, st.CardID)

	return &SubtaskDTO{
		ID:          st.ID,
		Title:       st.Title,
		IsCompleted: st.IsCompleted,
		Position:    st.Position,
		CardID:      st.CardID,
	}, nil
}

// CreateQuickNote adds a quick note to a card.
type CreateNoteInput struct {
	Content string `json:"content"`
}

func (s *CardService) CreateQuickNote(ctx context.Context, cardID uuid.UUID, input CreateNoteInput, userID uuid.UUID) (*NoteDTO, error) {
	if input.Content == "" {
		return nil, errors.New("note content is required")
	}

	note, err := s.repos.QuickNote.Create(ctx, input.Content, cardID, userID)
	if err != nil {
		return nil, err
	}

	return &NoteDTO{
		ID:        note.ID,
		Content:   note.Content,
		CardID:    note.CardID,
		CreatedAt: note.CreatedAt,
	}, nil
}

// CreateComment adds a comment to a card.
type CreateCommentInput struct {
	Content string `json:"content"`
}

func (s *CardService) CreateComment(ctx context.Context, cardID uuid.UUID, input CreateCommentInput, userID uuid.UUID) (*CommentDTO, error) {
	if input.Content == "" {
		return nil, errors.New("comment content is required")
	}

	comment, err := s.repos.Comment.Create(ctx, input.Content, cardID, userID)
	if err != nil {
		return nil, err
	}

	// Get user info for response
	user, _ := s.repos.User.GetByID(ctx, userID)

	dto := &CommentDTO{
		ID:        comment.ID,
		Content:   comment.Content,
		CardID:    comment.CardID,
		UserID:    comment.UserID,
		CreatedAt: comment.CreatedAt,
	}

	if user != nil {
		dto.Username = user.Username
		dto.AvatarURL = user.AvatarURL
	}

	return dto, nil
}

// recalculateProgress recalculates and updates a card's progress from its subtasks.
func (s *CardService) recalculateProgress(ctx context.Context, cardID uuid.UUID) {
	total, completed, err := s.repos.Subtask.CountByCard(ctx, cardID)
	if err != nil {
		return
	}

	progress := 0
	if total > 0 {
		progress = (completed * 100) / total
	}

	s.repos.Card.UpdateProgress(ctx, cardID, progress)
}
