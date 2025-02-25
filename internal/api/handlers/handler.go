package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mahendraiqbal/be_otto/internal/repository"
)

// @title           Voucher Management API
// @version         1.0
// @description     A RESTful API service for managing vouchers, brands, and redemptions.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// Common errors
var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrNotFound       = errors.New("resource not found")
	ErrServerError    = errors.New("internal server error")
)

// ErrorResponse represents an error response
// @Description Error response structure
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

// SuccessResponse represents a success response
// @Description Success response structure
type SuccessResponse struct {
	Message string      `json:"message" example:"operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

// Handler struct holds required services
type Handler struct {
	repo      repository.Repository
	validator *validator.Validate
}

// NewHandler creates a new handler instance
func NewHandler(repo repository.Repository) *Handler {
	return &Handler{
		repo:      repo,
		validator: validator.New(),
	}
}

// respondWithError sends an error response
func (h *Handler) respondWithError(w http.ResponseWriter, code int, message string) {
	response := ErrorResponse{
		Error: message,
	}
	h.respondWithJSON(w, code, response)
}

// respondWithJSON sends a JSON response
func (h *Handler) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// validateRequest validates the request body against a struct
func (h *Handler) validateRequest(data interface{}) error {
	if err := h.validator.Struct(data); err != nil {
		return err
	}
	return nil
}

// handleError handles common error cases
func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidRequest):
		h.respondWithError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, ErrNotFound):
		h.respondWithError(w, http.StatusNotFound, err.Error())
	default:
		h.respondWithError(w, http.StatusInternalServerError, ErrServerError.Error())
	}
}
