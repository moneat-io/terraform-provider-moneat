package apiclient

import (
	"fmt"
	"net/http"
)

// Dashboard represents a Moneat dashboard.
type Dashboard struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
}

// CreateDashboardRequest is the request body for creating a dashboard.
type CreateDashboardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// UpdateDashboardRequest is the request body for updating a dashboard.
type UpdateDashboardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// GetDashboard retrieves a dashboard by ID.
func (c *Client) GetDashboard(id string) (*Dashboard, error) {
	var dashboard Dashboard
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/dashboards/%s", id), nil, &dashboard)
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

// CreateDashboard creates a new dashboard.
func (c *Client) CreateDashboard(req CreateDashboardRequest) (*Dashboard, error) {
	var dashboard Dashboard
	err := c.doRequest(http.MethodPost, "/v1/dashboards", req, &dashboard)
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

// UpdateDashboard updates an existing dashboard.
func (c *Client) UpdateDashboard(id string, req UpdateDashboardRequest) (*Dashboard, error) {
	var dashboard Dashboard
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/dashboards/%s", id), req, &dashboard)
	if err != nil {
		return nil, err
	}
	return &dashboard, nil
}

// DeleteDashboard deletes a dashboard by ID.
func (c *Client) DeleteDashboard(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/dashboards/%s", id), nil, nil)
}
