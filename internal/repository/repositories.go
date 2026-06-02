package repository

import (
	"database/sql"

	"github.com/tendo-mulira/tnotes-teams/internal/ent"
)

// Repositories holds all repository instances.
type Repositories struct {
	User        *UserRepository
	Workspace   *WorkspaceRepository
	Board       *BoardRepository
	List        *ListRepository
	Card        *CardRepository
	Subtask     *SubtaskRepository
	Comment     *CommentRepository
	QuickNote   *QuickNoteRepository
	ActivityLog *ActivityLogRepository
}

// NewRepositories creates all repositories with shared Ent client and raw DB connection.
func NewRepositories(client *ent.Client, db *sql.DB) *Repositories {
	return &Repositories{
		User:        NewUserRepository(client, db),
		Workspace:   NewWorkspaceRepository(client, db),
		Board:       NewBoardRepository(client, db),
		List:        NewListRepository(client, db),
		Card:        NewCardRepository(client, db),
		Subtask:     NewSubtaskRepository(client, db),
		Comment:     NewCommentRepository(client, db),
		QuickNote:   NewQuickNoteRepository(client, db),
		ActivityLog: NewActivityLogRepository(client, db),
	}
}
