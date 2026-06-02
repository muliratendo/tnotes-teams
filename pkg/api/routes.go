package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tendo-mulira/tnotes-teams/internal/config"
	"github.com/tendo-mulira/tnotes-teams/internal/handler"
	internalmw "github.com/tendo-mulira/tnotes-teams/internal/middleware"
	"github.com/tendo-mulira/tnotes-teams/internal/service"
	"github.com/tendo-mulira/tnotes-teams/internal/websocket"
)

// SetupRoutes builds the application router and registers all middleware and routes.
func SetupRoutes(h *handler.Handlers, mw *internalmw.Middleware, hub *websocket.Hub, services *service.Services, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	// Standard middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(mw.CORS())

	// WebSocket upgrading endpoint
	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r, cfg)
	})

	// Public auth endpoints
	r.Group(func(r chi.Router) {
		r.Post("/api/auth/register", h.Auth.Register)
		r.Post("/api/auth/login", h.Auth.Login)
	})

	// Authenticated endpoints
	r.Group(func(r chi.Router) {
		r.Use(mw.Auth)

		r.Get("/api/auth/me", h.Auth.Me)

		// Workspaces CRUD
		r.Get("/api/workspaces", h.Workspace.List)
		r.Post("/api/workspaces", h.Workspace.Create)
		r.Get("/api/workspaces/{workspaceID}", h.Workspace.GetByID)

		// Workspace user actions (Admin required)
		r.With(mw.RequireRole("admin")).Post("/api/workspaces/{workspaceID}/members", h.Workspace.InviteMember)
		r.With(mw.RequireRole("admin")).Put("/api/workspaces/{workspaceID}/members/{userID}", h.Workspace.UpdateMemberRole)
		
		// Workspace members (Viewer+ required)
		r.With(mw.RequireRole("viewer")).Get("/api/workspaces/{workspaceID}/members", h.Workspace.ListMembers)

		// Workspace Boards (Viewer+ for list, Member+ for create, Admin for import)
		r.With(mw.RequireRole("viewer")).Get("/api/workspaces/{workspaceID}/boards", h.Board.List)
		r.With(mw.RequireRole("member")).Post("/api/workspaces/{workspaceID}/boards", h.Board.Create)
		r.With(mw.RequireRole("admin")).Post("/api/workspaces/{workspaceID}/import", h.Board.Import)

		// Specific Board operations
		r.With(mw.RequireRole("viewer")).Get("/api/boards/{boardID}", h.Board.GetByID)
		r.With(mw.RequireRole("member")).Put("/api/boards/{boardID}", h.Board.Update)
		r.With(mw.RequireRole("admin")).Delete("/api/boards/{boardID}", h.Board.Delete)
		r.With(mw.RequireRole("member")).Post("/api/boards/{boardID}/lists", h.Board.CreateList)

		// Data portability
		r.With(mw.RequireRole("viewer")).Get("/api/boards/{boardID}/export", h.Board.Export)
		r.With(mw.RequireRole("viewer")).Get("/api/boards/{boardID}/export/csv", h.Board.ExportCSV)
		
		// Offline Sync
		r.With(mw.RequireRole("member")).Post("/api/boards/{boardID}/sync", h.Sync.Sync)

		// Cards CRUD (Member+ for writes, Viewer+ for reads)
		r.With(mw.RequireRole("member")).Post("/api/lists/{listID}/cards", h.Card.Create)
		r.With(mw.RequireRole("viewer")).Get("/api/cards/{cardID}", h.Card.GetByID)
		r.With(mw.RequireRole("member")).Put("/api/cards/{cardID}", h.Card.Update)
		r.With(mw.RequireRole("member")).Put("/api/cards/{cardID}/move", h.Card.Move)
		r.With(mw.RequireRole("member")).Delete("/api/cards/{cardID}", h.Card.Delete)

		// Nested Card details (Member+ required)
		r.With(mw.RequireRole("member")).Post("/api/cards/{cardID}/subtasks", h.Card.CreateSubtask)
		r.With(mw.RequireRole("member")).Put("/api/subtasks/{subtaskID}", h.Card.ToggleSubtask)
		r.With(mw.RequireRole("member")).Post("/api/cards/{cardID}/notes", h.Card.AddNote)
		r.With(mw.RequireRole("member")).Post("/api/cards/{cardID}/comments", h.Card.AddComment)
	})

	return r
}
