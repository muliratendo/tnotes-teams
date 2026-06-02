package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/repository"
)

// ExportService handles board export/import.
type ExportService struct {
	repos *repository.Repositories
}

// NewExportService creates a new ExportService.
func NewExportService(repos *repository.Repositories) *ExportService {
	return &ExportService{repos: repos}
}

// BoardExport is the full export schema for a board.
type BoardExport struct {
	Version     string           `json:"version"`
	ExportedAt  time.Time        `json:"exported_at"`
	Board       BoardExportData  `json:"board"`
}

// BoardExportData contains the board data for export.
type BoardExportData struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	ColorTheme  string           `json:"color_theme"`
	Lists       []ListExportData `json:"lists"`
}

// ListExportData contains list data for export.
type ListExportData struct {
	Title    string           `json:"title"`
	Position int              `json:"position"`
	Cards    []CardExportData `json:"cards"`
}

// CardExportData contains card data for export.
type CardExportData struct {
	Title              string               `json:"title"`
	Description        string               `json:"description"`
	Position           int                  `json:"position"`
	DueDate            *time.Time           `json:"due_date,omitempty"`
	Labels             []string             `json:"labels"`
	ProgressPercentage int                  `json:"progress_percentage"`
	Subtasks           []SubtaskExportData  `json:"subtasks"`
	QuickNotes         []NoteExportData     `json:"quick_notes"`
	Comments           []CommentExportData  `json:"comments"`
}

// SubtaskExportData contains subtask data for export.
type SubtaskExportData struct {
	Title       string `json:"title"`
	IsCompleted bool   `json:"is_completed"`
	Position    int    `json:"position"`
}

// NoteExportData contains quick note data for export.
type NoteExportData struct {
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// CommentExportData contains comment data for export.
type CommentExportData struct {
	Content   string    `json:"content"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// ExportBoard generates a JSON export of a board with all its data.
func (s *ExportService) ExportBoard(ctx context.Context, boardID uuid.UUID) (*BoardExport, error) {
	board, err := s.repos.Board.GetWithFullData(ctx, boardID)
	if err != nil {
		return nil, errors.New("board not found")
	}

	export := &BoardExport{
		Version:    "1.0",
		ExportedAt: time.Now(),
		Board: BoardExportData{
			Name:        board.Name,
			Description: board.Description,
			ColorTheme:  board.ColorTheme,
			Lists:       make([]ListExportData, 0),
		},
	}

	if board.Edges.Lists != nil {
		for _, l := range board.Edges.Lists {
			listExport := ListExportData{
				Title:    l.Title,
				Position: l.Position,
				Cards:    make([]CardExportData, 0),
			}

			if l.Edges.Cards != nil {
				for _, c := range l.Edges.Cards {
					cardExport := CardExportData{
						Title:              c.Title,
						Description:        c.Description,
						Position:           c.Position,
						DueDate:            c.DueDate,
						Labels:             c.Labels,
						ProgressPercentage: c.ProgressPercentage,
						Subtasks:           make([]SubtaskExportData, 0),
						QuickNotes:         make([]NoteExportData, 0),
						Comments:           make([]CommentExportData, 0),
					}

					for _, st := range c.Edges.Subtasks {
						cardExport.Subtasks = append(cardExport.Subtasks, SubtaskExportData{
							Title:       st.Title,
							IsCompleted: st.IsCompleted,
							Position:    st.Position,
						})
					}

					for _, n := range c.Edges.QuickNotes {
						cardExport.QuickNotes = append(cardExport.QuickNotes, NoteExportData{
							Content:   n.Content,
							CreatedAt: n.CreatedAt,
						})
					}

					for _, cm := range c.Edges.Comments {
						username := ""
						if u := cm.Edges.User; u != nil {
							username = u.Username
						}
						cardExport.Comments = append(cardExport.Comments, CommentExportData{
							Content:   cm.Content,
							Username:  username,
							CreatedAt: cm.CreatedAt,
						})
					}

					listExport.Cards = append(listExport.Cards, cardExport)
				}
			}

			export.Board.Lists = append(export.Board.Lists, listExport)
		}
	}

	return export, nil
}

// ImportBoard creates a board from a JSON export.
func (s *ExportService) ImportBoard(ctx context.Context, data []byte, workspaceID, userID uuid.UUID) (*BoardDTO, error) {
	var export BoardExport
	if err := json.Unmarshal(data, &export); err != nil {
		return nil, errors.New("invalid JSON format")
	}

	if export.Board.Name == "" {
		return nil, errors.New("board name is required in import data")
	}

	// Create board
	board, err := s.repos.Board.Create(ctx, export.Board.Name, export.Board.Description, workspaceID, userID)
	if err != nil {
		return nil, err
	}

	// Create lists and cards
	for _, listData := range export.Board.Lists {
		list, err := s.repos.List.Create(ctx, listData.Title, board.ID, listData.Position)
		if err != nil {
			continue
		}

		for _, cardData := range listData.Cards {
			card, err := s.repos.Card.Create(ctx, cardData.Title, cardData.Description, list.ID, userID, cardData.Position, cardData.DueDate, cardData.Labels)
			if err != nil {
				continue
			}

			for _, stData := range cardData.Subtasks {
				s.repos.Subtask.Create(ctx, stData.Title, card.ID, stData.Position)
				if stData.IsCompleted {
					// We'd need the subtask ID to toggle it
				}
			}

			for _, noteData := range cardData.QuickNotes {
				s.repos.QuickNote.Create(ctx, noteData.Content, card.ID, userID)
			}
		}
	}

	return toBoardDTO(board), nil
}

// ExportBoardCSV generates a CSV export of cards.
func (s *ExportService) ExportBoardCSV(ctx context.Context, boardID uuid.UUID) (string, error) {
	board, err := s.repos.Board.GetWithFullData(ctx, boardID)
	if err != nil {
		return "", errors.New("board not found")
	}

	csv := "List,Card Title,Description,Due Date,Progress,Labels\n"

	if board.Edges.Lists != nil {
		for _, l := range board.Edges.Lists {
			if l.Edges.Cards != nil {
				for _, c := range l.Edges.Cards {
					dueDate := ""
					if c.DueDate != nil {
						dueDate = c.DueDate.Format("2006-01-02")
					}
					labels := ""
					if len(c.Labels) > 0 {
						for i, label := range c.Labels {
							if i > 0 {
								labels += ";"
							}
							labels += label
						}
					}
					csv += `"` + l.Title + `","` + c.Title + `","` + c.Description + `","` + dueDate + `",` + strconv.Itoa(c.ProgressPercentage) + `%,"` + labels + "\"\n"
				}
			}
		}
	}

	return csv, nil
}
