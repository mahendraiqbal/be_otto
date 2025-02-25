package api

import (
	"github.com/gorilla/mux"
	"github.com/mahendraiqbal/be_otto/internal/api/handlers"
)

func SetupRoutes(handler *handlers.Handler) *mux.Router {
	r := mux.NewRouter()

	// Brand routes
	r.HandleFunc("/brand", handler.CreateBrand).Methods("POST")

	// Voucher routes
	r.HandleFunc("/voucher", handler.CreateVoucher).Methods("POST")
	r.HandleFunc("/voucher", handler.GetVoucherByID).Methods("GET")
	r.HandleFunc("/voucher/brand", handler.GetVouchersByBrandID).Methods("GET")

	// Transaction routes
	r.HandleFunc("/transaction/redemption", handler.CreateRedemption).Methods("POST")
	r.HandleFunc("/transaction/redemption", handler.GetRedemptionByID).Methods("GET")

	return r
}
