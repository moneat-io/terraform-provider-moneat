package apiclient

import (
	"fmt"
	"net/http"
)

// CustomDataSource represents a Moneat custom data source.
type CustomDataSource struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	ConnectionString string `json:"connectionString"`
	CreatedAt        string `json:"createdAt,omitempty"`
}

// CreateCustomDataSourceRequest is the request body for creating a custom data source.
type CreateCustomDataSourceRequest struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	ConnectionString string `json:"connectionString"`
}

// UpdateCustomDataSourceRequest is the request body for updating a custom data source.
type UpdateCustomDataSourceRequest struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	ConnectionString string `json:"connectionString"`
}

// GetCustomDataSource retrieves a custom data source by ID.
func (c *Client) GetCustomDataSource(id string) (*CustomDataSource, error) {
	var ds CustomDataSource
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/datasources/%s", id), nil, &ds)
	if err != nil {
		return nil, err
	}
	return &ds, nil
}

// ListCustomDataSources retrieves all custom data sources.
func (c *Client) ListCustomDataSources() ([]CustomDataSource, error) {
	var datasources []CustomDataSource
	err := c.doRequest(http.MethodGet, "/v1/datasources", nil, &datasources)
	if err != nil {
		return nil, err
	}
	return datasources, nil
}

// CreateCustomDataSource creates a new custom data source.
func (c *Client) CreateCustomDataSource(req CreateCustomDataSourceRequest) (*CustomDataSource, error) {
	var ds CustomDataSource
	err := c.doRequest(http.MethodPost, "/v1/datasources", req, &ds)
	if err != nil {
		return nil, err
	}
	return &ds, nil
}

// UpdateCustomDataSource updates an existing custom data source.
func (c *Client) UpdateCustomDataSource(id string, req UpdateCustomDataSourceRequest) (*CustomDataSource, error) {
	var ds CustomDataSource
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/datasources/%s", id), req, &ds)
	if err != nil {
		return nil, err
	}
	return &ds, nil
}

// DeleteCustomDataSource deletes a custom data source by ID.
func (c *Client) DeleteCustomDataSource(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/datasources/%s", id), nil, nil)
}
