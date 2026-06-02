package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/utils"
)

// RequireRole creates middleware that checks the user has at least the specified role
// for the workspace referenced in the URL.
// Role hierarchy: admin > member > viewer
func (m *Middleware) RequireRole(minRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserID(r.Context())
			if !ok {
				utils.Unauthorized(w, "authentication required")
				return
			}

			// Try to get workspace_id from URL params
			workspaceIDStr := chi.URLParam(r, "workspaceID")
			if workspaceIDStr == "" {
				// Try to resolve workspace from board ID
				boardIDStr := chi.URLParam(r, "boardID")
				if boardIDStr != "" {
					boardID, err := uuid.Parse(boardIDStr)
					if err != nil {
						utils.ValidationError(w, "invalid board ID")
						return
					}
					board, err := m.repos.Board.GetByID(r.Context(), boardID)
					if err != nil {
						utils.NotFound(w, "board not found")
						return
					}
					workspaceIDStr = board.WorkspaceID.String()
				}
			}

			// Try to resolve from list ID
			if workspaceIDStr == "" {
				listIDStr := chi.URLParam(r, "listID")
				if listIDStr != "" {
					listID, err := uuid.Parse(listIDStr)
					if err == nil {
						list, err := m.repos.List.GetByID(r.Context(), listID)
						if err == nil {
							board, err := m.repos.Board.GetByID(r.Context(), list.BoardID)
							if err == nil {
								workspaceIDStr = board.WorkspaceID.String()
							}
						}
					}
				}
			}

			// Try to resolve from card ID
			if workspaceIDStr == "" {
				cardIDStr := chi.URLParam(r, "cardID")
				if cardIDStr != "" {
					cardID, err := uuid.Parse(cardIDStr)
					if err == nil {
						card, err := m.repos.Card.GetByID(r.Context(), cardID)
						if err == nil {
							list, err := m.repos.List.GetByID(r.Context(), card.ListID)
							if err == nil {
								board, err := m.repos.Board.GetByID(r.Context(), list.BoardID)
								if err == nil {
									workspaceIDStr = board.WorkspaceID.String()
								}
							}
						}
					}
				}
			}

			if workspaceIDStr == "" {
				utils.ValidationError(w, "workspace context required")
				return
			}

			workspaceID, err := uuid.Parse(workspaceIDStr)
			if err != nil {
				utils.ValidationError(w, "invalid workspace ID")
				return
			}

			// Check user's role in the workspace
			role, err := m.repos.Workspace.GetMemberRole(r.Context(), workspaceID, userID)
			if err != nil {
				utils.Forbidden(w, "you are not a member of this workspace")
				return
			}

			if !hasMinimumRole(role, minRole) {
				utils.Forbidden(w, "insufficient permissions: requires "+minRole+" role")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// hasMinimumRole checks if userRole meets the minimum required role.
func hasMinimumRole(userRole, minRole string) bool {
	roleLevel := map[string]int{
		"viewer": 1,
		"member": 2,
		"admin":  3,
	}

	userLevel, ok := roleLevel[userRole]
	if !ok {
		return false
	}

	minLevel, ok := roleLevel[minRole]
	if !ok {
		return false
	}

	return userLevel >= minLevel
}
