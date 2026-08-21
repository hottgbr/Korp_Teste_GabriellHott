package invoice

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(
	ctx context.Context,
	input CreateInvoiceInput,
) (*Invoice, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin invoice transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	var createdInvoice Invoice

	err = tx.QueryRow(
		ctx,
		`
			INSERT INTO invoices DEFAULT VALUES
			RETURNING id, number, status, created_at, updated_at
		`,
	).Scan(
		&createdInvoice.ID,
		&createdInvoice.Number,
		&createdInvoice.Status,
		&createdInvoice.CreatedAt,
		&createdInvoice.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create invoice: %w", err)
	}

	createdInvoice.Items = make([]InvoiceItem, 0, len(input.Items))

	for _, item := range input.Items {
		var createdItem InvoiceItem

		err = tx.QueryRow(
			ctx,
			`
				INSERT INTO invoice_items (
					invoice_id,
					product_id,
					quantity
				)
				VALUES ($1, $2, $3)
				RETURNING id, invoice_id, product_id, quantity
			`,
			createdInvoice.ID,
			item.ProductID,
			item.Quantity,
		).Scan(
			&createdItem.ID,
			&createdItem.InvoiceID,
			&createdItem.ProductID,
			&createdItem.Quantity,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"failed to create invoice item: %w",
				err,
			)
		}

		createdInvoice.Items = append(
			createdInvoice.Items,
			createdItem,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf(
			"failed to commit invoice transaction: %w",
			err,
		)
	}

	return &createdInvoice, nil
}
