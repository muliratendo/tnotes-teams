package handler

import (
	"net/http"

	"github.com/tendo-mulira/tnotes-teams/internal/middleware"
	"github.com/tendo-mulira/tnotes-teams/internal/service"
	"github.com/tendo-mulira/tnotes-teams/internal/utils"
)

// AuthHandler handles HTTP requests for user authentication.
type AuthHandler struct {
	service *service.AuthService
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register registers a new user.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input service.RegisterInput
	if err := utils.DecodeJSON(r, &input); err != nil {
		utils.ValidationError(w, "invalid request body")
		return
	}

	response, err := h.service.Register(r.Context(), input)
	if err != nil {
		utils.ValidationError(w, err.Error())
		return
	}

	utils.Created(w, response)
}

// Login authenticates a user.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input service.LoginInput
	if err := utils.DecodeJSON(r, &input); err != nil {
		utils.ValidationError(w, "invalid request body")
		return
	}

	response, err := h.service.Login(r.Context(), input)
	if err != nil {
		utils.ValidationError(w, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, response)
}

// Me returns the current authenticated user.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		utils.Unauthorized(w, "authentication required")
		return
	}

	user, err := h.service.GetCurrentUser(r.Context(), userID)
	if err != nil {
		utils.NotFound(w, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, user)
}
