package product

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrCodeRequired        = errors.New("product code is required")
	ErrDescriptionRequired = errors.New("product description is required")
	ErrInvalidStock        = errors.New("product stock cannot be negative")
	ErrInvalidQuantity     = errors.New("quantity must be greater than zero")
	ErrStockItemsRequired  = errors.New("at least one stock item is required")
	ErrInvalidProductID    = errors.New("product id must be greater than zero")
)

type Service struct {
	repository ProductRepository
}

func NewService(repository ProductRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	input CreateProductInput,
) (*Product, error) {
	input.Code = strings.TrimSpace(input.Code)
	input.Description = strings.TrimSpace(input.Description)

	if input.Code == "" {
		return nil, ErrCodeRequired
	}

	if input.Description == "" {
		return nil, ErrDescriptionRequired
	}

	if input.Stock < 0 {
		return nil, ErrInvalidStock
	}

	return s.repository.Create(ctx, input)
}

func (s *Service) List(ctx context.Context) ([]Product, error) {
	return s.repository.List(ctx)
}

func (s *Service) FindByID(
	ctx context.Context,
	id int64,
) (*Product, error) {
	return s.repository.FindByID(ctx, id)
}

func (s *Service) DecreaseStock(
	ctx context.Context,
	id int64,
	input DecreaseStockInput,
) (*Product, error) {
	if input.Quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	return s.repository.DecreaseStock(
		ctx,
		id,
		input.Quantity,
	)
}

func (s *Service) DecreaseStockBatch(
	ctx context.Context,
	input DecreaseStockBatchInput,
) ([]Product, error) {
	if len(input.Items) == 0 {
		return nil, ErrStockItemsRequired
	}

	for _, item := range input.Items {
		if item.ProductID <= 0 {
			return nil, ErrInvalidProductID
		}

		if item.Quantity <= 0 {
			return nil, ErrInvalidQuantity
		}
	}

	return s.repository.DecreaseStockBatch(
		ctx,
		input.Items,
	)
}
