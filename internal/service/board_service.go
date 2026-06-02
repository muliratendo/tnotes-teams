package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/ent"
	"github.com/tendo-mulira/tnotes-teams/internal/repository"
)

// BoardService handles board business logic.
type BoardService struct {
	repos *repository.Repositories
}

// NewBoardService creates a new BoardService.
func NewBoardService(repos *repository.Repositories) *BoardService {
	return &BoardService{repos: repos}
}

// BoardDTO is the public board data.
type BoardDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ColorTheme  string    `json:"color_theme"`
	IsArchived  bool      `json:"is_archived"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	Lists       []ListDTO `json:"lists,omitempty"`
}

// ListDTO is the public list data.
type ListDTO struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	Position int       `json:"position"`
	BoardID  uuid.UUID `json:"board_id"`
	Cards    []CardDTO `json:"cards,omitempty"`
}

// CardDTO is the public card data.
type CardDTO struct {
	ID                 uuid.UUID    `json:"id"`
	Title              string       `json:"title"`
	Description        string       `json:"description"`
	Position           int          `json:"position"`
	DueDate            *time.Time   `json:"due_date,omitempty"`
	Labels             []string     `json:"labels"`
	ProgressPercentage int          `json:"progress_percentage"`
	ListID             uuid.UUID    `json:"list_id"`
	CreatedBy          *uuid.UUID   `json:"created_by,omitempty"`
	CreatedAt          time.Time    `json:"created_at"`
	Subtasks           []SubtaskDTO `json:"subtasks,omitempty"`
	QuickNotes         []NoteDTO    `json:"quick_notes,omitempty"`
	Comments           []CommentDTO `json:"comments,omitempty"`
	CreatorName        string       `json:"creator_name,omitempty"`
}

// SubtaskDTO is the public subtask data.
type SubtaskDTO struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	IsCompleted bool      `json:"is_completed"`
	Position    int       `json:"position"`
	CardID      uuid.UUID `json:"card_id"`
}

// NoteDTO is the public quick note data.
type NoteDTO struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	CardID    uuid.UUID `json:"card_id"`
	CreatedAt time.Time `json:"created_at"`
}

// CommentDTO is the public comment data.
type CommentDTO struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	CardID    uuid.UUID `json:"card_id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username,omitempty"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateBoardInput is the input for creating a board.
type CreateBoardInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Create creates a new board.
func (s *BoardService) Create(ctx context.Context, workspaceID uuid.UUID, input CreateBoardInput, userID uuid.UUID) (*BoardDTO, error) {
	if input.Name == "" {
		return nil, errors.New("board name is required")
	}

	b, err := s.repos.Board.Create(ctx, input.Name, input.Description, workspaceID, userID)
	if err != nil {
		return nil, err
	}

	// Seed default Kanban columns
	defaultLists := []string{"To Do", "In Progress", "Done"}
	for i, title := range defaultLists {
		_, _ = s.repos.List.Create(ctx, title, b.ID, i)
	}

	s.repos.ActivityLog.Log(ctx, b.ID, userID, "board_created", "board", &b.ID, nil, map[string]interface{}{
		"name": b.Name,
	})

	return toBoardDTO(b), nil
}

// GetWithFullData loads a board with all nested data.
func (s *BoardService) GetWithFullData(ctx context.Context, id uuid.UUID) (*BoardDTO, error) {
	b, err := s.repos.Board.GetWithFullData(ctx, id)
	if err != nil {
		return nil, err
	}
	return toBoardDTOWithEdges(b), nil
}

// ListByWorkspace returns all boards for a workspace.
func (s *BoardService) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]BoardDTO, error) {
	boards, err := s.repos.Board.ListByWorkspace(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	dtos := make([]BoardDTO, 0, len(boards))
	for _, b := range boards {
		dtos = append(dtos, *toBoardDTO(b))
	}
	return dtos, nil
}

// Update updates a board.
func (s *BoardService) Update(ctx context.Context, id uuid.UUID, name, description, colorTheme *string, userID uuid.UUID) (*BoardDTO, error) {
	b, err := s.repos.Board.Update(ctx, id, name, description, colorTheme)
	if err != nil {
		return nil, err
	}

	s.repos.ActivityLog.Log(ctx, b.ID, userID, "board_updated", "board", &b.ID, nil, map[string]interface{}{
		"name": b.Name,
	})

	return toBoardDTO(b), nil
}

// CreateListInput is the input for creating a list column.
type CreateListInput struct {
	Title string `json:"title"`
}

// CreateList adds a new list column to a board.
func (s *BoardService) CreateList(ctx context.Context, boardID uuid.UUID, input CreateListInput, userID uuid.UUID) (*ListDTO, error) {
	if input.Title == "" {
		return nil, errors.New("list title is required")
	}

	if _, err := s.repos.Board.GetByID(ctx, boardID); err != nil {
		return nil, errors.New("board not found")
	}

	position, err := s.repos.List.GetNextPosition(ctx, boardID)
	if err != nil {
		return nil, err
	}

	l, err := s.repos.List.Create(ctx, input.Title, boardID, position)
	if err != nil {
		return nil, err
	}

	s.repos.ActivityLog.Log(ctx, boardID, userID, "list_created", "list", &l.ID, nil, map[string]interface{}{
		"title": l.Title,
	})

	return &ListDTO{
		ID:       l.ID,
		Title:    l.Title,
		Position: l.Position,
		BoardID:  l.BoardID,
	}, nil
}

// Delete deletes a board (admin only).
func (s *BoardService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	s.repos.ActivityLog.Log(ctx, id, userID, "board_deleted", "board", &id, nil, nil)
	return s.repos.Board.Delete(ctx, id)
}

// Archive archives a board.
func (s *BoardService) Archive(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*BoardDTO, error) {
	b, err := s.repos.Board.Archive(ctx, id)
	if err != nil {
		return nil, err
	}
	return toBoardDTO(b), nil
}

// CalculateProgress computes progress from subtask completion.
func CalculateProgress(subtasks []*ent.Subtask) int {
	if len(subtasks) == 0 {
		return 0
	}
	completed := 0
	for _, st := range subtasks {
		if st.IsCompleted {
			completed++
		}
	}
	return (completed * 100) / len(subtasks)
}

func toBoardDTO(b *ent.Board) *BoardDTO {
	return &BoardDTO{
		ID:          b.ID,
		Name:        b.Name,
		Description: b.Description,
		ColorTheme:  b.ColorTheme,
		IsArchived:  b.IsArchived,
		WorkspaceID: b.WorkspaceID,
		CreatedBy:   b.CreatedBy,
		CreatedAt:   b.CreatedAt,
	}
}

func toBoardDTOWithEdges(b *ent.Board) *BoardDTO {
	dto := toBoardDTO(b)
	dto.Lists = make([]ListDTO, 0)

	if b.Edges.Lists != nil {
		for _, l := range b.Edges.Lists {
			listDTO := ListDTO{
				ID:       l.ID,
				Title:    l.Title,
				Position: l.Position,
				BoardID:  l.BoardID,
				Cards:    make([]CardDTO, 0),
			}
			if l.Edges.Cards != nil {
				for _, c := range l.Edges.Cards {
					cardDTO := toCardDTO(c)
					listDTO.Cards = append(listDTO.Cards, *cardDTO)
				}
			}
			dto.Lists = append(dto.Lists, listDTO)
		}
	}
	return dto
}

func toCardDTO(c *ent.Card) *CardDTO {
	dto := &CardDTO{
		ID:                 c.ID,
		Title:              c.Title,
		Description:        c.Description,
		Position:           c.Position,
		DueDate:            c.DueDate,
		Labels:             c.Labels,
		ProgressPercentage: c.ProgressPercentage,
		ListID:             c.ListID,
		CreatedAt:          c.CreatedAt,
		Subtasks:           make([]SubtaskDTO, 0),
		QuickNotes:         make([]NoteDTO, 0),
		Comments:           make([]CommentDTO, 0),
	}

	if c.CreatedBy != uuid.Nil {
		dto.CreatedBy = &c.CreatedBy
	}

	if creator := c.Edges.Creator; creator != nil {
		dto.CreatorName = creator.Username
	}

	if c.Edges.Subtasks != nil {
		for _, st := range c.Edges.Subtasks {
			dto.Subtasks = append(dto.Subtasks, SubtaskDTO{
				ID:          st.ID,
				Title:       st.Title,
				IsCompleted: st.IsCompleted,
				Position:    st.Position,
				CardID:      st.CardID,
			})
		}
	}

	if c.Edges.QuickNotes != nil {
		for _, n := range c.Edges.QuickNotes {
			dto.QuickNotes = append(dto.QuickNotes, NoteDTO{
				ID:        n.ID,
				Content:   n.Content,
				CardID:    n.CardID,
				CreatedAt: n.CreatedAt,
			})
		}
	}

	if c.Edges.Comments != nil {
		for _, cm := range c.Edges.Comments {
			commentDTO := CommentDTO{
				ID:        cm.ID,
				Content:   cm.Content,
				CardID:    cm.CardID,
				UserID:    cm.UserID,
				CreatedAt: cm.CreatedAt,
			}
			if u := cm.Edges.User; u != nil {
				commentDTO.Username = u.Username
				commentDTO.AvatarURL = u.AvatarURL
			}
			dto.Comments = append(dto.Comments, commentDTO)
		}
	}

	return dto
}
