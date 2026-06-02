package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/ent"
	"github.com/tendo-mulira/tnotes-teams/internal/ent/quicknote"
)

// QuickNoteRepository handles quick note persistence.
type QuickNoteRepository struct {
	client *ent.Client
	db     *sql.DB
}

// NewQuickNoteRepository creates a new QuickNoteRepository.
func NewQuickNoteRepository(client *ent.Client, db *sql.DB) *QuickNoteRepository {
	return &QuickNoteRepository{client: client, db: db}
}

// Create creates a new quick note.
func (r *QuickNoteRepository) Create(ctx context.Context, content string, cardID, createdBy uuid.UUID) (*ent.QuickNote, error) {
	return r.client.QuickNote.Create().
		SetContent(content).
		SetCardID(cardID).
		SetCreatedBy(createdBy).
		Save(ctx)
}

// Update updates a quick note's content.
func (r *QuickNoteRepository) Update(ctx context.Context, id uuid.UUID, content string) (*ent.QuickNote, error) {
	return r.client.QuickNote.UpdateOneID(id).
		SetContent(content).
		Save(ctx)
}

// Delete permanently deletes a quick note.
func (r *QuickNoteRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.QuickNote.DeleteOneID(id).Exec(ctx)
}

// ListByCard returns all quick notes for a card.
func (r *QuickNoteRepository) ListByCard(ctx context.Context, cardID uuid.UUID) ([]*ent.QuickNote, error) {
	return r.client.QuickNote.Query().
		Where(quicknote.CardIDEQ(cardID)).
		Order(ent.Desc(quicknote.FieldCreatedAt)).
		All(ctx)
}
