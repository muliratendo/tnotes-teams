package handler

import (
	"github.com/tendo-mulira/tnotes-teams/internal/service"
)

// Handlers aggregates all API handlers.
type Handlers struct {
	Auth      *AuthHandler
	Workspace *WorkspaceHandler
	Board     *BoardHandler
	Card      *CardHandler
	Sync      *SyncHandler
}

// NewHandlers creates a new Handlers instance.
func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		Auth:      NewAuthHandler(services.Auth),
		Workspace: NewWorkspaceHandler(services.Workspace),
		Board:     NewBoardHandler(services.Board, services.Export),
		Card:      NewCardHandler(services.Card),
		Sync:      NewSyncHandler(services.Sync),
	}
}
