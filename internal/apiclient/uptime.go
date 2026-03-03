package apiclient

import (
	"fmt"
	"net/http"
)

// UptimeMonitor represents a Moneat uptime monitor.
type UptimeMonitor struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	URL             string `json:"url"`
	Type            string `json:"type"`
	IntervalSeconds int64  `json:"intervalSeconds"`
	Paused          bool   `json:"paused"`
	CreatedAt       string `json:"createdAt,omitempty"`
}

// CreateUptimeMonitorRequest is the request body for creating an uptime monitor.
type CreateUptimeMonitorRequest struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	Type            string `json:"type"`
	IntervalSeconds int64  `json:"intervalSeconds"`
}

// UpdateUptimeMonitorRequest is the request body for updating an uptime monitor.
type UpdateUptimeMonitorRequest struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	Type            string `json:"type"`
	IntervalSeconds int64  `json:"intervalSeconds"`
	Paused          bool   `json:"paused"`
}

// GetUptimeMonitor retrieves an uptime monitor by ID.
func (c *Client) GetUptimeMonitor(id string) (*UptimeMonitor, error) {
	var monitor UptimeMonitor
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/uptime/monitors/%s", id), nil, &monitor)
	if err != nil {
		return nil, err
	}
	return &monitor, nil
}

// CreateUptimeMonitor creates a new uptime monitor.
func (c *Client) CreateUptimeMonitor(req CreateUptimeMonitorRequest) (*UptimeMonitor, error) {
	var monitor UptimeMonitor
	err := c.doRequest(http.MethodPost, "/v1/uptime/monitors", req, &monitor)
	if err != nil {
		return nil, err
	}
	return &monitor, nil
}

// UpdateUptimeMonitor updates an existing uptime monitor.
func (c *Client) UpdateUptimeMonitor(id string, req UpdateUptimeMonitorRequest) (*UptimeMonitor, error) {
	var monitor UptimeMonitor
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/uptime/monitors/%s", id), req, &monitor)
	if err != nil {
		return nil, err
	}
	return &monitor, nil
}

// ListUptimeMonitors retrieves all uptime monitors.
func (c *Client) ListUptimeMonitors() ([]UptimeMonitor, error) {
	var monitors []UptimeMonitor
	err := c.doRequest(http.MethodGet, "/v1/uptime/monitors", nil, &monitors)
	if err != nil {
		return nil, err
	}
	return monitors, nil
}

// DeleteUptimeMonitor deletes an uptime monitor by ID.
func (c *Client) DeleteUptimeMonitor(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/uptime/monitors/%s", id), nil, nil)
}
