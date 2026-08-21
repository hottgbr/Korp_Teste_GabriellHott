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
