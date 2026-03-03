package apiclient

import (
	"fmt"
	"net/http"
)

// IncidentProvider represents a Moneat incident provider.
type IncidentProvider struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Config    string `json:"config"`
	CreatedAt string `json:"createdAt,omitempty"`
}

// CreateIncidentProviderRequest is the request body for creating an incident provider.
type CreateIncidentProviderRequest struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Config string `json:"config"`
}

// UpdateIncidentProviderRequest is the request body for updating an incident provider.
type UpdateIncidentProviderRequest struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Config string `json:"config"`
}

// GetIncidentProvider retrieves an incident provider by ID.
func (c *Client) GetIncidentProvider(id string) (*IncidentProvider, error) {
	var provider IncidentProvider
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/api/incident-providers/%s", id), nil, &provider)
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

// ListIncidentProviders retrieves all incident providers.
func (c *Client) ListIncidentProviders() ([]IncidentProvider, error) {
	var providers []IncidentProvider
	err := c.doRequest(http.MethodGet, "/api/incident-providers", nil, &providers)
	if err != nil {
		return nil, err
	}
	return providers, nil
}

// CreateIncidentProvider creates a new incident provider.
func (c *Client) CreateIncidentProvider(req CreateIncidentProviderRequest) (*IncidentProvider, error) {
	var provider IncidentProvider
	err := c.doRequest(http.MethodPost, "/api/incident-providers", req, &provider)
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

// UpdateIncidentProvider updates an existing incident provider.
func (c *Client) UpdateIncidentProvider(id string, req UpdateIncidentProviderRequest) (*IncidentProvider, error) {
	var provider IncidentProvider
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/api/incident-providers/%s", id), req, &provider)
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

// DeleteIncidentProvider deletes an incident provider by ID.
func (c *Client) DeleteIncidentProvider(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/api/incident-providers/%s", id), nil, nil)
}

// IncidentRoutingRule represents an incident routing rule.
type IncidentRoutingRule struct {
	ID            string `json:"id"`
	ProviderID    string `json:"providerId"`
	Name          string `json:"name"`
	Condition     string `json:"condition"`
	TargetService string `json:"targetService"`
	CreatedAt     string `json:"createdAt,omitempty"`
}

// CreateIncidentRoutingRuleRequest is the request body for creating an incident routing rule.
type CreateIncidentRoutingRuleRequest struct {
	Name          string `json:"name"`
	Condition     string `json:"condition"`
	TargetService string `json:"targetService"`
}

// UpdateIncidentRoutingRuleRequest is the request body for updating an incident routing rule.
type UpdateIncidentRoutingRuleRequest struct {
	Name          string `json:"name"`
	Condition     string `json:"condition"`
	TargetService string `json:"targetService"`
}

// GetIncidentRoutingRule retrieves an incident routing rule by ID.
func (c *Client) GetIncidentRoutingRule(providerID, ruleID string) (*IncidentRoutingRule, error) {
	var rule IncidentRoutingRule
	path := fmt.Sprintf("/api/incident-providers/%s/rules/%s", providerID, ruleID)
	err := c.doRequest(http.MethodGet, path, nil, &rule)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// ListIncidentRoutingRules retrieves all routing rules for an incident provider.
func (c *Client) ListIncidentRoutingRules(providerID string) ([]IncidentRoutingRule, error) {
	var rules []IncidentRoutingRule
	path := fmt.Sprintf("/api/incident-providers/%s/rules", providerID)
	err := c.doRequest(http.MethodGet, path, nil, &rules)
	if err != nil {
		return nil, err
	}
	return rules, nil
}

// CreateIncidentRoutingRule creates a new incident routing rule.
func (c *Client) CreateIncidentRoutingRule(providerID string, req CreateIncidentRoutingRuleRequest) (*IncidentRoutingRule, error) {
	var rule IncidentRoutingRule
	path := fmt.Sprintf("/api/incident-providers/%s/rules", providerID)
	err := c.doRequest(http.MethodPost, path, req, &rule)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// UpdateIncidentRoutingRule updates an existing incident routing rule.
func (c *Client) UpdateIncidentRoutingRule(providerID, ruleID string, req UpdateIncidentRoutingRuleRequest) (*IncidentRoutingRule, error) {
	var rule IncidentRoutingRule
	path := fmt.Sprintf("/api/incident-providers/%s/rules/%s", providerID, ruleID)
	err := c.doRequest(http.MethodPut, path, req, &rule)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// DeleteIncidentRoutingRule deletes an incident routing rule by ID.
func (c *Client) DeleteIncidentRoutingRule(providerID, ruleID string) error {
	path := fmt.Sprintf("/api/incident-providers/%s/rules/%s", providerID, ruleID)
	return c.doRequest(http.MethodDelete, path, nil, nil)
}
