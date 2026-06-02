package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/ent"
	"github.com/tendo-mulira/tnotes-teams/internal/ent/workspace"
	"github.com/tendo-mulira/tnotes-teams/internal/ent/workspacemember"
)

// WorkspaceRepository handles workspace persistence.
type WorkspaceRepository struct {
	client *ent.Client
	db     *sql.DB
}

// NewWorkspaceRepository creates a new WorkspaceRepository.
func NewWorkspaceRepository(client *ent.Client, db *sql.DB) *WorkspaceRepository {
	return &WorkspaceRepository{client: client, db: db}
}

// Create creates a new workspace.
func (r *WorkspaceRepository) Create(ctx context.Context, name, description string, createdBy uuid.UUID) (*ent.Workspace, error) {
	return r.client.Workspace.Create().
		SetName(name).
		SetDescription(description).
		SetCreatedBy(createdBy).
		Save(ctx)
}

// GetByID returns a workspace by ID.
func (r *WorkspaceRepository) GetByID(ctx context.Context, id uuid.UUID) (*ent.Workspace, error) {
	return r.client.Workspace.Get(ctx, id)
}

// ListByUser returns all workspaces a user is a member of.
func (r *WorkspaceRepository) ListByUser(ctx context.Context, userID uuid.UUID) ([]*ent.Workspace, error) {
	return r.client.Workspace.Query().
		Where(
			workspace.HasMembersWith(
				// Query through the edge
			),
		).
		All(ctx)
}

// ListByUserViaMembers returns workspaces using raw SQL for better join support.
func (r *WorkspaceRepository) ListByUserViaMembers(ctx context.Context, userID uuid.UUID) ([]*ent.Workspace, error) {
	// Use the workspace_member edge to find workspaces
	memberEntries, err := r.client.WorkspaceMember.Query().
		Where(workspacemember.UserIDEQ(userID)).
		WithWorkspace().
		All(ctx)
	if err != nil {
		return nil, err
	}

	workspaces := make([]*ent.Workspace, 0, len(memberEntries))
	for _, m := range memberEntries {
		if ws := m.Edges.Workspace; ws != nil {
			workspaces = append(workspaces, ws)
		}
	}
	return workspaces, nil
}

// AddMember adds a user to a workspace with a role.
func (r *WorkspaceRepository) AddMember(ctx context.Context, workspaceID, userID uuid.UUID, role string) (*ent.WorkspaceMember, error) {
	return r.client.WorkspaceMember.Create().
		SetWorkspaceID(workspaceID).
		SetUserID(userID).
		SetRole(workspacemember.Role(role)).
		Save(ctx)
}

// UpdateMemberRole updates a member's role in a workspace.
func (r *WorkspaceRepository) UpdateMemberRole(ctx context.Context, workspaceID, userID uuid.UUID, newRole string) error {
	_, err := r.client.WorkspaceMember.Update().
		Where(
			workspacemember.WorkspaceIDEQ(workspaceID),
			workspacemember.UserIDEQ(userID),
		).
		SetRole(workspacemember.Role(newRole)).
		Save(ctx)
	return err
}

// GetMemberRole returns the role of a user in a workspace.
func (r *WorkspaceRepository) GetMemberRole(ctx context.Context, workspaceID, userID uuid.UUID) (string, error) {
	member, err := r.client.WorkspaceMember.Query().
		Where(
			workspacemember.WorkspaceIDEQ(workspaceID),
			workspacemember.UserIDEQ(userID),
		).
		Only(ctx)
	if err != nil {
		return "", err
	}
	return string(member.Role), nil
}

// GetMembers returns all members of a workspace with user data.
func (r *WorkspaceRepository) GetMembers(ctx context.Context, workspaceID uuid.UUID) ([]*ent.WorkspaceMember, error) {
	return r.client.WorkspaceMember.Query().
		Where(workspacemember.WorkspaceIDEQ(workspaceID)).
		WithUser().
		All(ctx)
}

// RemoveMember removes a user from a workspace.
func (r *WorkspaceRepository) RemoveMember(ctx context.Context, workspaceID, userID uuid.UUID) error {
	_, err := r.client.WorkspaceMember.Delete().
		Where(
			workspacemember.WorkspaceIDEQ(workspaceID),
			workspacemember.UserIDEQ(userID),
		).
		Exec(ctx)
	return err
}
