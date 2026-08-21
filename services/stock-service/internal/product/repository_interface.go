package product

import "context"

type ProductRepository interface {
	Create(
		ctx context.Context,
		input CreateProductInput,
	) (*Product, error)

	List(ctx context.Context) ([]Product, error)

	FindByID(
		ctx context.Context,
		id int64,
	) (*Product, error)

	DecreaseStock(
		ctx context.Context,
		id int64,
		quantity int,
	) (*Product, error)
}
