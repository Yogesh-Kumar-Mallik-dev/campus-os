package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client handles network interaction with the Campus OS Go Logical Backend.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// HealthResponse represents the health payload from backend.
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}

// BLOCK_CLIENT_API_INIT_001
// Purpose: Constructs a new backend API client with default timeout.
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// BLOCK_CLIENT_API_HEALTH_001
// Purpose: Pings backend /healthz endpoint to check connectivity.
func (c *Client) CheckHealth(ctx context.Context) (*HealthResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/healthz", nil)
	if err != nil {
		return nil, fmt.Errorf("BLOCK_CLIENT_API_HEALTH_001: failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("BLOCK_CLIENT_API_HEALTH_001: backend unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("BLOCK_CLIENT_API_HEALTH_001: unexpected status %d", resp.StatusCode)
	}

	var res HealthResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("BLOCK_CLIENT_API_HEALTH_001: failed to decode health response: %w", err)
	}
	return &res, nil
}
