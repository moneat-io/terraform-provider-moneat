package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const defaultTimeout = 30 * time.Second

// Client is the Moneat API client.
type Client struct {
	BaseURL       string
	Token         string
	ResponseToken string
	HTTPClient    *http.Client
}

// NewClient creates a new Moneat API client.
func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

// ResponseConfigurationClient returns a client authenticated with the optional
// response automation key configured by the provider. It falls back to the
// provider's primary token for backwards-compatible single-token setups.
func (c *Client) ResponseConfigurationClient() *Client {
	if c.ResponseToken == "" || c.ResponseToken == c.Token {
		return c
	}
	return c.WithToken(c.ResponseToken)
}

// WithToken returns a client that shares the base URL and HTTP transport while
// using a different bearer token.
func (c *Client) WithToken(token string) *Client {
	return &Client{
		BaseURL:       c.BaseURL,
		Token:         token,
		ResponseToken: c.ResponseToken,
		HTTPClient:    c.HTTPClient,
	}
}

// APIError represents an error response from the Moneat API.
type APIError struct {
	StatusCode int
	Message    string
	Body       string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("moneat API error (HTTP %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("moneat API error (HTTP %d): %s", e.StatusCode, e.Body)
}

func (c *Client) doRequest(method, path string, body interface{}, result interface{}) error {
	return c.doRequestWithHeaders(method, path, body, result, nil)
}

func (c *Client) doRequestWithHeaders(
	method string,
	path string,
	body interface{},
	result interface{},
	extraHeaders map[string]string,
) error {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)

	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshaling request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	for name, value := range extraHeaders {
		req.Header.Set(name, value)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
		var errResp struct {
			Message string `json:"message"`
			Error   string `json:"error"`
		}
		if json.Unmarshal(respBody, &errResp) == nil {
			if errResp.Message != "" {
				apiErr.Message = errResp.Message
			} else if errResp.Error != "" {
				apiErr.Message = errResp.Error
			}
		}
		return apiErr
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("unmarshaling response: %w", err)
		}
	}

	return nil
}

// IsNotFound returns true if the error is a 404 Not Found.
func IsNotFound(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsConflict returns true when the API rejected a write because the resource
// version changed since it was read.
func IsConflict(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.StatusCode == http.StatusConflict
	}
	return false
}

func notFoundError(message string) error {
	return &APIError{
		StatusCode: http.StatusNotFound,
		Message:    message,
	}
}

func rawIDToString(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}

	var text string
	if json.Unmarshal(raw, &text) == nil {
		return text
	}

	var number json.Number
	if json.Unmarshal(raw, &number) == nil {
		return number.String()
	}

	return ""
}
