package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Project represents a Moneat project.
type Project struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Platform  string       `json:"platform"`
	Framework string       `json:"framework,omitempty"`
	Slug      string       `json:"slug,omitempty"`
	Keys      []ProjectKey `json:"keys,omitempty"`
	DSN       string       `json:"dsn,omitempty"`
	CreatedAt string       `json:"createdAt,omitempty"`
}

// ProjectKey represents a Sentry-compatible DSN target for a project.
type ProjectKey struct {
	PlatformTarget string `json:"platformTarget"`
	DSN            string `json:"dsn"`
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

// AddProjectTargetRequest is the request body for adding a project target.
type AddProjectTargetRequest struct {
	Target string `json:"target"`
}

// UnmarshalJSON supports both the older provider-facing ID shape and the
// current API shape, where resourceId is the stable external project ID.
func (p *Project) UnmarshalJSON(data []byte) error {
	var raw struct {
		ID         json.RawMessage `json:"id"`
		ResourceID string          `json:"resourceId"`
		Name       string          `json:"name"`
		Platform   string          `json:"platform"`
		Framework  string          `json:"framework"`
		Slug       string          `json:"slug"`
		Keys       []ProjectKey    `json:"keys"`
		DSN        string          `json:"dsn"`
		CreatedAt  string          `json:"createdAt"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	p.ID = raw.ResourceID
	if p.ID == "" {
		p.ID = rawIDToString(raw.ID)
	}
	p.Name = raw.Name
	p.Platform = raw.Platform
	if p.Platform == "" {
		p.Platform = raw.Framework
	}
	p.Framework = raw.Framework
	p.Slug = raw.Slug
	p.Keys = raw.Keys
	p.DSN = raw.DSN
	p.CreatedAt = raw.CreatedAt
	return nil
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

// AddProjectTarget creates a new DSN target for a project.
func (c *Client) AddProjectTarget(projectID string, req AddProjectTargetRequest) (*ProjectKey, error) {
	var key ProjectKey
	err := c.doRequest(http.MethodPost, fmt.Sprintf("/v1/projects/%s/targets", projectID), req, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// GetProjectTarget finds a project target from the project details response.
func (c *Client) GetProjectTarget(projectID, target string) (*ProjectKey, error) {
	project, err := c.GetProject(projectID)
	if err != nil {
		return nil, err
	}
	for _, key := range project.Keys {
		if key.PlatformTarget == target {
			return &key, nil
		}
	}
	return nil, notFoundError("Project target not found")
}
