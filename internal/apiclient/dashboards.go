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

// ListDashboards retrieves all dashboards.
func (c *Client) ListDashboards() ([]Dashboard, error) {
	var dashboards []Dashboard
	err := c.doRequest(http.MethodGet, "/v1/dashboards", nil, &dashboards)
	if err != nil {
		return nil, err
	}
	return dashboards, nil
}

// DeleteDashboard deletes a dashboard by ID.
func (c *Client) DeleteDashboard(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/dashboards/%s", id), nil, nil)
}

// DashboardAlert represents a Moneat dashboard alert.
type DashboardAlert struct {
	ID               string  `json:"id"`
	DashboardID      string  `json:"dashboardId"`
	WidgetID         string  `json:"widgetId"`
	Name             string  `json:"name"`
	Condition        string  `json:"condition"`
	Threshold        float64 `json:"threshold"`
	DurationSeconds  int64   `json:"durationSeconds"`
	IncidentSeverity string  `json:"incidentSeverity"`
}

// CreateDashboardAlertRequest is the request body for creating a dashboard alert.
type CreateDashboardAlertRequest struct {
	WidgetID         string  `json:"widgetId"`
	Name             string  `json:"name"`
	Condition        string  `json:"condition"`
	Threshold        float64 `json:"threshold"`
	DurationSeconds  int64   `json:"durationSeconds"`
	IncidentSeverity string  `json:"incidentSeverity"`
}

// UpdateDashboardAlertRequest is the request body for updating a dashboard alert.
type UpdateDashboardAlertRequest struct {
	WidgetID         string  `json:"widgetId"`
	Name             string  `json:"name"`
	Condition        string  `json:"condition"`
	Threshold        float64 `json:"threshold"`
	DurationSeconds  int64   `json:"durationSeconds"`
	IncidentSeverity string  `json:"incidentSeverity"`
}

// GetDashboardAlert retrieves a dashboard alert by ID.
func (c *Client) GetDashboardAlert(dashboardID, alertID string) (*DashboardAlert, error) {
	var alert DashboardAlert
	path := fmt.Sprintf("/v1/dashboards/%s/alerts/%s", dashboardID, alertID)
	err := c.doRequest(http.MethodGet, path, nil, &alert)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

// ListDashboardAlerts retrieves all alerts for a dashboard.
func (c *Client) ListDashboardAlerts(dashboardID string) ([]DashboardAlert, error) {
	var alerts []DashboardAlert
	path := fmt.Sprintf("/v1/dashboards/%s/alerts", dashboardID)
	err := c.doRequest(http.MethodGet, path, nil, &alerts)
	if err != nil {
		return nil, err
	}
	return alerts, nil
}

// CreateDashboardAlert creates a new dashboard alert.
func (c *Client) CreateDashboardAlert(dashboardID string, req CreateDashboardAlertRequest) (*DashboardAlert, error) {
	var alert DashboardAlert
	path := fmt.Sprintf("/v1/dashboards/%s/alerts", dashboardID)
	err := c.doRequest(http.MethodPost, path, req, &alert)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

// UpdateDashboardAlert updates an existing dashboard alert.
func (c *Client) UpdateDashboardAlert(dashboardID, alertID string, req UpdateDashboardAlertRequest) (*DashboardAlert, error) {
	var alert DashboardAlert
	path := fmt.Sprintf("/v1/dashboards/%s/alerts/%s", dashboardID, alertID)
	err := c.doRequest(http.MethodPut, path, req, &alert)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

// DeleteDashboardAlert deletes a dashboard alert by ID.
func (c *Client) DeleteDashboardAlert(dashboardID, alertID string) error {
	path := fmt.Sprintf("/v1/dashboards/%s/alerts/%s", dashboardID, alertID)
	return c.doRequest(http.MethodDelete, path, nil, nil)
}

// DashboardFolder represents a Moneat dashboard folder.
type DashboardFolder struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt,omitempty"`
}

// CreateDashboardFolderRequest is the request body for creating a dashboard folder.
type CreateDashboardFolderRequest struct {
	Name string `json:"name"`
}

// UpdateDashboardFolderRequest is the request body for updating a dashboard folder.
type UpdateDashboardFolderRequest struct {
	Name string `json:"name"`
}

// GetDashboardFolder retrieves a dashboard folder by ID.
func (c *Client) GetDashboardFolder(id string) (*DashboardFolder, error) {
	var folder DashboardFolder
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/dashboards/folders/%s", id), nil, &folder)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

// ListDashboardFolders retrieves all dashboard folders.
func (c *Client) ListDashboardFolders() ([]DashboardFolder, error) {
	var folders []DashboardFolder
	err := c.doRequest(http.MethodGet, "/v1/dashboards/folders", nil, &folders)
	if err != nil {
		return nil, err
	}
	return folders, nil
}

// CreateDashboardFolder creates a new dashboard folder.
func (c *Client) CreateDashboardFolder(req CreateDashboardFolderRequest) (*DashboardFolder, error) {
	var folder DashboardFolder
	err := c.doRequest(http.MethodPost, "/v1/dashboards/folders", req, &folder)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

// UpdateDashboardFolder updates an existing dashboard folder.
func (c *Client) UpdateDashboardFolder(id string, req UpdateDashboardFolderRequest) (*DashboardFolder, error) {
	var folder DashboardFolder
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/dashboards/folders/%s", id), req, &folder)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

// DeleteDashboardFolder deletes a dashboard folder by ID.
func (c *Client) DeleteDashboardFolder(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/dashboards/folders/%s", id), nil, nil)
}
