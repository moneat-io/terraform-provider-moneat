package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Workflow represents a Moneat alert or automation workflow.
type Workflow struct {
	ID              int             `json:"id"`
	Name            string          `json:"name"`
	TriggerName     string          `json:"trigger_name"`
	Enabled         bool            `json:"enabled"`
	Version         int             `json:"version"`
	SystemKey       string          `json:"system_key,omitempty"`
	Conditions      json.RawMessage `json:"conditions"`
	Steps           json.RawMessage `json:"steps"`
	OnceForTemplate []string        `json:"once_for_template"`
	CreatedAt       string          `json:"created_at,omitempty"`
	UpdatedAt       string          `json:"updated_at,omitempty"`
	LastRunAt       string          `json:"last_run_at,omitempty"`
	RunCount        int64           `json:"run_count,omitempty"`
}

// CreateWorkflowRequest is the request body for creating a workflow.
type CreateWorkflowRequest struct {
	Name            string          `json:"name"`
	TriggerName     string          `json:"trigger_name"`
	Enabled         bool            `json:"enabled"`
	Conditions      json.RawMessage `json:"conditions"`
	Steps           json.RawMessage `json:"steps"`
	OnceForTemplate []string        `json:"once_for_template"`
}

// UpdateWorkflowRequest is the request body for updating a workflow.
type UpdateWorkflowRequest struct {
	Name            string          `json:"name,omitempty"`
	Enabled         *bool           `json:"enabled,omitempty"`
	Conditions      json.RawMessage `json:"conditions,omitempty"`
	Steps           json.RawMessage `json:"steps,omitempty"`
	OnceForTemplate []string        `json:"once_for_template"`
}

// GetWorkflow retrieves a workflow by ID.
func (c *Client) GetWorkflow(id int) (*Workflow, error) {
	var workflow Workflow
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/workflows/%d", id), nil, &workflow)
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

// ListWorkflows retrieves all workflows.
func (c *Client) ListWorkflows() ([]Workflow, error) {
	var workflows []Workflow
	err := c.doRequest(http.MethodGet, "/v1/workflows", nil, &workflows)
	if err != nil {
		return nil, err
	}
	return workflows, nil
}

// CreateWorkflow creates a workflow.
func (c *Client) CreateWorkflow(req CreateWorkflowRequest) (*Workflow, error) {
	var workflow Workflow
	err := c.doRequest(http.MethodPost, "/v1/workflows", req, &workflow)
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

// UpdateWorkflow updates a workflow.
func (c *Client) UpdateWorkflow(id int, req UpdateWorkflowRequest) (*Workflow, error) {
	var workflow Workflow
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/workflows/%d", id), req, &workflow)
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

// DeleteWorkflow deletes a workflow.
func (c *Client) DeleteWorkflow(id int) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/workflows/%d", id), nil, nil)
}
