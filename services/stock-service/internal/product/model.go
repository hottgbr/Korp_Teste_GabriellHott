package product

import "time"

type Product struct {
	ID          int64     `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateProductInput struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Stock       int    `json:"stock"`
}

type DecreaseStockInput struct {
	Quantity int `json:"quantity"`
}

type DecreaseStockBatchItemInput struct {
	ProductID int64 `json:"productId"`
	Quantity  int   `json:"quantity"`
}

type DecreaseStockBatchInput struct {
	Items []DecreaseStockBatchItemInput `json:"items"`
}
