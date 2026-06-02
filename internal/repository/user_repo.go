package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/ent"
	"github.com/tendo-mulira/tnotes-teams/internal/ent/user"
)

// UserRepository handles user persistence.
type UserRepository struct {
	client *ent.Client
	db     *sql.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(client *ent.Client, db *sql.DB) *UserRepository {
	return &UserRepository{client: client, db: db}
}

// Create creates a new user.
func (r *UserRepository) Create(ctx context.Context, email, username, passwordHash string) (*ent.User, error) {
	return r.client.User.Create().
		SetEmail(email).
		SetUsername(username).
		SetPasswordHash(passwordHash).
		Save(ctx)
}

// GetByID returns a user by ID.
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*ent.User, error) {
	return r.client.User.Get(ctx, id)
}

// GetByEmail returns a user by email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*ent.User, error) {
	return r.client.User.Query().
		Where(user.EmailEQ(email)).
		Only(ctx)
}

// GetByUsername returns a user by username.
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*ent.User, error) {
	return r.client.User.Query().
		Where(user.UsernameEQ(username)).
		Only(ctx)
}

// UpdateAvatar updates a user's avatar URL.
func (r *UserRepository) UpdateAvatar(ctx context.Context, id uuid.UUID, avatarURL string) (*ent.User, error) {
	return r.client.User.UpdateOneID(id).
		SetAvatarURL(avatarURL).
		Save(ctx)
}
