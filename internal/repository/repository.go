package repository

import (
	"context"

	"github.com/mahendraiqbal/be_otto/internal/models"
)

type Repository interface {
	CreateBrand(ctx context.Context, brand *models.Brand) error
	BrandExists(ctx context.Context, id int) (bool, error)
	CreateVoucher(ctx context.Context, voucher *models.Voucher) error
	GetVoucherByID(ctx context.Context, id int) (*models.Voucher, error)
	GetVouchersByBrandID(ctx context.Context, brandID int) ([]models.Voucher, error)
	CreateRedemption(ctx context.Context, req *models.RedemptionRequest) (*models.RedemptionResponse, error)
	GetRedemptionByID(ctx context.Context, id int) (*models.RedemptionResponse, error)
}
