package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/ent"
	"github.com/tendo-mulira/tnotes-teams/internal/ent/subtask"
)

// SubtaskRepository handles subtask persistence.
type SubtaskRepository struct {
	client *ent.Client
	db     *sql.DB
}

// NewSubtaskRepository creates a new SubtaskRepository.
func NewSubtaskRepository(client *ent.Client, db *sql.DB) *SubtaskRepository {
	return &SubtaskRepository{client: client, db: db}
}

// Create creates a new subtask.
func (r *SubtaskRepository) Create(ctx context.Context, title string, cardID uuid.UUID, position int) (*ent.Subtask, error) {
	return r.client.Subtask.Create().
		SetTitle(title).
		SetCardID(cardID).
		SetPosition(position).
		Save(ctx)
}

// GetByID returns a subtask by ID.
func (r *SubtaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*ent.Subtask, error) {
	return r.client.Subtask.Get(ctx, id)
}

// Toggle toggles the completion status of a subtask.
func (r *SubtaskRepository) Toggle(ctx context.Context, id uuid.UUID) (*ent.Subtask, error) {
	s, err := r.client.Subtask.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.client.Subtask.UpdateOneID(id).
		SetIsCompleted(!s.IsCompleted).
		Save(ctx)
}

// SetCompleted sets the completion status of a subtask.
func (r *SubtaskRepository) SetCompleted(ctx context.Context, id uuid.UUID, completed bool) (*ent.Subtask, error) {
	return r.client.Subtask.UpdateOneID(id).
		SetIsCompleted(completed).
		Save(ctx)
}

// Delete permanently deletes a subtask.
func (r *SubtaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.Subtask.DeleteOneID(id).Exec(ctx)
}

// ListByCard returns all subtasks for a card ordered by position.
func (r *SubtaskRepository) ListByCard(ctx context.Context, cardID uuid.UUID) ([]*ent.Subtask, error) {
	return r.client.Subtask.Query().
		Where(subtask.CardIDEQ(cardID)).
		Order(ent.Asc(subtask.FieldPosition)).
		All(ctx)
}

// GetNextPosition returns the next position for a subtask.
func (r *SubtaskRepository) GetNextPosition(ctx context.Context, cardID uuid.UUID) (int, error) {
	subtasks, err := r.client.Subtask.Query().
		Where(subtask.CardIDEQ(cardID)).
		Order(ent.Desc(subtask.FieldPosition)).
		Limit(1).
		All(ctx)
	if err != nil {
		return 0, err
	}
	if len(subtasks) == 0 {
		return 0, nil
	}
	return subtasks[0].Position + 1, nil
}

// CountByCard returns total and completed counts for a card's subtasks.
func (r *SubtaskRepository) CountByCard(ctx context.Context, cardID uuid.UUID) (total int, completed int, err error) {
	subtasks, err := r.client.Subtask.Query().
		Where(subtask.CardIDEQ(cardID)).
		All(ctx)
	if err != nil {
		return 0, 0, err
	}
	total = len(subtasks)
	for _, s := range subtasks {
		if s.IsCompleted {
			completed++
		}
	}
	return total, completed, nil
}
