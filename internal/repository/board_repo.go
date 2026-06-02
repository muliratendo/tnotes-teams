package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/ent"
	"github.com/tendo-mulira/tnotes-teams/internal/ent/board"
	entlist "github.com/tendo-mulira/tnotes-teams/internal/ent/list"
)

// BoardRepository handles board persistence.
type BoardRepository struct {
	client *ent.Client
	db     *sql.DB
}

// NewBoardRepository creates a new BoardRepository.
func NewBoardRepository(client *ent.Client, db *sql.DB) *BoardRepository {
	return &BoardRepository{client: client, db: db}
}

// Create creates a new board.
func (r *BoardRepository) Create(ctx context.Context, name, description string, workspaceID, createdBy uuid.UUID) (*ent.Board, error) {
	return r.client.Board.Create().
		SetName(name).
		SetDescription(description).
		SetWorkspaceID(workspaceID).
		SetCreatedBy(createdBy).
		Save(ctx)
}

// GetByID returns a board by ID.
func (r *BoardRepository) GetByID(ctx context.Context, id uuid.UUID) (*ent.Board, error) {
	return r.client.Board.Get(ctx, id)
}

// GetWithFullData loads a board with all lists, cards, subtasks, notes, and comments.
func (r *BoardRepository) GetWithFullData(ctx context.Context, id uuid.UUID) (*ent.Board, error) {
	return r.client.Board.Query().
		Where(board.IDEQ(id)).
		WithLists(func(q *ent.ListQuery) {
			q.Order(ent.Asc(entlist.FieldPosition)).
				WithCards(func(cq *ent.CardQuery) {
					cq.Order(ent.Asc("position")).
						WithSubtasks(func(sq *ent.SubtaskQuery) {
							sq.Order(ent.Asc("position"))
						}).
						WithQuickNotes().
						WithComments(func(cmq *ent.CommentQuery) {
							cmq.WithUser().
								Order(ent.Desc("created_at"))
						}).
						WithCreator()
				})
		}).
		Only(ctx)
}

// ListByWorkspace returns all boards for a workspace.
func (r *BoardRepository) ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]*ent.Board, error) {
	return r.client.Board.Query().
		Where(
			board.WorkspaceIDEQ(workspaceID),
			board.IsArchivedEQ(false),
		).
		Order(ent.Desc(board.FieldCreatedAt)).
		All(ctx)
}

// Update updates a board.
func (r *BoardRepository) Update(ctx context.Context, id uuid.UUID, name, description, colorTheme *string) (*ent.Board, error) {
	update := r.client.Board.UpdateOneID(id)
	if name != nil {
		update.SetName(*name)
	}
	if description != nil {
		update.SetDescription(*description)
	}
	if colorTheme != nil {
		update.SetColorTheme(*colorTheme)
	}
	return update.Save(ctx)
}

// Delete permanently deletes a board and all its data (cascading).
func (r *BoardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.Board.DeleteOneID(id).Exec(ctx)
}

// Archive soft-deletes a board by setting is_archived to true.
func (r *BoardRepository) Archive(ctx context.Context, id uuid.UUID) (*ent.Board, error) {
	return r.client.Board.UpdateOneID(id).
		SetIsArchived(true).
		Save(ctx)
}

// Unarchive restores an archived board.
func (r *BoardRepository) Unarchive(ctx context.Context, id uuid.UUID) (*ent.Board, error) {
	return r.client.Board.UpdateOneID(id).
		SetIsArchived(false).
		Save(ctx)
}
