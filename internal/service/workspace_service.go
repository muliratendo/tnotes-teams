package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/ent"
	"github.com/tendo-mulira/tnotes-teams/internal/repository"
)

// WorkspaceService handles workspace business logic.
type WorkspaceService struct {
	repos *repository.Repositories
}

// NewWorkspaceService creates a new WorkspaceService.
func NewWorkspaceService(repos *repository.Repositories) *WorkspaceService {
	return &WorkspaceService{repos: repos}
}

// WorkspaceDTO is the public workspace data.
type WorkspaceDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ThemeColor  string    `json:"theme_color"`
	CreatedBy   uuid.UUID `json:"created_by"`
	Role        string    `json:"role,omitempty"`
}

// MemberDTO is the public member data.
type MemberDTO struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Role      string    `json:"role"`
}

// CreateWorkspaceInput is the input for creating a workspace.
type CreateWorkspaceInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Create creates a new workspace and adds the creator as admin.
func (s *WorkspaceService) Create(ctx context.Context, input CreateWorkspaceInput, userID uuid.UUID) (*WorkspaceDTO, error) {
	if input.Name == "" {
		return nil, errors.New("workspace name is required")
	}

	ws, err := s.repos.Workspace.Create(ctx, input.Name, input.Description, userID)
	if err != nil {
		return nil, err
	}

	// Add creator as admin
	_, err = s.repos.Workspace.AddMember(ctx, ws.ID, userID, "admin")
	if err != nil {
		return nil, err
	}

	return &WorkspaceDTO{
		ID:          ws.ID,
		Name:        ws.Name,
		Description: ws.Description,
		ThemeColor:  ws.ThemeColor,
		CreatedBy:   ws.CreatedBy,
		Role:        "admin",
	}, nil
}

// GetByID returns a workspace by ID.
func (s *WorkspaceService) GetByID(ctx context.Context, id uuid.UUID) (*WorkspaceDTO, error) {
	ws, err := s.repos.Workspace.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &WorkspaceDTO{
		ID:          ws.ID,
		Name:        ws.Name,
		Description: ws.Description,
		ThemeColor:  ws.ThemeColor,
		CreatedBy:   ws.CreatedBy,
	}, nil
}

// ListByUser returns all workspaces a user belongs to, including their role.
func (s *WorkspaceService) ListByUser(ctx context.Context, userID uuid.UUID) ([]WorkspaceDTO, error) {
	workspaces, err := s.repos.Workspace.ListByUserViaMembers(ctx, userID)
	if err != nil {
		return nil, err
	}

	dtos := make([]WorkspaceDTO, 0, len(workspaces))
	for _, ws := range workspaces {
		role, _ := s.repos.Workspace.GetMemberRole(ctx, ws.ID, userID)
		dtos = append(dtos, WorkspaceDTO{
			ID:          ws.ID,
			Name:        ws.Name,
			Description: ws.Description,
			ThemeColor:  ws.ThemeColor,
			CreatedBy:   ws.CreatedBy,
			Role:        role,
		})
	}
	return dtos, nil
}

// InviteMember adds a user to a workspace by email.
type InviteMemberInput struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (s *WorkspaceService) InviteMember(ctx context.Context, workspaceID uuid.UUID, input InviteMemberInput) (*MemberDTO, error) {
	if input.Email == "" {
		return nil, errors.New("email is required")
	}
	if input.Role == "" {
		input.Role = "member"
	}

	// Find user by email
	user, err := s.repos.User.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("user not found with that email")
	}

	// Check if already a member
	existingRole, err := s.repos.Workspace.GetMemberRole(ctx, workspaceID, user.ID)
	if err == nil && existingRole != "" {
		return nil, errors.New("user is already a member of this workspace")
	}

	// Add member
	_, err = s.repos.Workspace.AddMember(ctx, workspaceID, user.ID, input.Role)
	if err != nil {
		return nil, err
	}

	return &MemberDTO{
		UserID:    user.ID,
		Email:     user.Email,
		Username:  user.Username,
		AvatarURL: user.AvatarURL,
		Role:      input.Role,
	}, nil
}

// UpdateMemberRole changes a member's role.
func (s *WorkspaceService) UpdateMemberRole(ctx context.Context, workspaceID, userID uuid.UUID, newRole string) error {
	validRoles := map[string]bool{"admin": true, "member": true, "viewer": true}
	if !validRoles[newRole] {
		return errors.New("invalid role")
	}
	return s.repos.Workspace.UpdateMemberRole(ctx, workspaceID, userID, newRole)
}

// GetMembers returns all members of a workspace.
func (s *WorkspaceService) GetMembers(ctx context.Context, workspaceID uuid.UUID) ([]MemberDTO, error) {
	members, err := s.repos.Workspace.GetMembers(ctx, workspaceID)
	if err != nil {
		return nil, err
	}

	dtos := make([]MemberDTO, 0, len(members))
	for _, m := range members {
		dto := MemberDTO{
			UserID: m.UserID,
			Role:   string(m.Role),
		}
		if u := m.Edges.User; u != nil {
			dto.Email = u.Email
			dto.Username = u.Username
			dto.AvatarURL = u.AvatarURL
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}

func toWorkspaceDTO(ws *ent.Workspace) WorkspaceDTO {
	return WorkspaceDTO{
		ID:          ws.ID,
		Name:        ws.Name,
		Description: ws.Description,
		ThemeColor:  ws.ThemeColor,
		CreatedBy:   ws.CreatedBy,
	}
}
