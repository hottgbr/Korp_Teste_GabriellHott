package product

import (
	"context"
	"errors"
	"testing"
)

type fakeRepository struct {
	createFn func(
		context.Context,
		CreateProductInput,
	) (*Product, error)
}

func (f *fakeRepository) Create(
	ctx context.Context,
	input CreateProductInput,
) (*Product, error) {
	if f.createFn != nil {
		return f.createFn(ctx, input)
	}

	return nil, nil
}

func (f *fakeRepository) List(
	ctx context.Context,
) ([]Product, error) {
	return nil, nil
}

func (f *fakeRepository) FindByID(
	ctx context.Context,
	id int64,
) (*Product, error) {
	return nil, nil
}

func (f *fakeRepository) DecreaseStock(
	ctx context.Context,
	id int64,
	quantity int,
) (*Product, error) {
	return nil, nil
}

func TestServiceCreateReturnsErrorWhenCodeIsEmpty(t *testing.T) {
	repository := &fakeRepository{}
	service := NewService(repository)

	input := CreateProductInput{
		Code:        "",
		Description: "Notebook",
		Stock:       10,
	}

	product, err := service.Create(
		context.Background(),
		input,
	)

	if !errors.Is(err, ErrCodeRequired) {
		t.Fatalf(
			"expected ErrCodeRequired, got %v",
			err,
		)
	}

	if product != nil {
		t.Fatalf(
			"expected product to be nil, got %+v",
			product,
		)
	}
}

func TestServiceCreateReturnsCreatedProduct(t *testing.T) {
	expectedProduct := &Product{
		ID:          1,
		Code:        "PROD-002",
		Description: "Mouse",
		Stock:       20,
	}

	repository := &fakeRepository{
		createFn: func(
			ctx context.Context,
			input CreateProductInput,
		) (*Product, error) {
			return expectedProduct, nil
		},
	}

	service := NewService(repository)

	input := CreateProductInput{
		Code:        "PROD-002",
		Description: "Mouse",
		Stock:       20,
	}

	product, err := service.Create(
		context.Background(),
		input,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if product != expectedProduct {
		t.Fatalf(
			"expected product %+v, got %+v",
			expectedProduct,
			product,
		)
	}
}

func TestServiceCreateValidations(t *testing.T) {
	tests := []struct {
		name          string
		input         CreateProductInput
		expectedError error
	}{
		{
			name: "returns error when code is empty",
			input: CreateProductInput{
				Code:        "",
				Description: "Notebook",
				Stock:       10,
			},
			expectedError: ErrCodeRequired,
		},
		{
			name: "returns error when description is empty",
			input: CreateProductInput{
				Code:        "PROD-001",
				Description: "",
				Stock:       10,
			},
			expectedError: ErrDescriptionRequired,
		},
		{
			name: "returns error when stock is negative",
			input: CreateProductInput{
				Code:        "PROD-001",
				Description: "Notebook",
				Stock:       -1,
			},
			expectedError: ErrInvalidStock,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repository := &fakeRepository{}
			service := NewService(repository)

			product, err := service.Create(
				context.Background(),
				test.input,
			)

			if !errors.Is(err, test.expectedError) {
				t.Fatalf(
					"expected %v, got %v",
					test.expectedError,
					err,
				)
			}

			if product != nil {
				t.Fatalf(
					"expected product to be nil, got %+v",
					product,
				)
			}
		})
	}
}

func (f *fakeRepository) DecreaseStockBatch(
	ctx context.Context,
	items []DecreaseStockBatchItemInput,
) ([]Product, error) {
	return nil, nil
}
