package invoice

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrInvoiceItemsRequired = errors.New(
		"invoice must contain at least one item",
	)

	ErrInvalidProductID = errors.New(
		"product id must be greater than zero",
	)

	ErrInvalidQuantity = errors.New(
		"item quantity must be greater than zero",
	)
	ErrInvoiceAlreadyClosed = errors.New(
		"invoice is already closed",
	)
	ErrStockUpdateFailed = errors.New(
		"failed to update product stock",
	)
)

type Service struct {
	repository  *Repository
	stockClient StockClient
}

func NewService(
	repository *Repository,
	stockClient StockClient,
) *Service {
	return &Service{
		repository:  repository,
		stockClient: stockClient,
	}
}

func (s *Service) Create(
	ctx context.Context,
	input CreateInvoiceInput,
) (*Invoice, error) {
	if len(input.Items) == 0 {
		return nil, ErrInvoiceItemsRequired
	}

	for _, item := range input.Items {
		if item.ProductID <= 0 {
			return nil, ErrInvalidProductID
		}

		if item.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}
	}

	return s.repository.Create(ctx, input)
}

func (s *Service) List(
	ctx context.Context,
) ([]Invoice, error) {
	return s.repository.List(ctx)
}

func (s *Service) FindByID(
	ctx context.Context,
	id int64,
) (*Invoice, error) {
	return s.repository.FindByID(ctx, id)
}
func (r *Repository) Close(
	ctx context.Context,
	id int64,
) (*Invoice, error) {
	var invoice Invoice

	err := r.db.QueryRow(
		ctx,
		`
			UPDATE invoices
			SET
				status = 'CLOSED',
				updated_at = NOW()
			WHERE id = $1
			  AND status = 'OPEN'
			RETURNING id, number, status, created_at, updated_at
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
			"failed to close invoice: %w",
			err,
		)
	}

	return &invoice, nil
}
func (s *Service) Close(
	ctx context.Context,
	id int64,
) (*Invoice, error) {
	invoice, err := s.repository.FindByID(
		ctx,
		id,
	)
	if err != nil {
		return nil, err
	}

	if invoice.Status == StatusClosed {
		return nil, ErrInvoiceAlreadyClosed
	}

	for _, item := range invoice.Items {
		err := s.stockClient.DecreaseStock(
			ctx,
			item.ProductID,
			item.Quantity,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"%w: %v",
				ErrStockUpdateFailed,
				err,
			)
		}
	}

	closedInvoice, err := s.repository.Close(
		ctx,
		id,
	)
	if err != nil {
		return nil, err
	}

	closedInvoice.Items = invoice.Items

	return closedInvoice, nil
}
