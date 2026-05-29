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

// SyntheticVariable represents a global synthetic test variable.
type SyntheticVariable struct {
	ID             int    `json:"id"`
	OrganizationID int    `json:"organizationId,omitempty"`
	Name           string `json:"name"`
	Value          string `json:"value"`
	IsSecret       bool   `json:"isSecret"`
	CreatedAt      int64  `json:"createdAt,omitempty"`
	UpdatedAt      int64  `json:"updatedAt,omitempty"`
}

// SyntheticVariableRequest is the request body for creating or updating a synthetic variable.
type SyntheticVariableRequest struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	IsSecret bool   `json:"isSecret"`
}

// ListSyntheticVariables retrieves all synthetic variables.
func (c *Client) ListSyntheticVariables() ([]SyntheticVariable, error) {
	var variables []SyntheticVariable
	err := c.doRequest(http.MethodGet, "/v1/synthetics/variables", nil, &variables)
	if err != nil {
		return nil, err
	}
	return variables, nil
}

// GetSyntheticVariable retrieves a synthetic variable by ID.
func (c *Client) GetSyntheticVariable(id int) (*SyntheticVariable, error) {
	variables, err := c.ListSyntheticVariables()
	if err != nil {
		return nil, err
	}
	for _, variable := range variables {
		if variable.ID == id {
			return &variable, nil
		}
	}
	return nil, notFoundError("Synthetic variable not found")
}

// CreateSyntheticVariable creates a synthetic variable.
func (c *Client) CreateSyntheticVariable(req SyntheticVariableRequest) (*SyntheticVariable, error) {
	var variable SyntheticVariable
	err := c.doRequest(http.MethodPost, "/v1/synthetics/variables", req, &variable)
	if err != nil {
		return nil, err
	}
	return &variable, nil
}

// UpdateSyntheticVariable updates a synthetic variable.
func (c *Client) UpdateSyntheticVariable(id int, req SyntheticVariableRequest) (*SyntheticVariable, error) {
	var variable SyntheticVariable
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/synthetics/variables/%d", id), req, &variable)
	if err != nil {
		return nil, err
	}
	return &variable, nil
}

// DeleteSyntheticVariable deletes a synthetic variable.
func (c *Client) DeleteSyntheticVariable(id int) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/synthetics/variables/%d", id), nil, nil)
}
