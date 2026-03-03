package apiclient

import (
	"fmt"
	"net/http"
)

// AuthToken represents a Moneat auth token.
type AuthToken struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Token     string   `json:"token"`
	Scopes    []string `json:"scopes"`
	CreatedAt string   `json:"createdAt,omitempty"`
}

// CreateAuthTokenRequest is the request body for creating an auth token.
type CreateAuthTokenRequest struct {
	Name   string   `json:"name"`
	Scopes []string `json:"scopes"`
}

// UpdateAuthTokenRequest is the request body for updating an auth token.
type UpdateAuthTokenRequest struct {
	Name   string   `json:"name"`
	Scopes []string `json:"scopes"`
}

// GetAuthToken retrieves an auth token by ID.
func (c *Client) GetAuthToken(id string) (*AuthToken, error) {
	var token AuthToken
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/auth-tokens/%s", id), nil, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// ListAuthTokens retrieves all auth tokens.
func (c *Client) ListAuthTokens() ([]AuthToken, error) {
	var tokens []AuthToken
	err := c.doRequest(http.MethodGet, "/v1/auth-tokens", nil, &tokens)
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

// CreateAuthToken creates a new auth token.
func (c *Client) CreateAuthToken(req CreateAuthTokenRequest) (*AuthToken, error) {
	var token AuthToken
	err := c.doRequest(http.MethodPost, "/v1/auth-tokens", req, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// UpdateAuthToken updates an existing auth token.
func (c *Client) UpdateAuthToken(id string, req UpdateAuthTokenRequest) (*AuthToken, error) {
	var token AuthToken
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/auth-tokens/%s", id), req, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// DeleteAuthToken deletes an auth token by ID.
func (c *Client) DeleteAuthToken(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/auth-tokens/%s", id), nil, nil)
}

// AgentAPIKey represents a Moneat agent API key.
type AgentAPIKey struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Key       string `json:"key"`
	CreatedAt string `json:"createdAt,omitempty"`
}

// CreateAgentAPIKeyRequest is the request body for creating an agent API key.
type CreateAgentAPIKeyRequest struct {
	Name string `json:"name"`
}

// UpdateAgentAPIKeyRequest is the request body for updating an agent API key.
type UpdateAgentAPIKeyRequest struct {
	Name string `json:"name"`
}

// GetAgentAPIKey retrieves an agent API key by ID.
func (c *Client) GetAgentAPIKey(id string) (*AgentAPIKey, error) {
	var key AgentAPIKey
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/agent-api-keys/%s", id), nil, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// ListAgentAPIKeys retrieves all agent API keys.
func (c *Client) ListAgentAPIKeys() ([]AgentAPIKey, error) {
	var keys []AgentAPIKey
	err := c.doRequest(http.MethodGet, "/v1/agent-api-keys", nil, &keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// CreateAgentAPIKey creates a new agent API key.
func (c *Client) CreateAgentAPIKey(req CreateAgentAPIKeyRequest) (*AgentAPIKey, error) {
	var key AgentAPIKey
	err := c.doRequest(http.MethodPost, "/v1/agent-api-keys", req, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// UpdateAgentAPIKey updates an existing agent API key.
func (c *Client) UpdateAgentAPIKey(id string, req UpdateAgentAPIKeyRequest) (*AgentAPIKey, error) {
	var key AgentAPIKey
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/agent-api-keys/%s", id), req, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// DeleteAgentAPIKey deletes an agent API key by ID.
func (c *Client) DeleteAgentAPIKey(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/agent-api-keys/%s", id), nil, nil)
}

// SSOConfig represents a Moneat SSO configuration.
type SSOConfig struct {
	ID          string `json:"id"`
	Provider    string `json:"provider"`
	EntityID    string `json:"entityId"`
	SSOURL      string `json:"ssoUrl"`
	Certificate string `json:"certificate"`
	Enabled     bool   `json:"enabled"`
}

// UpdateSSOConfigRequest is the request body for updating SSO configuration.
type UpdateSSOConfigRequest struct {
	Provider    string `json:"provider"`
	EntityID    string `json:"entityId"`
	SSOURL      string `json:"ssoUrl"`
	Certificate string `json:"certificate"`
	Enabled     bool   `json:"enabled"`
}

// GetSSOConfig retrieves the SSO configuration.
func (c *Client) GetSSOConfig() (*SSOConfig, error) {
	var config SSOConfig
	err := c.doRequest(http.MethodGet, "/v1/sso/config", nil, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// UpdateSSOConfig updates the SSO configuration.
func (c *Client) UpdateSSOConfig(req UpdateSSOConfigRequest) (*SSOConfig, error) {
	var config SSOConfig
	err := c.doRequest(http.MethodPut, "/v1/sso/config", req, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
