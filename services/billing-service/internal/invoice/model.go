package invoice

import "time"

type Status string

const (
	StatusOpen   Status = "OPEN"
	StatusClosed Status = "CLOSED"
)

type Invoice struct {
	ID        int64         `json:"id"`
	Number    int64         `json:"number"`
	Status    Status        `json:"status"`
	Items     []InvoiceItem `json:"items"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

type InvoiceItem struct {
	ID        int64 `json:"id"`
	InvoiceID int64 `json:"invoiceId"`
	ProductID int64 `json:"productId"`
	Quantity  int   `json:"quantity"`
}

type CreateInvoiceItemInput struct {
	ProductID int64 `json:"productId"`
	Quantity  int   `json:"quantity"`
}

type CreateInvoiceInput struct {
	Items []CreateInvoiceItemInput `json:"items"`
}