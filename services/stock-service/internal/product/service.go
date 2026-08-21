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
