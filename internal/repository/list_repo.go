package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/ent"
	entlist "github.com/tendo-mulira/tnotes-teams/internal/ent/list"
)

// ListRepository handles list persistence.
type ListRepository struct {
	client *ent.Client
	db     *sql.DB
}

// NewListRepository creates a new ListRepository.
func NewListRepository(client *ent.Client, db *sql.DB) *ListRepository {
	return &ListRepository{client: client, db: db}
}

// Create creates a new list.
func (r *ListRepository) Create(ctx context.Context, title string, boardID uuid.UUID, position int) (*ent.List, error) {
	return r.client.List.Create().
		SetTitle(title).
		SetBoardID(boardID).
		SetPosition(position).
		Save(ctx)
}

// GetByID returns a list by ID.
func (r *ListRepository) GetByID(ctx context.Context, id uuid.UUID) (*ent.List, error) {
	return r.client.List.Get(ctx, id)
}

// ListByBoard returns all lists for a board ordered by position.
func (r *ListRepository) ListByBoard(ctx context.Context, boardID uuid.UUID) ([]*ent.List, error) {
	return r.client.List.Query().
		Where(entlist.BoardIDEQ(boardID)).
		Order(ent.Asc(entlist.FieldPosition)).
		All(ctx)
}

// Update updates a list's title.
func (r *ListRepository) Update(ctx context.Context, id uuid.UUID, title string) (*ent.List, error) {
	return r.client.List.UpdateOneID(id).
		SetTitle(title).
		Save(ctx)
}

// UpdatePosition updates a list's position.
func (r *ListRepository) UpdatePosition(ctx context.Context, id uuid.UUID, position int) (*ent.List, error) {
	return r.client.List.UpdateOneID(id).
		SetPosition(position).
		Save(ctx)
}

// Delete permanently deletes a list.
func (r *ListRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.List.DeleteOneID(id).Exec(ctx)
}

// GetNextPosition returns the next available position for a new list in a board.
func (r *ListRepository) GetNextPosition(ctx context.Context, boardID uuid.UUID) (int, error) {
	lists, err := r.client.List.Query().
		Where(entlist.BoardIDEQ(boardID)).
		Order(ent.Desc(entlist.FieldPosition)).
		Limit(1).
		All(ctx)
	if err != nil {
		return 0, err
	}
	if len(lists) == 0 {
		return 0, nil
	}
	return lists[0].Position + 1, nil
}
