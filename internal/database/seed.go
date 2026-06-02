package database

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/ent"
	"github.com/tendo-mulira/tnotes-teams/internal/ent/user"
	"github.com/tendo-mulira/tnotes-teams/internal/ent/workspacemember"
	"github.com/tendo-mulira/tnotes-teams/internal/utils"
)

// SeedDevelopment inserts demo users and a sample board when the database is empty.
func SeedDevelopment(ctx context.Context, client *ent.Client) error {
	count, err := client.User.Query().Count(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	log.Println("Seeding development database...")

	hash, _ := utils.HashPassword("password123")

	admin, err := client.User.Create().
		SetEmail("admin@tnotes.dev").
		SetUsername("admin").
		SetPasswordHash(hash).
		SetAvatarURL("https://api.dicebear.com/7.x/bottts/svg?seed=admin").
		Save(ctx)
	if err != nil {
		return err
	}

	ws, err := client.Workspace.Create().
		SetName("TNotes Demo Team").
		SetDescription("Sample workspace for local development").
		SetThemeColor("#6366F1").
		SetCreatedBy(admin.ID).
		Save(ctx)
	if err != nil {
		return err
	}

	_, err = client.WorkspaceMember.Create().
		SetWorkspaceID(ws.ID).
		SetUserID(admin.ID).
		SetRole("admin").
		Save(ctx)
	if err != nil {
		return err
	}

	board, err := client.Board.Create().
		SetName("Sprint Board").
		SetDescription("Demo Kanban with sample cards").
		SetWorkspaceID(ws.ID).
		SetCreatedBy(admin.ID).
		Save(ctx)
	if err != nil {
		return err
	}

	listTitles := []string{"To Do", "In Progress", "Done"}
	var listIDs []uuid.UUID
	for i, title := range listTitles {
		l, err := client.List.Create().
			SetTitle(title).
			SetBoardID(board.ID).
			SetPosition(i).
			Save(ctx)
		if err != nil {
			return err
		}
		listIDs = append(listIDs, l.ID)
	}

	cardSpecs := []struct {
		listIdx int
		title   string
		desc    string
	}{
		{0, "Design WebSocket hub", "Architecture for real-time co-presence"},
		{0, "IndexedDB offline queue", "Persist mutations while offline"},
		{1, "Implement RBAC middleware", "Admin / Member / Viewer enforcement"},
		{2, "Project scaffold", "Go + Ent + sqlc + React"},
	}

	for i, spec := range cardSpecs {
		c, err := client.Card.Create().
			SetTitle(spec.title).
			SetDescription(spec.desc).
			SetListID(listIDs[spec.listIdx]).
			SetPosition(i).
			SetCreatedBy(admin.ID).
			SetProgressPercentage(0).
			Save(ctx)
		if err != nil {
			return err
		}

		subtasks := []string{"Research", "Implement", "Test"}
		for j, st := range subtasks {
			done := j == 0
			_, err = client.Subtask.Create().
				SetTitle(st).
				SetCardID(c.ID).
				SetPosition(j).
				SetIsCompleted(done).
				Save(ctx)
			if err != nil {
				return err
			}
		}
	}

	// Mock teammates for co-presence simulation
	mockUsers := []struct {
		email    string
		username string
		role     string
	}{
		{"sarah@tnotes.dev", "sarah", "member"},
		{"alex@tnotes.dev", "alex", "member"},
		{"michael@tnotes.dev", "michael", "viewer"},
	}

	for _, mu := range mockUsers {
		u, err := client.User.Create().
			SetEmail(mu.email).
			SetUsername(mu.username).
			SetPasswordHash(hash).
			SetAvatarURL("https://api.dicebear.com/7.x/bottts/svg?seed=" + mu.username).
			Save(ctx)
		if err != nil {
			if ent.IsConstraintError(err) {
				u, _ = client.User.Query().Where(user.EmailEQ(mu.email)).Only(ctx)
			} else {
				continue
			}
		}
		_, _ = client.WorkspaceMember.Create().
			SetWorkspaceID(ws.ID).
			SetUserID(u.ID).
			SetRole(workspacemember.Role(mu.role)).
			Save(ctx)
	}

	log.Println("Development seed complete (login: admin@tnotes.dev / password123)")
	return nil
}
