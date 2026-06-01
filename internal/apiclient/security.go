package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// DetectionRule represents a Moneat security detection rule.
type DetectionRule struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Source         string   `json:"source"`
	Filter         string   `json:"filter"`
	GroupBy        []string `json:"group_by"`
	WindowSeconds  int      `json:"window_seconds"`
	Type           string   `json:"type"`
	ThresholdCount *int     `json:"threshold_count,omitempty"`
	Severity       string   `json:"severity"`
	SignalTitle    string   `json:"signal_title"`
	SignalMessage  string   `json:"signal_message"`
	Suppressions   []string `json:"suppressions"`
	Enabled        bool     `json:"enabled"`
	Tags           []string `json:"tags"`
	CreatedAt      string   `json:"created_at,omitempty"`
	UpdatedAt      string   `json:"updated_at,omitempty"`
}

// DetectionRuleRequest is the request body for creating or updating a detection rule.
type DetectionRuleRequest struct {
	Name           string   `json:"name,omitempty"`
	Description    string   `json:"description,omitempty"`
	Source         string   `json:"source,omitempty"`
	Filter         string   `json:"filter,omitempty"`
	GroupBy        []string `json:"group_by,omitempty"`
	WindowSeconds  int      `json:"window_seconds,omitempty"`
	Type           string   `json:"type,omitempty"`
	ThresholdCount *int     `json:"threshold_count,omitempty"`
	Severity       string   `json:"severity,omitempty"`
	SignalTitle    string   `json:"signal_title,omitempty"`
	SignalMessage  string   `json:"signal_message,omitempty"`
	Suppressions   []string `json:"suppressions,omitempty"`
	Enabled        *bool    `json:"enabled,omitempty"`
	Tags           []string `json:"tags,omitempty"`
}

type detectionRuleListResponse struct {
	Rules []DetectionRule `json:"rules"`
}

// GetDetectionRule retrieves a detection rule by ID.
func (c *Client) GetDetectionRule(id int) (*DetectionRule, error) {
	var rule DetectionRule
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/security/detection/rules/%d", id), nil, &rule)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// ListDetectionRules retrieves detection rules.
func (c *Client) ListDetectionRules() ([]DetectionRule, error) {
	var response detectionRuleListResponse
	err := c.doRequest(http.MethodGet, "/v1/security/detection/rules", nil, &response)
	if err != nil {
		return nil, err
	}
	return response.Rules, nil
}

// CreateDetectionRule creates a detection rule.
func (c *Client) CreateDetectionRule(req DetectionRuleRequest) (*DetectionRule, error) {
	var rule DetectionRule
	err := c.doRequest(http.MethodPost, "/v1/security/detection/rules", req, &rule)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// UpdateDetectionRule updates a detection rule.
func (c *Client) UpdateDetectionRule(id int, req DetectionRuleRequest) (*DetectionRule, error) {
	var rule DetectionRule
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/security/detection/rules/%d", id), req, &rule)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

// DeleteDetectionRule deletes a detection rule.
func (c *Client) DeleteDetectionRule(id int) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/security/detection/rules/%d", id), nil, nil)
}

// GetSecurityDetectionCoverage retrieves MITRE detection coverage as JSON.
func (c *Client) GetSecurityDetectionCoverage() (json.RawMessage, error) {
	var response json.RawMessage
	err := c.doRequest(http.MethodGet, "/v1/security/detection/coverage", nil, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// ListSecuritySignals retrieves security signals as JSON.
func (c *Client) ListSecuritySignals(params map[string]string) (json.RawMessage, error) {
	var response json.RawMessage
	err := c.doRequest(http.MethodGet, "/v1/security/signals"+queryString(params), nil, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// GetVulnerabilitySummary retrieves vulnerability summary as JSON.
func (c *Client) GetVulnerabilitySummary() (json.RawMessage, error) {
	var response json.RawMessage
	err := c.doRequest(http.MethodGet, "/v1/security/vulnerabilities/summary", nil, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// ListVulnerabilityFindings retrieves vulnerability findings as JSON.
func (c *Client) ListVulnerabilityFindings(params map[string]string) (json.RawMessage, error) {
	var response json.RawMessage
	path := "/v1/security/vulnerabilities/findings" + queryString(params)
	err := c.doRequest(http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

// ListVulnerabilityInventory retrieves vulnerability package inventory as JSON.
func (c *Client) ListVulnerabilityInventory(params map[string]string) (json.RawMessage, error) {
	var response json.RawMessage
	path := "/v1/security/vulnerabilities/inventory" + queryString(params)
	err := c.doRequest(http.MethodGet, path, nil, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func queryString(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	values := url.Values{}
	for key, value := range params {
		if value != "" {
			values.Set(key, value)
		}
	}
	encoded := values.Encode()
	if encoded == "" {
		return ""
	}
	return "?" + encoded
}
