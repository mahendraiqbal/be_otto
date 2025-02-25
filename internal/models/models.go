package models

import (
	"time"
)

type Brand struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Voucher struct {
	ID          int       `json:"id"`
	BrandID     int       `json:"brand_id" validate:"required"`
	Code        string    `json:"code" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Description string    `json:"description"`
	PointCost   int       `json:"point_cost" validate:"required,gt=0"`
	Stock       int       `json:"stock" validate:"required,gte=0"`
	ValidUntil  time.Time `json:"valid_until"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type RedemptionRequest struct {
	CustomerID int              `json:"customer_id" validate:"required"`
	Items      []RedemptionItem `json:"items" validate:"required,dive"`
}

type RedemptionItem struct {
	VoucherID int `json:"voucher_id" validate:"required"`
	Quantity  int `json:"quantity" validate:"required,min=1"`
}

type RedemptionResponse struct {
	ID          int                    `json:"id"`
	CustomerID  int                    `json:"customer_id"`
	TotalPoints int                    `json:"total_points"`
	Status      string                 `json:"status"`
	Items       []RedemptionItemDetail `json:"items"`
	CreatedAt   time.Time              `json:"created_at"`
}

type RedemptionItemDetail struct {
	VoucherID     int    `json:"voucher_id"`
	VoucherName   string `json:"voucher_name"`
	Quantity      int    `json:"quantity"`
	PointsPerUnit int    `json:"points_per_unit"`
	TotalPoints   int    `json:"total_points"`
}
