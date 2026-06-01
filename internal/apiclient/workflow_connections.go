package apiclient

import (
	"fmt"
	"net/http"
)

// WorkflowConnection represents a vaulted workflow connection.
type WorkflowConnection struct {
	ID             int               `json:"id"`
	Type           string            `json:"type"`
	Name           string            `json:"name"`
	IdentifierTags map[string]string `json:"identifier_tags"`
	LastFour       string            `json:"last_four,omitempty"`
	CreatedAt      string            `json:"created_at,omitempty"`
	UpdatedAt      string            `json:"updated_at,omitempty"`
}

// CreateWorkflowConnectionRequest is the request body for creating a workflow connection.
type CreateWorkflowConnectionRequest struct {
	Type           string            `json:"type"`
	Name           string            `json:"name"`
	IdentifierTags map[string]string `json:"identifier_tags"`
	Secret         string            `json:"secret"`
}

// RotateWorkflowConnectionRequest is the request body for rotating a workflow connection secret.
type RotateWorkflowConnectionRequest struct {
	Secret string `json:"secret"`
}

// WorkflowConnectionGroup represents a workflow connection group.
type WorkflowConnectionGroup struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	ConnectionType      string `json:"connection_type"`
	MemberConnectionIDs []int  `json:"member_connection_ids"`
	SelectionStrategy   string `json:"selection_strategy"`
	CreatedAt           string `json:"created_at,omitempty"`
	UpdatedAt           string `json:"updated_at,omitempty"`
}

// CreateWorkflowConnectionGroupRequest is the request body for creating a connection group.
type CreateWorkflowConnectionGroupRequest struct {
	Name                string `json:"name"`
	ConnectionType      string `json:"connection_type"`
	MemberConnectionIDs []int  `json:"member_connection_ids"`
	SelectionStrategy   string `json:"selection_strategy"`
}

// ListWorkflowConnections retrieves workflow connections.
func (c *Client) ListWorkflowConnections() ([]WorkflowConnection, error) {
	var connections []WorkflowConnection
	err := c.doRequest(http.MethodGet, "/v1/workflows/connections", nil, &connections)
	if err != nil {
		return nil, err
	}
	return connections, nil
}

// GetWorkflowConnection retrieves a workflow connection by ID.
func (c *Client) GetWorkflowConnection(id int) (*WorkflowConnection, error) {
	var connection WorkflowConnection
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/workflows/connections/%d", id), nil, &connection)
	if err != nil {
		return nil, err
	}
	return &connection, nil
}

// CreateWorkflowConnection creates a workflow connection.
func (c *Client) CreateWorkflowConnection(req CreateWorkflowConnectionRequest) (*WorkflowConnection, error) {
	var connection WorkflowConnection
	err := c.doRequest(http.MethodPost, "/v1/workflows/connections", req, &connection)
	if err != nil {
		return nil, err
	}
	return &connection, nil
}

// RotateWorkflowConnection rotates a workflow connection secret.
func (c *Client) RotateWorkflowConnection(id int, req RotateWorkflowConnectionRequest) (*WorkflowConnection, error) {
	var connection WorkflowConnection
	path := fmt.Sprintf("/v1/workflows/connections/%d/rotate", id)
	err := c.doRequest(http.MethodPut, path, req, &connection)
	if err != nil {
		return nil, err
	}
	return &connection, nil
}

// DeleteWorkflowConnection deletes a workflow connection.
func (c *Client) DeleteWorkflowConnection(id int) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/workflows/connections/%d", id), nil, nil)
}

// ListWorkflowConnectionGroups retrieves workflow connection groups.
func (c *Client) ListWorkflowConnectionGroups() ([]WorkflowConnectionGroup, error) {
	var groups []WorkflowConnectionGroup
	err := c.doRequest(http.MethodGet, "/v1/workflows/connection-groups", nil, &groups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// GetWorkflowConnectionGroup retrieves a workflow connection group by ID.
func (c *Client) GetWorkflowConnectionGroup(id int) (*WorkflowConnectionGroup, error) {
	groups, err := c.ListWorkflowConnectionGroups()
	if err != nil {
		return nil, err
	}
	for _, group := range groups {
		if group.ID == id {
			return &group, nil
		}
	}
	return nil, notFoundError("workflow connection group not found")
}

// CreateWorkflowConnectionGroup creates a workflow connection group.
func (c *Client) CreateWorkflowConnectionGroup(
	req CreateWorkflowConnectionGroupRequest,
) (*WorkflowConnectionGroup, error) {
	var group WorkflowConnectionGroup
	err := c.doRequest(http.MethodPost, "/v1/workflows/connection-groups", req, &group)
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// DeleteWorkflowConnectionGroup deletes a workflow connection group.
func (c *Client) DeleteWorkflowConnectionGroup(id int) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/workflows/connection-groups/%d", id), nil, nil)
}
