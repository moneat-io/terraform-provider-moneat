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
	path := fmt.Sprintf("/v1/monitor/hosts/%s/alerts/%s", systemID, alertID)
	err := c.doRequest(http.MethodGet, path, nil, &alert)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

// CreateSystemAlert creates a new system alert.
func (c *Client) CreateSystemAlert(systemID string, req CreateSystemAlertRequest) (*SystemAlert, error) {
	var alert SystemAlert
	path := fmt.Sprintf("/v1/monitor/hosts/%s/alerts", systemID)
	err := c.doRequest(http.MethodPost, path, req, &alert)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

// UpdateSystemAlert updates an existing system alert.
func (c *Client) UpdateSystemAlert(systemID, alertID string, req UpdateSystemAlertRequest) (*SystemAlert, error) {
	var alert SystemAlert
	path := fmt.Sprintf("/v1/monitor/hosts/%s/alerts/%s", systemID, alertID)
	err := c.doRequest(http.MethodPut, path, req, &alert)
	if err != nil {
		return nil, err
	}
	return &alert, nil
}

// DeleteSystemAlert deletes a system alert by ID.
func (c *Client) DeleteSystemAlert(systemID, alertID string) error {
	path := fmt.Sprintf("/v1/monitor/hosts/%s/alerts/%s", systemID, alertID)
	return c.doRequest(http.MethodDelete, path, nil, nil)
}

// SilencePeriod represents a Moneat alert silence period.
type SilencePeriod struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Reason    string `json:"reason,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}

// CreateSilencePeriodRequest is the request body for creating a silence period.
type CreateSilencePeriodRequest struct {
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Reason    string `json:"reason,omitempty"`
}

// UpdateSilencePeriodRequest is the request body for updating a silence period.
type UpdateSilencePeriodRequest struct {
	Name      string `json:"name"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Reason    string `json:"reason,omitempty"`
}

// GetSilencePeriod retrieves a silence period by ID.
func (c *Client) GetSilencePeriod(id string) (*SilencePeriod, error) {
	var period SilencePeriod
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/monitor/silence-periods/%s", id), nil, &period)
	if err != nil {
		return nil, err
	}
	return &period, nil
}

// ListSilencePeriods retrieves all silence periods.
func (c *Client) ListSilencePeriods() ([]SilencePeriod, error) {
	var periods []SilencePeriod
	err := c.doRequest(http.MethodGet, "/v1/monitor/silence-periods", nil, &periods)
	if err != nil {
		return nil, err
	}
	return periods, nil
}

// CreateSilencePeriod creates a new silence period.
func (c *Client) CreateSilencePeriod(req CreateSilencePeriodRequest) (*SilencePeriod, error) {
	var period SilencePeriod
	err := c.doRequest(http.MethodPost, "/v1/monitor/silence-periods", req, &period)
	if err != nil {
		return nil, err
	}
	return &period, nil
}

// UpdateSilencePeriod updates an existing silence period.
func (c *Client) UpdateSilencePeriod(id string, req UpdateSilencePeriodRequest) (*SilencePeriod, error) {
	var period SilencePeriod
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/monitor/silence-periods/%s", id), req, &period)
	if err != nil {
		return nil, err
	}
	return &period, nil
}

// DeleteSilencePeriod deletes a silence period by ID.
func (c *Client) DeleteSilencePeriod(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/monitor/silence-periods/%s", id), nil, nil)
}
