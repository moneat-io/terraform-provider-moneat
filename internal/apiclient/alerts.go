package apiclient

import (
	"fmt"
	"net/http"
)

// SystemAlert represents a Moneat system alert.
type SystemAlert struct {
	ID              string  `json:"id"`
	SystemID        string  `json:"systemId"`
	Metric          string  `json:"metric"`
	Condition       string  `json:"condition"`
	Threshold       float64 `json:"threshold"`
	DurationSeconds int64   `json:"durationSeconds"`
	Enabled         bool    `json:"enabled"`
	CreatedAt       string  `json:"createdAt,omitempty"`
}

// CreateSystemAlertRequest is the request body for creating a system alert.
type CreateSystemAlertRequest struct {
	Metric          string  `json:"metric"`
	Condition       string  `json:"condition"`
	Threshold       float64 `json:"threshold"`
	DurationSeconds int64   `json:"durationSeconds"`
	Enabled         bool    `json:"enabled"`
}

// UpdateSystemAlertRequest is the request body for updating a system alert.
type UpdateSystemAlertRequest struct {
	Metric          string  `json:"metric"`
	Condition       string  `json:"condition"`
	Threshold       float64 `json:"threshold"`
	DurationSeconds int64   `json:"durationSeconds"`
	Enabled         bool    `json:"enabled"`
}

// GetSystemAlert retrieves a system alert by ID.
func (c *Client) GetSystemAlert(systemID, alertID string) (*SystemAlert, error) {
	var alert SystemAlert
	path := fmt.Sprintf("/v1/monitor/systems/%s/alerts/%s", systemID, alertID)
	err := c.doRequest(http.MethodGet, path, nil, &alert)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

// CreateSystemAlert creates a new system alert.
func (c *Client) CreateSystemAlert(systemID string, req CreateSystemAlertRequest) (*SystemAlert, error) {
	var alert SystemAlert
	path := fmt.Sprintf("/v1/monitor/systems/%s/alerts", systemID)
	err := c.doRequest(http.MethodPost, path, req, &alert)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

// UpdateSystemAlert updates an existing system alert.
func (c *Client) UpdateSystemAlert(systemID, alertID string, req UpdateSystemAlertRequest) (*SystemAlert, error) {
	var alert SystemAlert
	path := fmt.Sprintf("/v1/monitor/systems/%s/alerts/%s", systemID, alertID)
	err := c.doRequest(http.MethodPut, path, req, &alert)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

// DeleteSystemAlert deletes a system alert by ID.
func (c *Client) DeleteSystemAlert(systemID, alertID string) error {
	path := fmt.Sprintf("/v1/monitor/systems/%s/alerts/%s", systemID, alertID)
	return c.doRequest(http.MethodDelete, path, nil, nil)
}
