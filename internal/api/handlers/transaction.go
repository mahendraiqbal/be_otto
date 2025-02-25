package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mahendraiqbal/be_otto/internal/models"
)

func (h *Handler) CreateRedemption(w http.ResponseWriter, r *http.Request) {
	var req models.RedemptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	response, err := h.repo.CreateRedemption(r.Context(), &req)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create redemption: "+err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusCreated, response)
}

func (h *Handler) GetRedemptionByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("transactionId")
	redemptionID, err := strconv.Atoi(id)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid transaction ID")
		return
	}

	redemption, err := h.repo.GetRedemptionByID(r.Context(), redemptionID)
	if err != nil {
		h.respondWithError(w, http.StatusNotFound, "Failed to get transaction: "+err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusOK, redemption)
}
