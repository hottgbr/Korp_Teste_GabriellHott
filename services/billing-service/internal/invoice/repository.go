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

func (r *Repository) List(
	ctx context.Context,
) ([]Invoice, error) {
	rows, err := r.db.Query(
		ctx,
		`
			SELECT id, number, status, created_at, updated_at
			FROM invoices
			ORDER BY number
		`,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to list invoices: %w",
			err,
		)
	}
	defer rows.Close()

	invoices := make([]Invoice, 0)

	for rows.Next() {
		var invoice Invoice

		if err := rows.Scan(
			&invoice.ID,
			&invoice.Number,
			&invoice.Status,
			&invoice.CreatedAt,
			&invoice.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf(
				"failed to scan invoice: %w",
				err,
			)
		}

		invoices = append(invoices, invoice)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"failed while reading invoices: %w",
			err,
		)
	}

	return invoices, nil
}

func (r *Repository) FindByID(
	ctx context.Context,
	id int64,
) (*Invoice, error) {
	var invoice Invoice

	err := r.db.QueryRow(
		ctx,
		`
			SELECT id, number, status, created_at, updated_at
			FROM invoices
			WHERE id = $1
		`,
		id,
	).Scan(
		&invoice.ID,
		&invoice.Number,
		&invoice.Status,
		&invoice.CreatedAt,
		&invoice.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to find invoice: %w",
			err,
		)
	}

	rows, err := r.db.Query(
		ctx,
		`
			SELECT id, invoice_id, product_id, quantity
			FROM invoice_items
			WHERE invoice_id = $1
			ORDER BY id
		`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to list invoice items: %w",
			err,
		)
	}
	defer rows.Close()

	invoice.Items = make([]InvoiceItem, 0)

	for rows.Next() {
		var item InvoiceItem

		if err := rows.Scan(
			&item.ID,
			&item.InvoiceID,
			&item.ProductID,
			&item.Quantity,
		); err != nil {
			return nil, fmt.Errorf(
				"failed to scan invoice item: %w",
				err,
			)
		}

		invoice.Items = append(
			invoice.Items,
			item,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"failed while reading invoice items: %w",
			err,
		)
	}

	return &invoice, nil
}
