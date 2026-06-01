package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// McpAPIKey represents a Moneat MCP API key.
type McpAPIKey struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	KeyPrefix        string   `json:"keyPrefix"`
	Key              string   `json:"key,omitempty"`
	EnabledTools     []string `json:"enabledTools"`
	EnabledResources []string `json:"enabledResources"`
	LastUsedAt       string   `json:"lastUsedAt,omitempty"`
	ExpiresAt        string   `json:"expiresAt,omitempty"`
	CreatedAt        string   `json:"createdAt,omitempty"`
}

// CreateMcpAPIKeyRequest is the request body for creating an MCP API key.
type CreateMcpAPIKeyRequest struct {
	Name             string   `json:"name"`
	EnabledTools     []string `json:"enabledTools"`
	EnabledResources []string `json:"enabledResources"`
	ExpiresInDays    *int     `json:"expiresInDays,omitempty"`
}

// UpdateMcpAPIKeyRequest is the request body for updating an MCP API key.
type UpdateMcpAPIKeyRequest struct {
	Name             string   `json:"name,omitempty"`
	EnabledTools     []string `json:"enabledTools"`
	EnabledResources []string `json:"enabledResources"`
	ExpiresInDays    *int     `json:"expiresInDays,omitempty"`
}

type mcpAPIKeysResponse struct {
	Keys []McpAPIKey `json:"keys"`
}

// ListMcpAPIKeys retrieves all MCP API keys.
func (c *Client) ListMcpAPIKeys() ([]McpAPIKey, error) {
	var response mcpAPIKeysResponse
	err := c.doRequest(http.MethodGet, "/v1/mcp/api-keys", nil, &response)
	if err != nil {
		return nil, err
	}
	return response.Keys, nil
}

// GetMcpAPIKey retrieves an MCP API key by ID.
func (c *Client) GetMcpAPIKey(id int) (*McpAPIKey, error) {
	keys, err := c.ListMcpAPIKeys()
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		if key.ID == id {
			return &key, nil
		}
	}
	return nil, notFoundError("MCP API key not found")
}

// CreateMcpAPIKey creates an MCP API key.
func (c *Client) CreateMcpAPIKey(req CreateMcpAPIKeyRequest) (*McpAPIKey, error) {
	var key McpAPIKey
	err := c.doRequest(http.MethodPost, "/v1/mcp/api-keys", req, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// UpdateMcpAPIKey updates an MCP API key.
func (c *Client) UpdateMcpAPIKey(id int, req UpdateMcpAPIKeyRequest) (*McpAPIKey, error) {
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/mcp/api-keys/%d", id), req, nil)
	if err != nil {
		return nil, err
	}
	return c.GetMcpAPIKey(id)
}

// DeleteMcpAPIKey revokes an MCP API key.
func (c *Client) DeleteMcpAPIKey(id int) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/mcp/api-keys/%d", id), nil, nil)
}

// GetMcpToolCatalog retrieves the MCP tool and resource catalog as JSON.
func (c *Client) GetMcpToolCatalog() (json.RawMessage, error) {
	var catalog json.RawMessage
	err := c.doRequest(http.MethodGet, "/v1/mcp/tool-catalog", nil, &catalog)
	if err != nil {
		return nil, err
	}
	return catalog, nil
}
