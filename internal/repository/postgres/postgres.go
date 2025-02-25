package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/mahendraiqbal/be_otto/internal/models"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(connStr string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return &PostgresRepository{
		db: db,
	}, nil
}

func (r *PostgresRepository) CreateBrand(ctx context.Context, brand *models.Brand) error {
	query := `
		INSERT INTO brands (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRowContext(
		ctx,
		query,
		brand.Name,
		brand.Description,
	).Scan(&brand.ID, &brand.CreatedAt, &brand.UpdatedAt)
}

func (r *PostgresRepository) CreateVoucher(ctx context.Context, voucher *models.Voucher) error {
	query := `
		INSERT INTO vouchers (brand_id, code, name, description, point_cost, stock, valid_until)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		voucher.BrandID,
		voucher.Code,
		voucher.Name,
		voucher.Description,
		voucher.PointCost,
		voucher.Stock,
		voucher.ValidUntil,
	).Scan(&voucher.ID, &voucher.CreatedAt, &voucher.UpdatedAt)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				return fmt.Errorf("brand_id %d does not exist", voucher.BrandID)
			case "unique_violation":
				return fmt.Errorf("voucher code %s already exists", voucher.Code)
			}
		}
		return fmt.Errorf("error creating voucher: %v", err)
	}
	return nil
}

func (r *PostgresRepository) GetVoucherByID(ctx context.Context, id int) (*models.Voucher, error) {
	voucher := &models.Voucher{}
	query := `
		SELECT id, brand_id, code, name, description, point_cost, stock, valid_until, created_at, updated_at
		FROM vouchers
		WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&voucher.ID,
		&voucher.BrandID,
		&voucher.Code,
		&voucher.Name,
		&voucher.Description,
		&voucher.PointCost,
		&voucher.Stock,
		&voucher.ValidUntil,
		&voucher.CreatedAt,
		&voucher.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("voucher not found")
	}
	if err != nil {
		return nil, err
	}

	return voucher, nil
}

func (r *PostgresRepository) GetVouchersByBrandID(ctx context.Context, brandID int) ([]models.Voucher, error) {
	query := `
		SELECT id, brand_id, code, name, description, point_cost, stock, valid_until, created_at, updated_at
		FROM vouchers
		WHERE brand_id = $1`

	rows, err := r.db.QueryContext(ctx, query, brandID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vouchers []models.Voucher
	for rows.Next() {
		var v models.Voucher
		err := rows.Scan(
			&v.ID,
			&v.BrandID,
			&v.Code,
			&v.Name,
			&v.Description,
			&v.PointCost,
			&v.Stock,
			&v.ValidUntil,
			&v.CreatedAt,
			&v.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		vouchers = append(vouchers, v)
	}

	return vouchers, nil
}

func (r *PostgresRepository) CreateRedemption(ctx context.Context, req *models.RedemptionRequest) (*models.RedemptionResponse, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// First check if customer exists and has enough points
	var customerPoints int
	err = tx.QueryRowContext(ctx, `
		SELECT points FROM customers WHERE id = $1`,
		req.CustomerID,
	).Scan(&customerPoints)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer points: %v", err)
	}

	// Create redemption record
	var redemptionID int
	err = tx.QueryRowContext(ctx, `
		INSERT INTO redemptions (customer_id, total_points, status)
		VALUES ($1, 0, 'PENDING')
		RETURNING id`,
		req.CustomerID,
	).Scan(&redemptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to create redemption: %v", err)
	}

	totalPoints := 0
	for _, item := range req.Items {
		// Get voucher details and check stock
		var pointCost int
		var voucherName string
		var currentStock int
		err := tx.QueryRowContext(ctx, `
			SELECT name, point_cost, stock FROM vouchers WHERE id = $1`,
			item.VoucherID,
		).Scan(&voucherName, &pointCost, &currentStock)
		if err != nil {
			return nil, fmt.Errorf("failed to get voucher details: %v", err)
		}

		// Check stock availability
		if currentStock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for voucher %s", voucherName)
		}

		itemTotalPoints := pointCost * item.Quantity
		totalPoints += itemTotalPoints

		// Create redemption item
		_, err = tx.ExecContext(ctx, `
			INSERT INTO redemption_items (redemption_id, voucher_id, quantity, points_per_unit, total_points)
			VALUES ($1, $2, $3, $4, $5)`,
			redemptionID, item.VoucherID, item.Quantity, pointCost, itemTotalPoints,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create redemption item: %v", err)
		}

		// Update voucher stock
		_, err = tx.ExecContext(ctx, `
			UPDATE vouchers
			SET stock = stock - $1
			WHERE id = $2`,
			item.Quantity, item.VoucherID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to update voucher stock: %v", err)
		}
	}

	// Check if customer has enough points
	if customerPoints < totalPoints {
		return nil, fmt.Errorf("insufficient points: customer has %d, needs %d", customerPoints, totalPoints)
	}

	// Update customer points
	_, err = tx.ExecContext(ctx, `
		UPDATE customers
		SET points = points - $1
		WHERE id = $2`,
		totalPoints, req.CustomerID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update customer points: %v", err)
	}

	// Update total points in redemption
	_, err = tx.ExecContext(ctx, `
		UPDATE redemptions
		SET total_points = $1, status = 'COMPLETED'
		WHERE id = $2`,
		totalPoints, redemptionID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update redemption: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// Return redemption response
	return r.GetRedemptionByID(ctx, redemptionID)
}

func (r *PostgresRepository) GetRedemptionByID(ctx context.Context, id int) (*models.RedemptionResponse, error) {
	var resp models.RedemptionResponse

	// Get redemption details
	err := r.db.QueryRowContext(ctx, `
		SELECT id, customer_id, total_points, status, created_at
		FROM redemptions
		WHERE id = $1`,
		id,
	).Scan(&resp.ID, &resp.CustomerID, &resp.TotalPoints, &resp.Status, &resp.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("redemption with ID %d not found", id)
		}
		return nil, fmt.Errorf("error fetching redemption: %v", err)
	}

	// Get redemption items
	rows, err := r.db.QueryContext(ctx, `
		SELECT ri.voucher_id, v.name, ri.quantity, ri.points_per_unit, ri.total_points
		FROM redemption_items ri
		JOIN vouchers v ON v.id = ri.voucher_id
		WHERE ri.redemption_id = $1`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("error fetching redemption items: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item models.RedemptionItemDetail
		err := rows.Scan(
			&item.VoucherID,
			&item.VoucherName,
			&item.Quantity,
			&item.PointsPerUnit,
			&item.TotalPoints,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning redemption item: %v", err)
		}
		resp.Items = append(resp.Items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating redemption items: %v", err)
	}

	return &resp, nil
}

func (r *PostgresRepository) BrandExists(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM brands WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking brand existence: %v", err)
	}
	return exists, nil
}
