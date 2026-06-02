package handler

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/middleware"
	"github.com/tendo-mulira/tnotes-teams/internal/service"
	"github.com/tendo-mulira/tnotes-teams/internal/utils"
)

// BoardHandler handles board API endpoints.
type BoardHandler struct {
	boardService  *service.BoardService
	exportService *service.ExportService
}

// NewBoardHandler creates a new BoardHandler.
func NewBoardHandler(boardService *service.BoardService, exportService *service.ExportService) *BoardHandler {
	return &BoardHandler{
		boardService:  boardService,
		exportService: exportService,
	}
}

// Create creates a new board.
func (h *BoardHandler) Create(w http.ResponseWriter, r *http.Request) {
	wsIDStr := chi.URLParam(r, "workspaceID")
	wsID, err := uuid.Parse(wsIDStr)
	if err != nil {
		utils.ValidationError(w, "invalid workspace ID")
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	var input service.CreateBoardInput
	if err := utils.DecodeJSON(r, &input); err != nil {
		utils.ValidationError(w, "invalid request body")
		return
	}

	b, err := h.boardService.Create(r.Context(), wsID, input, userID)
	if err != nil {
		utils.ValidationError(w, err.Error())
		return
	}

	utils.Created(w, b)
}

// List lists all boards for a workspace.
func (h *BoardHandler) List(w http.ResponseWriter, r *http.Request) {
	wsIDStr := chi.URLParam(r, "workspaceID")
	wsID, err := uuid.Parse(wsIDStr)
	if err != nil {
		utils.ValidationError(w, "invalid workspace ID")
		return
	}

	boards, err := h.boardService.ListByWorkspace(r.Context(), wsID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, boards)
}

// GetByID returns the full board data (lists, cards, subtasks, notes).
func (h *BoardHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "boardID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ValidationError(w, "invalid board ID")
		return
	}

	b, err := h.boardService.GetWithFullData(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "board not found")
		return
	}

	utils.JSON(w, http.StatusOK, b)
}

// Update updates a board's settings.
func (h *BoardHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "boardID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ValidationError(w, "invalid board ID")
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	var input struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		ColorTheme  *string `json:"color_theme"`
	}

	if err := utils.DecodeJSON(r, &input); err != nil {
		utils.ValidationError(w, "invalid request body")
		return
	}

	b, err := h.boardService.Update(r.Context(), id, input.Name, input.Description, input.ColorTheme, userID)
	if err != nil {
		utils.ValidationError(w, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, b)
}

// Delete permanently deletes a board.
func (h *BoardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "boardID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ValidationError(w, "invalid board ID")
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	err = h.boardService.Delete(r.Context(), id, userID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(w, "board deleted successfully", nil)
}

// CreateList adds a new list column to a board.
func (h *BoardHandler) CreateList(w http.ResponseWriter, r *http.Request) {
	boardIDStr := chi.URLParam(r, "boardID")
	boardID, err := uuid.Parse(boardIDStr)
	if err != nil {
		utils.ValidationError(w, "invalid board ID")
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	var input service.CreateListInput
	if err := utils.DecodeJSON(r, &input); err != nil {
		utils.ValidationError(w, "invalid request body")
		return
	}

	l, err := h.boardService.CreateList(r.Context(), boardID, input, userID)
	if err != nil {
		utils.ValidationError(w, err.Error())
		return
	}

	utils.Created(w, l)
}

// Export JSON formats a board configuration.
func (h *BoardHandler) Export(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "boardID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ValidationError(w, "invalid board ID")
		return
	}

	export, err := h.exportService.ExportBoard(r.Context(), id)
	if err != nil {
		utils.NotFound(w, err.Error())
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=board-export.json")
	utils.JSON(w, http.StatusOK, export)
}

// ExportCSV formats a board configuration as CSV.
func (h *BoardHandler) ExportCSV(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "boardID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ValidationError(w, "invalid board ID")
		return
	}

	csv, err := h.exportService.ExportBoardCSV(r.Context(), id)
	if err != nil {
		utils.NotFound(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=board-export.csv")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(csv))
}

// Import restores a board from JSON.
func (h *BoardHandler) Import(w http.ResponseWriter, r *http.Request) {
	wsIDStr := chi.URLParam(r, "workspaceID")
	wsID, err := uuid.Parse(wsIDStr)
	if err != nil {
		utils.ValidationError(w, "invalid workspace ID")
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.ValidationError(w, "unable to read payload")
		return
	}

	board, err := h.exportService.ImportBoard(r.Context(), body, wsID, userID)
	if err != nil {
		utils.ValidationError(w, err.Error())
		return
	}

	utils.Created(w, board)
}
