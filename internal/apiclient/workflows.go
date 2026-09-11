package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Workflow is the UUID-addressed workflow contract exposed by the Moneat API.
// JSON fields remain raw messages so Terraform can preserve the graph and
// typed workflow configuration without inventing a second provider model.
type Workflow struct {
	ID                string          `json:"id"`
	Name              string          `json:"name"`
	TriggerName       string          `json:"trigger_name"`
	Enabled           bool            `json:"enabled"`
	Version           int             `json:"version"`
	Published         bool            `json:"published"`
	SystemKey         string          `json:"system_key,omitempty"`
	Conditions        json.RawMessage `json:"conditions"`
	Steps             json.RawMessage `json:"steps"`
	Graph             json.RawMessage `json:"graph"`
	OnceForTemplate   []string        `json:"once_for_template"`
	RunOnce           json.RawMessage `json:"run_once,omitempty"`
	InputSchema       json.RawMessage `json:"input_schema,omitempty"`
	TriggerNames      []string        `json:"trigger_names,omitempty"`
	Schedules         json.RawMessage `json:"schedules,omitempty"`
	OwnerUserID       string          `json:"owner_user_id,omitempty"`
	ExecutionIdentity json.RawMessage `json:"execution_identity,omitempty"`
	CreatedAt         string          `json:"created_at,omitempty"`
	UpdatedAt         string          `json:"updated_at,omitempty"`
	LastRunAt         string          `json:"last_run_at,omitempty"`
	RunCount          int64           `json:"run_count,omitempty"`
}

// CreateWorkflowRequest is the complete declarative workflow payload.
type CreateWorkflowRequest struct {
	Name              string          `json:"name"`
	TriggerName       string          `json:"trigger_name"`
	Enabled           bool            `json:"enabled"`
	Conditions        json.RawMessage `json:"conditions"`
	Steps             json.RawMessage `json:"steps"`
	Graph             json.RawMessage `json:"graph,omitempty"`
	OnceForTemplate   []string        `json:"once_for_template"`
	RunOnce           json.RawMessage `json:"run_once,omitempty"`
	InputSchema       json.RawMessage `json:"input_schema,omitempty"`
	TriggerNames      []string        `json:"trigger_names,omitempty"`
	Schedules         json.RawMessage `json:"schedules,omitempty"`
	ExecutionIdentity json.RawMessage `json:"execution_identity,omitempty"`
}

// UpdateWorkflowRequest is the complete mutable workflow payload. Expected
// version is sent on every Terraform update to prevent lost writes.
type UpdateWorkflowRequest struct {
	Name              string          `json:"name,omitempty"`
	Enabled           *bool           `json:"enabled,omitempty"`
	Conditions        json.RawMessage `json:"conditions,omitempty"`
	Steps             json.RawMessage `json:"steps,omitempty"`
	Graph             json.RawMessage `json:"graph,omitempty"`
	OnceForTemplate   []string        `json:"once_for_template,omitempty"`
	ExpectedVersion   *int            `json:"expected_version,omitempty"`
	RunOnce           json.RawMessage `json:"run_once,omitempty"`
	InputSchema       json.RawMessage `json:"input_schema,omitempty"`
	TriggerNames      []string        `json:"trigger_names,omitempty"`
	Schedules         json.RawMessage `json:"schedules,omitempty"`
	ExecutionIdentity json.RawMessage `json:"execution_identity,omitempty"`
}

func workflowPath(id string) string {
	return "/v1/workflows/" + url.PathEscape(id)
}

// GetWorkflow retrieves a workflow by UUID resource ID.
func (c *Client) GetWorkflow(id string) (*Workflow, error) {
	var workflow Workflow
	if err := c.doRequest(http.MethodGet, workflowPath(id), nil, &workflow); err != nil {
		return nil, err
	}
	return &workflow, nil
}

// ListWorkflows retrieves all workflows in the current organization.
func (c *Client) ListWorkflows() ([]Workflow, error) {
	var workflows []Workflow
	if err := c.doRequest(http.MethodGet, "/v1/workflows", nil, &workflows); err != nil {
		return nil, err
	}
	return workflows, nil
}

// CreateWorkflow creates a workflow and returns its UUID resource ID.
func (c *Client) CreateWorkflow(req CreateWorkflowRequest) (*Workflow, error) {
	var workflow Workflow
	if err := c.doRequest(http.MethodPost, "/v1/workflows", req, &workflow); err != nil {
		return nil, err
	}
	return &workflow, nil
}

// UpdateWorkflow updates a workflow using optimistic versioning.
func (c *Client) UpdateWorkflow(id string, req UpdateWorkflowRequest) (*Workflow, error) {
	var workflow Workflow
	if err := c.doRequest(http.MethodPut, workflowPath(id), req, &workflow); err != nil {
		return nil, err
	}
	return &workflow, nil
}

// DeleteWorkflow deletes a workflow by UUID resource ID.
func (c *Client) DeleteWorkflow(id string) error {
	return c.doRequest(http.MethodDelete, workflowPath(id), nil, nil)
}

func lifecyclePath(id, action string, expectedVersion *int) string {
	path := workflowPath(id) + "/" + action
	if expectedVersion != nil {
		path += fmt.Sprintf("?expected_version=%d", *expectedVersion)
	}
	return path
}

// PublishWorkflow publishes the latest validated workflow version.
func (c *Client) PublishWorkflow(id string, expectedVersion *int) (*Workflow, error) {
	var workflow Workflow
	if err := c.doRequest(http.MethodPost, lifecyclePath(id, "publish", expectedVersion), nil, &workflow); err != nil {
		return nil, err
	}
	return &workflow, nil
}

// UnpublishWorkflow removes publication from the latest workflow version.
func (c *Client) UnpublishWorkflow(id string, expectedVersion *int) (*Workflow, error) {
	var workflow Workflow
	if err := c.doRequest(http.MethodPost, lifecyclePath(id, "unpublish", expectedVersion), nil, &workflow); err != nil {
		return nil, err
	}
	return &workflow, nil
}

// GetWorkflowCatalog returns the server-owned workflow action catalog.
func (c *Client) GetWorkflowCatalog() (json.RawMessage, error) {
	var value json.RawMessage
	if err := c.doRequest(http.MethodGet, "/v1/workflows/catalog", nil, &value); err != nil {
		return nil, err
	}
	return value, nil
}

// ListWorkflowBlueprints returns the portable workflow blueprints catalog.
func (c *Client) ListWorkflowBlueprints() (json.RawMessage, error) {
	var value json.RawMessage
	if err := c.doRequest(http.MethodGet, "/v1/workflows/blueprints", nil, &value); err != nil {
		return nil, err
	}
	return value, nil
}

// GetWorkflowBlueprint returns one blueprint by catalog key.
func (c *Client) GetWorkflowBlueprint(key string) (json.RawMessage, error) {
	var value json.RawMessage
	path := "/v1/workflows/blueprints/" + url.PathEscape(key)
	if err := c.doRequest(http.MethodGet, path, nil, &value); err != nil {
		return nil, err
	}
	return value, nil
}
