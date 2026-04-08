package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DefaultTimeout is the default timeout for HTTP requests
const DefaultTimeout = 30 * time.Second

// HTTPClient wraps http.Client to provide simplified methods
type HTTPClient struct {
	Client *http.Client
}

// NewHTTPClient creates a new HTTPClient with a default timeout
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		Client: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// NewHTTPClientWithTimeout creates a new HTTPClient with a specified timeout
func NewHTTPClientWithTimeout(timeout time.Duration) *HTTPClient {
	return &HTTPClient{
		Client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Get performs a GET request
func (c *HTTPClient) Get(ctx context.Context, url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}

	c.addHeaders(req, headers)

	return c.do(req)
}

// Post performs a POST request with a JSON body
func (c *HTTPClient) Post(ctx context.Context, url string, body any, headers map[string]string) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	c.addHeaders(req, headers)

	return c.do(req)
}

// do executes the request and handles the response
func (c *HTTPClient) do(req *http.Request) ([]byte, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request execution failed: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Consider 2xx as success
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}

func (c *HTTPClient) addHeaders(req *http.Request, headers map[string]string) {
	for k, v := range headers {
		req.Header.Set(k, v)
	}
}
