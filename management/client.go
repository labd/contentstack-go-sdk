package management

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

type Auth struct {
	AuthToken string
}

type ClientConfig struct {
	BaseURL         string
	HTTPClient      *http.Client
	AuthToken       string
	OrganizationUID string
	RateLimit       float64
	RateBurst       int
	MaxRetries      int
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Client struct {
	authToken   string
	baseURL     *url.URL
	httpClient  *http.Client
	rateLimiter *rate.Limiter
	maxRetries  int
}

type ErrorMessage struct {
	ErrorMessage string              `json:"error_message"`
	ErrorCode    int                 `json:"error_code"`
	Errors       map[string][]string `json:"errors"`
}

func (e *ErrorMessage) Error() string {
	return e.ErrorMessage
}

func NewClient(cfg ClientConfig) (*Client, error) {
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("missing BaseURL")
	}

	url, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, err
	}

	// If a custom httpClient is passed use that
	var httpClient *http.Client
	if cfg.HTTPClient != nil {
		httpClient = cfg.HTTPClient
	} else {
		httpClient = &http.Client{}
	}

	rateLimit := cfg.RateLimit
	if rateLimit <= 0 {
		rateLimit = 10.0
	}

	rateBurst := cfg.RateBurst
	if rateBurst <= 0 {
		rateBurst = 10
	}

	rateLimiter := rate.NewLimiter(rate.Limit(rateLimit), rateBurst)

	maxRetries := cfg.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}

	client := &Client{
		baseURL:     url,
		authToken:   cfg.AuthToken,
		httpClient:  httpClient,
		rateLimiter: rateLimiter,
		maxRetries:  maxRetries,
	}

	return client, nil
}

func NewClientWithToken(auth *Auth) *Client {
	rateLimiter := rate.NewLimiter(rate.Limit(10.0), 10)

	return &Client{
		authToken:   auth.AuthToken,
		httpClient:  &http.Client{},
		rateLimiter: rateLimiter,
		maxRetries:  3,
	}
}

func (c *Client) head(ctx context.Context, path string, queryParams url.Values, headers http.Header) (*http.Response, error) {
	return c.execute(ctx, http.MethodHead, path, queryParams, headers, nil)
}

func (c *Client) get(ctx context.Context, path string, queryParams url.Values, headers http.Header) (*http.Response, error) {
	return c.execute(ctx, http.MethodGet, path, queryParams, headers, nil)
}

func (c *Client) post(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	return c.execute(ctx, http.MethodPost, path, queryParams, headers, body)
}

func (c *Client) put(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	return c.execute(ctx, http.MethodPut, path, queryParams, headers, body)
}

func (c *Client) delete(ctx context.Context, path string, queryParams url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	return c.execute(ctx, http.MethodDelete, path, queryParams, headers, body)
}

func (c *Client) createEndpoint(p string) (*url.URL, error) {
	url, err := url.Parse(p)
	if err != nil {
		return nil, err
	}
	return c.baseURL.ResolveReference(url), nil
}

func (c *Client) execute(ctx context.Context, method string, path string, params url.Values, headers http.Header, body io.Reader) (*http.Response, error) {
	return c.executeWithRetry(ctx, method, path, params, headers, body, 0)
}

func (c *Client) executeWithRetry(ctx context.Context, method string, path string, params url.Values, headers http.Header, body io.Reader, attempt int) (*http.Response, error) {
	if c.rateLimiter != nil {
		if err := c.rateLimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("rate limiting wait failed: %w", err)
		}
	}

	endpoint, err := c.createEndpoint(path)
	if err != nil {
		return nil, err
	}

	if params != nil {
		endpoint.RawQuery = params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), body)
	if err != nil {
		return nil, fmt.Errorf("Creating new request: %w", err)
	}

	if headers != nil {
		req.Header = headers
	}
	req.Header.Set("Accept", "application/json; charset=utf-8")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 429 && attempt < c.maxRetries {
		resp.Body.Close()
		waitTime := c.calculateBackoffWait(attempt, resp)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(waitTime):
		}
		return c.executeWithRetry(ctx, method, path, params, headers, body, attempt+1)
	}

	return resp, nil
}

// calculateBackoffWait calculates the wait time for exponential backoff
func (c *Client) calculateBackoffWait(attempt int, resp *http.Response) time.Duration {
	// Check if the server provided a Retry-After header
	if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
		if seconds, err := strconv.Atoi(retryAfter); err == nil {
			waitTime := time.Duration(seconds) * time.Second
			// Cap at 60 seconds maximum
			if waitTime > 60*time.Second {
				waitTime = 60 * time.Second
			}
			return waitTime
		}
	}

	// Simple exponential backoff: 1s, 2s, 4s, 8s, 16s, capped at 30s
	waitTime := time.Duration(1<<uint(attempt)) * time.Second
	if waitTime > 30*time.Second {
		waitTime = 30 * time.Second
	}

	return waitTime
}

func (c *Client) processResponse(r *http.Response, dst interface{}) error {
	content, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	switch r.StatusCode {
	case 200, 201:
		if err = json.Unmarshal(content, &dst); err != nil {
			return err
		}
		return nil
	case 401:
		result := make(map[string]interface{})
		if err = json.Unmarshal(content, &result); err != nil {
			return err
		}
		return fmt.Errorf("Not authorized")
	case 404:
		return &ErrorMessage{
			ErrorMessage: "Resource not found",
			ErrorCode:    404,
		}
	case 422:
		result := ErrorMessage{}
		if err = json.Unmarshal(content, &result); err != nil {

			// Contentstack can return 'invalid' json, so try to process that
			// {"error_message":"xx","error_code":115,"errors":[{}]}
			tmp := struct {
				ErrorMessage string `json:"error_message"`
				ErrorCode    int    `json:"error_code"`
			}{}

			if err = json.Unmarshal(content, &tmp); err != nil {
				return err
			}
			return &ErrorMessage{
				ErrorMessage: tmp.ErrorMessage,
				ErrorCode:    tmp.ErrorCode,
			}
		}
		return &result
	case 429:
		return &ErrorMessage{
			ErrorMessage: "Rate limit exceeded. All retry attempts have been exhausted. Please reduce request frequency or increase rate limiting configuration.",
			ErrorCode:    429,
		}
	default:
		return fmt.Errorf("Unhandled StatusCode: %d", r.StatusCode)
	}
}
