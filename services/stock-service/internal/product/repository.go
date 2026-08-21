package product

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrProductCodeExists = errors.New("product code already exists")
	ErrInsufficientStock = errors.New("insufficient product stock")
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(
	ctx context.Context,
	input CreateProductInput,
) (*Product, error) {
	query := `
		INSERT INTO products (code, description, stock)
		VALUES ($1, $2, $3)
		RETURNING id, code, description, stock, created_at, updated_at
	`

	var product Product

	err := r.db.QueryRow(
		ctx,
		query,
		input.Code,
		input.Description,
		input.Stock,
	).Scan(
		&product.ID,
		&product.Code,
		&product.Description,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrProductCodeExists
		}

		return nil, fmt.Errorf("failed to create product: %w", err)
	}
	return &product, nil
}

func (r *Repository) List(ctx context.Context) ([]Product, error) {
	query := `
		SELECT id, code, description, stock, created_at, updated_at
		FROM products
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	products := make([]Product, 0)

	for rows.Next() {
		var product Product

		if err := rows.Scan(
			&product.ID,
			&product.Code,
			&product.Description,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed while reading products: %w", err)
	}

	return products, nil
}

func (r *Repository) FindByID(
	ctx context.Context,
	id int64,
) (*Product, error) {
	query := `
		SELECT id, code, description, stock, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var product Product

	err := r.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.Code,
		&product.Description,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrProductNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	return &product, nil
}

func (r *Repository) DecreaseStock(
	ctx context.Context,
	id int64,
	quantity int,
) (*Product, error) {
	query := `
		UPDATE products
		SET
			stock = stock - $1,
			updated_at = NOW()
		WHERE id = $2
		  AND stock >= $1
		RETURNING id, code, description, stock, created_at, updated_at
	`

	var product Product

	err := r.db.QueryRow(
		ctx,
		query,
		quantity,
		id,
	).Scan(
		&product.ID,
		&product.Code,
		&product.Description,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		_, findErr := r.FindByID(ctx, id)

		if errors.Is(findErr, ErrProductNotFound) {
			return nil, ErrProductNotFound
		}

		if findErr != nil {
			return nil, fmt.Errorf(
				"failed to verify product after stock update: %w",
				findErr,
			)
		}

		return nil, ErrInsufficientStock
	}

	if err != nil {
		return nil, fmt.Errorf("failed to decrease product stock: %w", err)
	}

	return &product, nil
}

func (r *Repository) DecreaseStockBatch(
	ctx context.Context,
	items []DecreaseStockBatchItemInput,
) ([]Product, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to begin stock transaction: %w",
			err,
		)
	}

	defer tx.Rollback(ctx)

	updatedProducts := make([]Product, 0, len(items))

	for _, item := range items {
		var product Product

		err := tx.QueryRow(
			ctx,
			`
				UPDATE products
				SET
					stock = stock - $1,
					updated_at = NOW()
				WHERE id = $2
				  AND stock >= $1
				RETURNING
					id,
					code,
					description,
					stock,
					created_at,
					updated_at
			`,
			item.Quantity,
			item.ProductID,
		).Scan(
			&product.ID,
			&product.Code,
			&product.Description,
			&product.Stock,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if errors.Is(err, pgx.ErrNoRows) {
			var exists bool

			checkErr := tx.QueryRow(
				ctx,
				`
					SELECT EXISTS(
						SELECT 1
						FROM products
						WHERE id = $1
					)
				`,
				item.ProductID,
			).Scan(&exists)

			if checkErr != nil {
				return nil, fmt.Errorf(
					"failed to verify product: %w",
					checkErr,
				)
			}

			if !exists {
				return nil, ErrProductNotFound
			}

			return nil, ErrInsufficientStock
		}

		if err != nil {
			return nil, fmt.Errorf(
				"failed to decrease product stock: %w",
				err,
			)
		}

		updatedProducts = append(
			updatedProducts,
			product,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(
			"failed to commit stock transaction: %w",
			err,
		)
	}

	return updatedProducts, nil
}
