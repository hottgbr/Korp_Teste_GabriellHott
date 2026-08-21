package invoice

import (
	"context"
	"errors"
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
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
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
