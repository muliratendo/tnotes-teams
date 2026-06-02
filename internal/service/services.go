package service

import (
	"github.com/tendo-mulira/tnotes-teams/internal/config"
	"github.com/tendo-mulira/tnotes-teams/internal/repository"
	"github.com/tendo-mulira/tnotes-teams/internal/websocket"
)

// Services holds all service instances.
type Services struct {
	Auth      *AuthService
	Workspace *WorkspaceService
	Board     *BoardService
	Card      *CardService
	Sync      *SyncService
	Export    *ExportService
}

// NewServices creates all services.
func NewServices(repos *repository.Repositories, cfg *config.Config, hub *websocket.Hub) *Services {
	return &Services{
		Auth:      NewAuthService(repos, cfg),
		Workspace: NewWorkspaceService(repos),
		Board:     NewBoardService(repos),
		Card:      NewCardService(repos, hub),
		Sync:      NewSyncService(repos, hub),
		Export:    NewExportService(repos),
	}
}
