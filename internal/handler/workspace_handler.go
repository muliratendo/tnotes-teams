package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/middleware"
	"github.com/tendo-mulira/tnotes-teams/internal/service"
	"github.com/tendo-mulira/tnotes-teams/internal/utils"
)

// WorkspaceHandler handles workspace API endpoints.
type WorkspaceHandler struct {
	service *service.WorkspaceService
}

// NewWorkspaceHandler creates a new WorkspaceHandler.
func NewWorkspaceHandler(service *service.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{service: service}
}

// Create creates a new workspace.
func (h *WorkspaceHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	var input service.CreateWorkspaceInput
	if err := utils.DecodeJSON(r, &input); err != nil {
		utils.ValidationError(w, "invalid request body")
		return
	}

	ws, err := h.service.Create(r.Context(), input, userID)
	if err != nil {
		utils.ValidationError(w, err.Error())
		return
	}

	utils.Created(w, ws)
}

// List lists all workspaces for the current user.
func (h *WorkspaceHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	workspaces, err := h.service.ListByUser(r.Context(), userID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, workspaces)
}

// GetByID gets details of a workspace.
func (h *WorkspaceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "workspaceID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ValidationError(w, "invalid workspace ID")
		return
	}

	ws, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "workspace not found")
		return
	}

	utils.JSON(w, http.StatusOK, ws)
}

// InviteMember invites a user to the workspace.
func (h *WorkspaceHandler) InviteMember(w http.ResponseWriter, r *http.Request) {
	wsIDStr := chi.URLParam(r, "workspaceID")
	wsID, err := uuid.Parse(wsIDStr)
	if err != nil {
		utils.ValidationError(w, "invalid workspace ID")
		return
	}

	var input service.InviteMemberInput
	if err := utils.DecodeJSON(r, &input); err != nil {
		utils.ValidationError(w, "invalid request body")
		return
	}

	member, err := h.service.InviteMember(r.Context(), wsID, input)
	if err != nil {
		utils.ValidationError(w, err.Error())
		return
	}

	utils.Created(w, member)
}

// UpdateMemberRole updates the role of a workspace member.
func (h *WorkspaceHandler) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	wsIDStr := chi.URLParam(r, "workspaceID")
	wsID, err := uuid.Parse(wsIDStr)
	if err != nil {
		utils.ValidationError(w, "invalid workspace ID")
		return
	}

	memberIDStr := chi.URLParam(r, "userID")
	memberID, err := uuid.Parse(memberIDStr)
	if err != nil {
		utils.ValidationError(w, "invalid user ID")
		return
	}

	type updateRoleInput struct {
		Role string `json:"role"`
	}

	var input updateRoleInput
	if err := utils.DecodeJSON(r, &input); err != nil {
		utils.ValidationError(w, "invalid request body")
		return
	}

	err = h.service.UpdateMemberRole(r.Context(), wsID, memberID, input.Role)
	if err != nil {
		utils.ValidationError(w, err.Error())
		return
	}

	utils.Success(w, "member role updated successfully", nil)
}

// ListMembers lists all members of the workspace.
func (h *WorkspaceHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	wsIDStr := chi.URLParam(r, "workspaceID")
	wsID, err := uuid.Parse(wsIDStr)
	if err != nil {
		utils.ValidationError(w, "invalid workspace ID")
		return
	}

	members, err := h.service.GetMembers(r.Context(), wsID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, members)
}
