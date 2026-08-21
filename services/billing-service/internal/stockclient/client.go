package stockclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

type decreaseStockRequest struct {
	Quantity int `json:"quantity"`
}

func (c *Client) DecreaseStock(
	ctx context.Context,
	productID int64,
	quantity int,
) error {
	body, err := json.Marshal(
		decreaseStockRequest{
			Quantity: quantity,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"failed to encode stock request: %w",
			err,
		)
	}

	url := fmt.Sprintf(
		"%s/products/%d/stock",
		c.baseURL,
		productID,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPatch,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return fmt.Errorf(
			"failed to create stock request: %w",
			err,
		)
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf(
			"stock service unavailable: %w",
			err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"stock service returned status %d",
			resp.StatusCode,
		)
	}

	return nil
}
