package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/tendo-mulira/tnotes-teams/internal/middleware"
	"github.com/tendo-mulira/tnotes-teams/internal/service"
	"github.com/tendo-mulira/tnotes-teams/internal/utils"
)

// CardHandler handles Card and nested sub-resource endpoints.
type CardHandler struct {
	service *service.CardService
}

// NewCardHandler creates a new CardHandler.
func NewCardHandler(service *service.CardService) *CardHandler {
	return &CardHandler{service: service}
}

// Create creates a new card in a list.
func (h *CardHandler) Create(w http.ResponseWriter, r *http.Request) {
	listIDStr := chi.URLParam(r, "listID")
	listID, err := uuid.Parse(listIDStr)
	if err != nil {
		utils.ValidationError(w, "invalid list ID")
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	var input service.CreateCardInput
	if err := utils.DecodeJSON(r, &input); err != nil {
		utils.ValidationError(w, "invalid request body")
		return
	}

	card, err := h.service.CreateCard(r.Context(), listID, input, userID)
	if err != nil {
		utils.ValidationError(w, err.Error())
		return
	}

	utils.Created(w, card)
}

// GetByID returns the details of a single card.
func (h *CardHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "cardID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ValidationError(w, "invalid card ID")
		return
	}

	card, err := h.service.GetCardWithDetails(r.Context(), id)
	if err != nil {
		utils.NotFound(w, "card not found")
		return
	}

	utils.JSON(w, http.StatusOK, card)
}

// Update updates card scalar fields.
func (h *CardHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "cardID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ValidationError(w, "invalid card ID")
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	var input service.UpdateCardInput
	if err := utils.DecodeJSON(r, &input); err != nil {
		utils.ValidationError(w, "invalid request body")
		return
	}

	card, err := h.service.UpdateCard(r.Context(), id, input, userID)
	if err != nil {
		utils.ValidationError(w, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, card)
}

// Move handles drag-and-drop position changes.
func (h *CardHandler) Move(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "cardID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ValidationError(w, "invalid card ID")
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	var input service.MoveCardInput
	if err := utils.DecodeJSON(r, &input); err != nil {
		utils.ValidationError(w, "invalid request body")
		return
	}

	card, err := h.service.MoveCard(r.Context(), id, input, userID)
	if err != nil {
		utils.ValidationError(w, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, card)
}

// Delete permanently removes a card.
func (h *CardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "cardID")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ValidationError(w, "invalid card ID")
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	err = h.service.DeleteCard(r.Context(), id, userID)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(w, "card deleted successfully", nil)
}

// CreateSubtask creates a new subtask checklist item.
func (h *CardHandler) CreateSubtask(w http.ResponseWriter, r *http.Request) {
	cardIDStr := chi.URLParam(r, "cardID")
	cardID, err := uuid.Parse(cardIDStr)
	if err != nil {
		utils.ValidationError(w, "invalid card ID")
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	var input service.CreateSubtaskInput
	if err := utils.DecodeJSON(r, &input); err != nil {
		utils.ValidationError(w, "invalid request body")
		return
	}

	st, err := h.service.CreateSubtask(r.Context(), cardID, input, userID)
	if err != nil {
		utils.ValidationError(w, err.Error())
		return
	}

	utils.Created(w, st)
}

// ToggleSubtask toggles completion status.
func (h *CardHandler) ToggleSubtask(w http.ResponseWriter, r *http.Request) {
	subtaskIDStr := chi.URLParam(r, "subtaskID")
	subtaskID, err := uuid.Parse(subtaskIDStr)
	if err != nil {
		utils.ValidationError(w, "invalid subtask ID")
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	st, err := h.service.ToggleSubtask(r.Context(), subtaskID, userID)
	if err != nil {
		utils.ValidationError(w, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, st)
}

// AddNote adds a quick note.
func (h *CardHandler) AddNote(w http.ResponseWriter, r *http.Request) {
	cardIDStr := chi.URLParam(r, "cardID")
	cardID, err := uuid.Parse(cardIDStr)
	if err != nil {
		utils.ValidationError(w, "invalid card ID")
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	var input service.CreateNoteInput
	if err := utils.DecodeJSON(r, &input); err != nil {
		utils.ValidationError(w, "invalid request body")
		return
	}

	note, err := h.service.CreateQuickNote(r.Context(), cardID, input, userID)
	if err != nil {
		utils.ValidationError(w, err.Error())
		return
	}

	utils.Created(w, note)
}

// AddComment uploads a comment.
func (h *CardHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	cardIDStr := chi.URLParam(r, "cardID")
	cardID, err := uuid.Parse(cardIDStr)
	if err != nil {
		utils.ValidationError(w, "invalid card ID")
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	var input service.CreateCommentInput
	if err := utils.DecodeJSON(r, &input); err != nil {
		utils.ValidationError(w, "invalid request body")
		return
	}

	comment, err := h.service.CreateComment(r.Context(), cardID, input, userID)
	if err != nil {
		utils.ValidationError(w, err.Error())
		return
	}

	utils.Created(w, comment)
}
