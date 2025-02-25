package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mahendraiqbal/be_otto/internal/models"
)

func (h *Handler) CreateBrand(w http.ResponseWriter, r *http.Request) {
	var brand models.Brand
	if err := json.NewDecoder(r.Body).Decode(&brand); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validator.Struct(brand); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Validation failed")
		return
	}

	if err := h.repo.CreateBrand(r.Context(), &brand); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create brand")
		return
	}

	h.respondWithJSON(w, http.StatusCreated, brand)
}
