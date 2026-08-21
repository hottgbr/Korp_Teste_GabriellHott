package invoice

import "context"

type StockClient interface {
	DecreaseStock(
		ctx context.Context,
		productID int64,
		quantity int,
	) error
}
