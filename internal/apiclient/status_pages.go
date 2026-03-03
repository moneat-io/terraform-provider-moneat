package apiclient

import (
	"fmt"
	"net/http"
)

// StatusPage represents a Moneat status page.
type StatusPage struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	IsPublic    bool   `json:"isPublic"`
	CreatedAt   string `json:"createdAt,omitempty"`
}

// CreateStatusPageRequest is the request body for creating a status page.
type CreateStatusPageRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	IsPublic    bool   `json:"isPublic"`
}

// UpdateStatusPageRequest is the request body for updating a status page.
type UpdateStatusPageRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description,omitempty"`
	IsPublic    bool   `json:"isPublic"`
}

// StatusPageCustomDomain represents a custom domain on a status page.
type StatusPageCustomDomain struct {
	ID       string `json:"id"`
	PageID   string `json:"pageId"`
	Domain   string `json:"domain"`
	Verified bool   `json:"verified"`
}

// CreateStatusPageCustomDomainRequest is the request body for adding a custom domain.
type CreateStatusPageCustomDomainRequest struct {
	Domain string `json:"domain"`
}

// GetStatusPage retrieves a status page by ID.
func (c *Client) GetStatusPage(id string) (*StatusPage, error) {
	var page StatusPage
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/status-pages/%s", id), nil, &page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// CreateStatusPage creates a new status page.
func (c *Client) CreateStatusPage(req CreateStatusPageRequest) (*StatusPage, error) {
	var page StatusPage
	err := c.doRequest(http.MethodPost, "/v1/status-pages", req, &page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// UpdateStatusPage updates an existing status page.
func (c *Client) UpdateStatusPage(id string, req UpdateStatusPageRequest) (*StatusPage, error) {
	var page StatusPage
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/status-pages/%s", id), req, &page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// ListStatusPages retrieves all status pages.
func (c *Client) ListStatusPages() ([]StatusPage, error) {
	var pages []StatusPage
	err := c.doRequest(http.MethodGet, "/v1/status-pages", nil, &pages)
	if err != nil {
		return nil, err
	}
	return pages, nil
}

// DeleteStatusPage deletes a status page by ID.
func (c *Client) DeleteStatusPage(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/status-pages/%s", id), nil, nil)
}

// GetStatusPageCustomDomain retrieves a custom domain by ID.
func (c *Client) GetStatusPageCustomDomain(pageID, domainID string) (*StatusPageCustomDomain, error) {
	var domain StatusPageCustomDomain
	path := fmt.Sprintf("/v1/status-pages/%s/domains/%s", pageID, domainID)
	err := c.doRequest(http.MethodGet, path, nil, &domain)
	if err != nil {
		return nil, err
	}
	return &domain, nil
}

// CreateStatusPageCustomDomain adds a custom domain to a status page.
func (c *Client) CreateStatusPageCustomDomain(pageID string, req CreateStatusPageCustomDomainRequest) (*StatusPageCustomDomain, error) {
	var domain StatusPageCustomDomain
	path := fmt.Sprintf("/v1/status-pages/%s/domains", pageID)
	err := c.doRequest(http.MethodPost, path, req, &domain)
	if err != nil {
		return nil, err
	}
	return &domain, nil
}

// DeleteStatusPageCustomDomain removes a custom domain from a status page.
func (c *Client) DeleteStatusPageCustomDomain(pageID, domainID string) error {
	path := fmt.Sprintf("/v1/status-pages/%s/domains/%s", pageID, domainID)
	return c.doRequest(http.MethodDelete, path, nil, nil)
}

// StatusPageMonitor represents a monitor linked to a status page.
type StatusPageMonitor struct {
	ID          string `json:"id"`
	PageID      string `json:"pageId"`
	MonitorID   string `json:"monitorId"`
	DisplayName string `json:"displayName"`
	SortOrder   int64  `json:"sortOrder"`
}

// CreateStatusPageMonitorRequest is the request body for adding a monitor to a status page.
type CreateStatusPageMonitorRequest struct {
	MonitorID   string `json:"monitorId"`
	DisplayName string `json:"displayName"`
	SortOrder   int64  `json:"sortOrder"`
}

// UpdateStatusPageMonitorRequest is the request body for updating a status page monitor.
type UpdateStatusPageMonitorRequest struct {
	DisplayName string `json:"displayName"`
	SortOrder   int64  `json:"sortOrder"`
}

// GetStatusPageMonitor retrieves a status page monitor by ID.
func (c *Client) GetStatusPageMonitor(pageID, monitorID string) (*StatusPageMonitor, error) {
	var monitor StatusPageMonitor
	path := fmt.Sprintf("/v1/status-pages/%s/monitors/%s", pageID, monitorID)
	err := c.doRequest(http.MethodGet, path, nil, &monitor)
	if err != nil {
		return nil, err
	}
	return &monitor, nil
}

// ListStatusPageMonitors retrieves all monitors for a status page.
func (c *Client) ListStatusPageMonitors(pageID string) ([]StatusPageMonitor, error) {
	var monitors []StatusPageMonitor
	path := fmt.Sprintf("/v1/status-pages/%s/monitors", pageID)
	err := c.doRequest(http.MethodGet, path, nil, &monitors)
	if err != nil {
		return nil, err
	}
	return monitors, nil
}

// CreateStatusPageMonitor adds a monitor to a status page.
func (c *Client) CreateStatusPageMonitor(pageID string, req CreateStatusPageMonitorRequest) (*StatusPageMonitor, error) {
	var monitor StatusPageMonitor
	path := fmt.Sprintf("/v1/status-pages/%s/monitors", pageID)
	err := c.doRequest(http.MethodPost, path, req, &monitor)
	if err != nil {
		return nil, err
	}
	return &monitor, nil
}

// UpdateStatusPageMonitor updates a status page monitor.
func (c *Client) UpdateStatusPageMonitor(pageID, monitorID string, req UpdateStatusPageMonitorRequest) (*StatusPageMonitor, error) {
	var monitor StatusPageMonitor
	path := fmt.Sprintf("/v1/status-pages/%s/monitors/%s", pageID, monitorID)
	err := c.doRequest(http.MethodPut, path, req, &monitor)
	if err != nil {
		return nil, err
	}
	return &monitor, nil
}

// DeleteStatusPageMonitor removes a monitor from a status page.
func (c *Client) DeleteStatusPageMonitor(pageID, monitorID string) error {
	path := fmt.Sprintf("/v1/status-pages/%s/monitors/%s", pageID, monitorID)
	return c.doRequest(http.MethodDelete, path, nil, nil)
}

// StatusPageIncident represents an incident on a status page.
type StatusPageIncident struct {
	ID          string `json:"id"`
	PageID      string `json:"pageId"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	ImpactLevel string `json:"impactLevel"`
	Message     string `json:"message,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
}

// CreateStatusPageIncidentRequest is the request body for creating a status page incident.
type CreateStatusPageIncidentRequest struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	ImpactLevel string `json:"impactLevel"`
	Message     string `json:"message,omitempty"`
}

// UpdateStatusPageIncidentRequest is the request body for updating a status page incident.
type UpdateStatusPageIncidentRequest struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	ImpactLevel string `json:"impactLevel"`
	Message     string `json:"message,omitempty"`
}

// GetStatusPageIncident retrieves a status page incident by ID.
func (c *Client) GetStatusPageIncident(pageID, incidentID string) (*StatusPageIncident, error) {
	var incident StatusPageIncident
	path := fmt.Sprintf("/v1/status-pages/%s/incidents/%s", pageID, incidentID)
	err := c.doRequest(http.MethodGet, path, nil, &incident)
	if err != nil {
		return nil, err
	}
	return &incident, nil
}

// ListStatusPageIncidents retrieves all incidents for a status page.
func (c *Client) ListStatusPageIncidents(pageID string) ([]StatusPageIncident, error) {
	var incidents []StatusPageIncident
	path := fmt.Sprintf("/v1/status-pages/%s/incidents", pageID)
	err := c.doRequest(http.MethodGet, path, nil, &incidents)
	if err != nil {
		return nil, err
	}
	return incidents, nil
}

// CreateStatusPageIncident creates a new status page incident.
func (c *Client) CreateStatusPageIncident(pageID string, req CreateStatusPageIncidentRequest) (*StatusPageIncident, error) {
	var incident StatusPageIncident
	path := fmt.Sprintf("/v1/status-pages/%s/incidents", pageID)
	err := c.doRequest(http.MethodPost, path, req, &incident)
	if err != nil {
		return nil, err
	}
	return &incident, nil
}

// UpdateStatusPageIncident updates an existing status page incident.
func (c *Client) UpdateStatusPageIncident(pageID, incidentID string, req UpdateStatusPageIncidentRequest) (*StatusPageIncident, error) {
	var incident StatusPageIncident
	path := fmt.Sprintf("/v1/status-pages/%s/incidents/%s", pageID, incidentID)
	err := c.doRequest(http.MethodPut, path, req, &incident)
	if err != nil {
		return nil, err
	}
	return &incident, nil
}

// DeleteStatusPageIncident deletes a status page incident by ID.
func (c *Client) DeleteStatusPageIncident(pageID, incidentID string) error {
	path := fmt.Sprintf("/v1/status-pages/%s/incidents/%s", pageID, incidentID)
	return c.doRequest(http.MethodDelete, path, nil, nil)
}
