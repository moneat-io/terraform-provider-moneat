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
