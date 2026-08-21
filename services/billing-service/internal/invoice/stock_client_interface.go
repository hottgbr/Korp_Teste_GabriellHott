package invoice

import "context"

type StockItem struct {
	ProductID int64
	Quantity  int
}

type StockClient interface {
	DecreaseStockBatch(
		ctx context.Context,
		items []StockItem,
	) error
}