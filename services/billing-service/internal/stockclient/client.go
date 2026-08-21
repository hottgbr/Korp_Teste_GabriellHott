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

type decreaseStockBatchItem struct {
	ProductID int64 `json:"productId"`
	Quantity  int   `json:"quantity"`
}

type decreaseStockBatchRequest struct {
	Items []decreaseStockBatchItem `json:"items"`
}

func (c *Client) DecreaseStockBatch(
	ctx context.Context,
	items []decreaseStockBatchItem,
) error {
	body, err := json.Marshal(
		decreaseStockBatchRequest{
			Items: items,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"failed to encode stock request: %w",
			err,
		)
	}

	url := fmt.Sprintf(
		"%s/products/stock",
		c.baseURL,
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