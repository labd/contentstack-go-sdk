package management

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Auth struct {
	AuthToken string
}

type ClientConfig struct {
	BaseURL         string
	HTTPClient      *http.Client
	AuthToken       string
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Client struct {
	authToken  string
	baseURL    *url.URL
	httpClient *http.Client
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

	client := &Client{
		baseURL:    url,
		authToken:  cfg.AuthToken,
		httpClient: httpClient,
	}

	return client, nil
}

func NewClientWithToken(auth *Auth) *Client {
	return &Client{}
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

	return resp, nil
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
	case 422:
		result := ErrorMessage{}
		if err = json.Unmarshal(content, &result); err != nil {
			return err
		}
		return &result
	default:
		return fmt.Errorf("Unhandled StatusCode: %d", r.StatusCode)
	}
}
