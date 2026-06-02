package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/config"
	"github.com/tendo-mulira/tnotes-teams/internal/ent"
	"github.com/tendo-mulira/tnotes-teams/internal/repository"
	"github.com/tendo-mulira/tnotes-teams/internal/utils"
)

// AuthService handles authentication and user management.
type AuthService struct {
	repos *repository.Repositories
	cfg   *config.Config
}

// NewAuthService creates a new AuthService.
func NewAuthService(repos *repository.Repositories, cfg *config.Config) *AuthService {
	return &AuthService{repos: repos, cfg: cfg}
}

// RegisterInput represents registration request data.
type RegisterInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginInput represents login request data.
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse is the response after successful auth.
type AuthResponse struct {
	Token string   `json:"token"`
	User  UserDTO  `json:"user"`
}

// UserDTO is the public user data.
type UserDTO struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url,omitempty"`
}

// Register creates a new user account.
func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*AuthResponse, error) {
	// Validate input
	if strings.TrimSpace(input.Email) == "" {
		return nil, errors.New("email is required")
	}
	if strings.TrimSpace(input.Username) == "" {
		return nil, errors.New("username is required")
	}
	if len(input.Password) < 6 {
		return nil, errors.New("password must be at least 6 characters")
	}

	// Hash password
	hash, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, errors.New("failed to process password")
	}

	// Create user
	user, err := s.repos.User.Create(ctx, input.Email, input.Username, hash)
	if err != nil {
		if ent.IsConstraintError(err) {
			return nil, errors.New("email or username already taken")
		}
		return nil, err
	}

	// Default workspace for new accounts
	ws, err := s.repos.Workspace.Create(ctx, "My Workspace", "Your team Kanban workspace", user.ID)
	if err == nil {
		_, _ = s.repos.Workspace.AddMember(ctx, ws.ID, user.ID, "admin")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Username, s.cfg.JWTSecret, s.cfg.JWTExpiration)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		Token: token,
		User:  toUserDTO(user),
	}, nil
}

// Login authenticates a user.
func (s *AuthService) Login(ctx context.Context, input LoginInput) (*AuthResponse, error) {
	if strings.TrimSpace(input.Email) == "" {
		return nil, errors.New("email is required")
	}
	if strings.TrimSpace(input.Password) == "" {
		return nil, errors.New("password is required")
	}

	// Find user
	user, err := s.repos.User.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !utils.CheckPassword(input.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email, user.Username, s.cfg.JWTSecret, s.cfg.JWTExpiration)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &AuthResponse{
		Token: token,
		User:  toUserDTO(user),
	}, nil
}

// GetCurrentUser returns the authenticated user's data.
func (s *AuthService) GetCurrentUser(ctx context.Context, userID uuid.UUID) (*UserDTO, error) {
	user, err := s.repos.User.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	dto := toUserDTO(user)
	return &dto, nil
}

func toUserDTO(u *ent.User) UserDTO {
	return UserDTO{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		AvatarURL: u.AvatarURL,
	}
}
