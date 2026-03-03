package apiclient

import (
	"fmt"
	"net/http"
)

// Project represents a Moneat project.
type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Platform  string `json:"platform"`
	Framework string `json:"framework,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}

// CreateProjectRequest is the request body for creating a project.
type CreateProjectRequest struct {
	Name      string `json:"name"`
	Platform  string `json:"platform"`
	Framework string `json:"framework,omitempty"`
}

// UpdateProjectRequest is the request body for updating a project.
type UpdateProjectRequest struct {
	Name      string `json:"name"`
	Platform  string `json:"platform"`
	Framework string `json:"framework,omitempty"`
}

// GetProject retrieves a project by ID.
func (c *Client) GetProject(id string) (*Project, error) {
	var project Project
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/projects/%s", id), nil, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// ListProjects retrieves all projects.
func (c *Client) ListProjects() ([]Project, error) {
	var projects []Project
	err := c.doRequest(http.MethodGet, "/v1/projects", nil, &projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

// CreateProject creates a new project.
func (c *Client) CreateProject(req CreateProjectRequest) (*Project, error) {
	var project Project
	err := c.doRequest(http.MethodPost, "/v1/projects", req, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// UpdateProject updates an existing project.
func (c *Client) UpdateProject(id string, req UpdateProjectRequest) (*Project, error) {
	var project Project
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/projects/%s", id), req, &project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

// DeleteProject deletes a project by ID.
func (c *Client) DeleteProject(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/projects/%s", id), nil, nil)
}
