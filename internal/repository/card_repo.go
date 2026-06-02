package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/ent"
	"github.com/tendo-mulira/tnotes-teams/internal/ent/card"
)

// CardRepository handles card persistence.
type CardRepository struct {
	client *ent.Client
	db     *sql.DB
}

// NewCardRepository creates a new CardRepository.
func NewCardRepository(client *ent.Client, db *sql.DB) *CardRepository {
	return &CardRepository{client: client, db: db}
}

// Create creates a new card.
func (r *CardRepository) Create(ctx context.Context, title, description string, listID, createdBy uuid.UUID, position int, dueDate *time.Time, labels []string) (*ent.Card, error) {
	builder := r.client.Card.Create().
		SetTitle(title).
		SetDescription(description).
		SetListID(listID).
		SetCreatedBy(createdBy).
		SetPosition(position)

	if dueDate != nil {
		builder.SetDueDate(*dueDate)
	}
	if labels != nil {
		builder.SetLabels(labels)
	}

	return builder.Save(ctx)
}

// GetByID returns a card by ID.
func (r *CardRepository) GetByID(ctx context.Context, id uuid.UUID) (*ent.Card, error) {
	return r.client.Card.Get(ctx, id)
}

// GetWithDetails loads a card with all associated data.
func (r *CardRepository) GetWithDetails(ctx context.Context, id uuid.UUID) (*ent.Card, error) {
	return r.client.Card.Query().
		Where(card.IDEQ(id)).
		WithSubtasks(func(q *ent.SubtaskQuery) {
			q.Order(ent.Asc("position"))
		}).
		WithQuickNotes(func(q *ent.QuickNoteQuery) {
			q.Order(ent.Desc("created_at"))
		}).
		WithComments(func(q *ent.CommentQuery) {
			q.WithUser().
				Order(ent.Asc("created_at"))
		}).
		WithCreator().
		WithList().
		Only(ctx)
}

// ListByList returns all cards for a list ordered by position.
func (r *CardRepository) ListByList(ctx context.Context, listID uuid.UUID) ([]*ent.Card, error) {
	return r.client.Card.Query().
		Where(card.ListIDEQ(listID)).
		Order(ent.Asc(card.FieldPosition)).
		All(ctx)
}

// Update updates a card's fields.
func (r *CardRepository) Update(ctx context.Context, id uuid.UUID, title, description *string, dueDate *time.Time, labels []string) (*ent.Card, error) {
	update := r.client.Card.UpdateOneID(id)
	if title != nil {
		update.SetTitle(*title)
	}
	if description != nil {
		update.SetDescription(*description)
	}
	if dueDate != nil {
		update.SetDueDate(*dueDate)
	}
	if labels != nil {
		update.SetLabels(labels)
	}
	return update.Save(ctx)
}

// Move moves a card to a new list and position.
func (r *CardRepository) Move(ctx context.Context, id uuid.UUID, newListID uuid.UUID, newPosition int) (*ent.Card, error) {
	return r.client.Card.UpdateOneID(id).
		SetListID(newListID).
		SetPosition(newPosition).
		Save(ctx)
}

// UpdateProgress updates a card's progress percentage.
func (r *CardRepository) UpdateProgress(ctx context.Context, id uuid.UUID, progress int) (*ent.Card, error) {
	return r.client.Card.UpdateOneID(id).
		SetProgressPercentage(progress).
		Save(ctx)
}

// Delete permanently deletes a card.
func (r *CardRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.Card.DeleteOneID(id).Exec(ctx)
}

// GetNextPosition returns the next available position for a new card in a list.
func (r *CardRepository) GetNextPosition(ctx context.Context, listID uuid.UUID) (int, error) {
	cards, err := r.client.Card.Query().
		Where(card.ListIDEQ(listID)).
		Order(ent.Desc(card.FieldPosition)).
		Limit(1).
		All(ctx)
	if err != nil {
		return 0, err
	}
	if len(cards) == 0 {
		return 0, nil
	}
	return cards[0].Position + 1, nil
}
