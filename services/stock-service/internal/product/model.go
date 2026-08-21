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
