package apiclient

import (
	"fmt"
	"net/http"
)

// SyntheticTest represents a Moneat synthetic test.
type SyntheticTest struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	URL             string `json:"url"`
	IntervalSeconds int64  `json:"intervalSeconds"`
	Enabled         bool   `json:"enabled"`
	CreatedAt       string `json:"createdAt,omitempty"`
}

// CreateSyntheticTestRequest is the request body for creating a synthetic test.
type CreateSyntheticTestRequest struct {
	Name            string `json:"name"`
	Type            string `json:"type"`
	URL             string `json:"url"`
	IntervalSeconds int64  `json:"intervalSeconds"`
	Enabled         bool   `json:"enabled"`
}

// UpdateSyntheticTestRequest is the request body for updating a synthetic test.
type UpdateSyntheticTestRequest struct {
	Name            string `json:"name"`
	Type            string `json:"type"`
	URL             string `json:"url"`
	IntervalSeconds int64  `json:"intervalSeconds"`
	Enabled         bool   `json:"enabled"`
}

// GetSyntheticTest retrieves a synthetic test by ID.
func (c *Client) GetSyntheticTest(id string) (*SyntheticTest, error) {
	var test SyntheticTest
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/synthetics/tests/%s", id), nil, &test)
	if err != nil {
		return nil, err
	}
	return &test, nil
}

// ListSyntheticTests retrieves all synthetic tests.
func (c *Client) ListSyntheticTests() ([]SyntheticTest, error) {
	var tests []SyntheticTest
	err := c.doRequest(http.MethodGet, "/v1/synthetics/tests", nil, &tests)
	if err != nil {
		return nil, err
	}
	return tests, nil
}

// CreateSyntheticTest creates a new synthetic test.
func (c *Client) CreateSyntheticTest(req CreateSyntheticTestRequest) (*SyntheticTest, error) {
	var test SyntheticTest
	err := c.doRequest(http.MethodPost, "/v1/synthetics/tests", req, &test)
	if err != nil {
		return nil, err
	}
	return &test, nil
}

// UpdateSyntheticTest updates an existing synthetic test.
func (c *Client) UpdateSyntheticTest(id string, req UpdateSyntheticTestRequest) (*SyntheticTest, error) {
	var test SyntheticTest
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/synthetics/tests/%s", id), req, &test)
	if err != nil {
		return nil, err
	}
	return &test, nil
}

// DeleteSyntheticTest deletes a synthetic test by ID.
func (c *Client) DeleteSyntheticTest(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/synthetics/tests/%s", id), nil, nil)
}
