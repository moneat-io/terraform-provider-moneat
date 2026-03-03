package apiclient

import (
	"fmt"
	"net/http"
)

// LogIndex represents a Moneat log index.
type LogIndex struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	FilterQuery   string `json:"filterQuery"`
	RetentionDays int64  `json:"retentionDays"`
	CreatedAt     string `json:"createdAt,omitempty"`
}

// CreateLogIndexRequest is the request body for creating a log index.
type CreateLogIndexRequest struct {
	Name          string `json:"name"`
	FilterQuery   string `json:"filterQuery"`
	RetentionDays int64  `json:"retentionDays"`
}

// UpdateLogIndexRequest is the request body for updating a log index.
type UpdateLogIndexRequest struct {
	Name          string `json:"name"`
	FilterQuery   string `json:"filterQuery"`
	RetentionDays int64  `json:"retentionDays"`
}

// GetLogIndex retrieves a log index by ID.
func (c *Client) GetLogIndex(id string) (*LogIndex, error) {
	var index LogIndex
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/logs/indexes/%s", id), nil, &index)
	if err != nil {
		return nil, err
	}
	return &index, nil
}

// ListLogIndexes retrieves all log indexes.
func (c *Client) ListLogIndexes() ([]LogIndex, error) {
	var indexes []LogIndex
	err := c.doRequest(http.MethodGet, "/v1/logs/indexes", nil, &indexes)
	if err != nil {
		return nil, err
	}
	return indexes, nil
}

// CreateLogIndex creates a new log index.
func (c *Client) CreateLogIndex(req CreateLogIndexRequest) (*LogIndex, error) {
	var index LogIndex
	err := c.doRequest(http.MethodPost, "/v1/logs/indexes", req, &index)
	if err != nil {
		return nil, err
	}
	return &index, nil
}

// UpdateLogIndex updates an existing log index.
func (c *Client) UpdateLogIndex(id string, req UpdateLogIndexRequest) (*LogIndex, error) {
	var index LogIndex
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/logs/indexes/%s", id), req, &index)
	if err != nil {
		return nil, err
	}
	return &index, nil
}

// DeleteLogIndex deletes a log index by ID.
func (c *Client) DeleteLogIndex(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/logs/indexes/%s", id), nil, nil)
}

// LogAPIKey represents a Moneat log API key.
type LogAPIKey struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Key       string `json:"key"`
	CreatedAt string `json:"createdAt,omitempty"`
}

// CreateLogAPIKeyRequest is the request body for creating a log API key.
type CreateLogAPIKeyRequest struct {
	Name string `json:"name"`
}

// UpdateLogAPIKeyRequest is the request body for updating a log API key.
type UpdateLogAPIKeyRequest struct {
	Name string `json:"name"`
}

// GetLogAPIKey retrieves a log API key by ID.
func (c *Client) GetLogAPIKey(id string) (*LogAPIKey, error) {
	var key LogAPIKey
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/logs/api-keys/%s", id), nil, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// ListLogAPIKeys retrieves all log API keys.
func (c *Client) ListLogAPIKeys() ([]LogAPIKey, error) {
	var keys []LogAPIKey
	err := c.doRequest(http.MethodGet, "/v1/logs/api-keys", nil, &keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// CreateLogAPIKey creates a new log API key.
func (c *Client) CreateLogAPIKey(req CreateLogAPIKeyRequest) (*LogAPIKey, error) {
	var key LogAPIKey
	err := c.doRequest(http.MethodPost, "/v1/logs/api-keys", req, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// UpdateLogAPIKey updates an existing log API key.
func (c *Client) UpdateLogAPIKey(id string, req UpdateLogAPIKeyRequest) (*LogAPIKey, error) {
	var key LogAPIKey
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/logs/api-keys/%s", id), req, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// DeleteLogAPIKey deletes a log API key by ID.
func (c *Client) DeleteLogAPIKey(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/logs/api-keys/%s", id), nil, nil)
}
