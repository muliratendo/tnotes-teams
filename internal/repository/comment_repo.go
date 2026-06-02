package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/ent"
	"github.com/tendo-mulira/tnotes-teams/internal/ent/comment"
)

// CommentRepository handles comment persistence.
type CommentRepository struct {
	client *ent.Client
	db     *sql.DB
}

// NewCommentRepository creates a new CommentRepository.
func NewCommentRepository(client *ent.Client, db *sql.DB) *CommentRepository {
	return &CommentRepository{client: client, db: db}
}

// Create creates a new comment.
func (r *CommentRepository) Create(ctx context.Context, content string, cardID, userID uuid.UUID) (*ent.Comment, error) {
	return r.client.Comment.Create().
		SetContent(content).
		SetCardID(cardID).
		SetUserID(userID).
		Save(ctx)
}

// ListByCard returns all comments for a card with user data.
func (r *CommentRepository) ListByCard(ctx context.Context, cardID uuid.UUID) ([]*ent.Comment, error) {
	return r.client.Comment.Query().
		Where(comment.CardIDEQ(cardID)).
		WithUser().
		Order(ent.Asc(comment.FieldCreatedAt)).
		All(ctx)
}

// Delete permanently deletes a comment.
func (r *CommentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.client.Comment.DeleteOneID(id).Exec(ctx)
}
