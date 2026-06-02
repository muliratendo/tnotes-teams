package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/middleware"
	"github.com/tendo-mulira/tnotes-teams/internal/service"
	"github.com/tendo-mulira/tnotes-teams/internal/utils"
)

// SyncHandler handles synchronization requests for offline mutations.
type SyncHandler struct {
	service *service.SyncService
}

// NewSyncHandler creates a new SyncHandler.
func NewSyncHandler(service *service.SyncService) *SyncHandler {
	return &SyncHandler{service: service}
}

// SyncProcesses syncs a batch of offline operations.
func (h *SyncHandler) Sync(w http.ResponseWriter, r *http.Request) {
	boardIDStr := chi.URLParam(r, "boardID")
	_, err := uuid.Parse(boardIDStr)
	if err != nil {
		utils.ValidationError(w, "invalid board ID")
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	var input service.SyncBatchInput
	if err := utils.DecodeJSON(r, &input); err != nil {
		utils.ValidationError(w, "invalid request body")
		return
	}

	response, err := h.service.ProcessBatch(r.Context(), input, userID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, response)
}
