package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type ClientConfig struct {
	BasePath   string
	Timeout    time.Duration
	Retry      int
	RetryDelay time.Duration
}

type Client struct {
	config *ClientConfig
	client *http.Client
}

func NewClient(cfg *ClientConfig) *Client {
	return &Client{
		config: cfg,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// Do executes an HTTP request against the client's base path, retrying on
// transient failures (429, 502, 503, 504). The response body is JSON-decoded into dst.
// headers are merged on top of the default Accept/Content-Type headers; pass nil for none.
func (c *Client) Do(ctx context.Context, method, path string, headers http.Header, body, dst any) error {
	url := c.config.BasePath + path

	attempts := max(c.config.Retry, 1)

	var lastErr error
	for i := range attempts {
		if i > 0 && c.config.RetryDelay > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(c.config.RetryDelay):
			}
		}

		req, err := buildRequest(ctx, method, url, headers, body)
		if err != nil {
			return err
		}

		slog.DebugContext(ctx, "http request", "method", method, "url", url, "attempt", i+1)

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			slog.WarnContext(ctx, "http request failed", "method", method, "url", url, "attempt", i+1, "error", err)
			continue
		}

		status := resp.StatusCode
		lastErr = readResponse(ctx, resp, dst)
		if isTransient(status) {
			slog.WarnContext(ctx, "transient response, retrying", "method", method, "url", url, "status", status, "attempt", i+1)
		} else {
			break
		}
	}

	return lastErr
}

func buildRequest(ctx context.Context, method, url string, headers http.Header, body any) (*http.Request, error) {
	var r io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		slog.DebugContext(ctx, "request body", "body", string(b))
		r = bytes.NewReader(b)
	} else {
		slog.DebugContext(ctx, "request body", "body", nil)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, r)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	if headers != nil {
		req.Header = headers
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func readResponse(ctx context.Context, resp *http.Response, dst any) error {
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body := string(bytes.TrimSpace(b))
		slog.ErrorContext(ctx, "http error response", "status", resp.StatusCode, "body", body)
		return fmt.Errorf("upstream %d: %s", resp.StatusCode, body)
	}

	var bodyLog any = nil
	if len(b) > 0 {
		bodyLog = string(b)
	}
	slog.DebugContext(ctx, "response body", "body", bodyLog)

	if dst == nil {
		return nil
	}
	return json.NewDecoder(bytes.NewReader(b)).Decode(dst)
}

func isTransient(status int) bool {
	return status == http.StatusTooManyRequests ||
		status == http.StatusBadGateway ||
		status == http.StatusServiceUnavailable ||
		status == http.StatusGatewayTimeout
}
