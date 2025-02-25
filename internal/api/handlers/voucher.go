package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mahendraiqbal/be_otto/internal/models"
)

func (h *Handler) CreateVoucher(w http.ResponseWriter, r *http.Request) {
	var voucher models.Voucher
	if err := json.NewDecoder(r.Body).Decode(&voucher); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.validator.Struct(voucher); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Validation failed: "+err.Error())
		return
	}

	// Check if brand exists
	exists, err := h.repo.BrandExists(r.Context(), voucher.BrandID)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to check brand: "+err.Error())
		return
	}
	if !exists {
		h.respondWithError(w, http.StatusBadRequest, "Brand does not exist")
		return
	}

	if err := h.repo.CreateVoucher(r.Context(), &voucher); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to create voucher: "+err.Error())
		return
	}

	h.respondWithJSON(w, http.StatusCreated, voucher)
}

func (h *Handler) GetVoucherByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	voucherID, err := strconv.Atoi(id)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid voucher ID")
		return
	}

	voucher, err := h.repo.GetVoucherByID(r.Context(), voucherID)
	if err != nil {
		h.respondWithError(w, http.StatusNotFound, "Voucher not found")
		return
	}

	h.respondWithJSON(w, http.StatusOK, voucher)
}

func (h *Handler) GetVouchersByBrandID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	brandID, err := strconv.Atoi(id)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "Invalid brand ID")
		return
	}

	vouchers, err := h.repo.GetVouchersByBrandID(r.Context(), brandID)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "Failed to fetch vouchers")
		return
	}

	h.respondWithJSON(w, http.StatusOK, vouchers)
}
