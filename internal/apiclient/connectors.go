package apiclient

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// ConnectorExternalAccount identifies the provider-side account to import.
type ConnectorExternalAccount struct {
	ProjectID         string `json:"projectId,omitempty"`
	CustomerID        string `json:"customerId,omitempty"`
	ManagerCustomerID string `json:"managerCustomerId,omitempty"`
}

// ConnectorInstallation is the redacted connector installation contract. The
// API never returns the plaintext secret; Terraform retains it as sensitive
// state and sends it only on create/rotation.
type ConnectorInstallation struct {
	ID                           string            `json:"id"`
	ProviderID                   string            `json:"providerId"`
	UseID                        string            `json:"useId"`
	Name                         string            `json:"name"`
	CredentialType               string            `json:"credentialType"`
	AuthProfileID                string            `json:"authProfileId"`
	IdentifierTags               map[string]string `json:"identifierTags"`
	ExternalProjectID            string            `json:"externalProjectId,omitempty"`
	ExternalProjectName          string            `json:"externalProjectName,omitempty"`
	Status                       string            `json:"status"`
	StatusReason                 string            `json:"statusReason,omitempty"`
	Enabled                      bool              `json:"enabled"`
	APISecretLastFour            string            `json:"apiSecretLastFour,omitempty"`
	SecretVersion                int               `json:"secretVersion,omitempty"`
	KeyID                        string            `json:"keyId,omitempty"`
	SecretCreatedAt              string            `json:"secretCreatedAt,omitempty"`
	WorkflowConfig               json.RawMessage   `json:"workflowConfig,omitempty"`
	WebhookTokenPrefix           string            `json:"webhookTokenPrefix,omitempty"`
	LastTestedAt                 string            `json:"lastTestedAt,omitempty"`
	LastTestResult               string            `json:"lastTestResult,omitempty"`
	LastSuccessfulProviderCallAt string            `json:"lastSuccessfulProviderCallAt,omitempty"`
	LastError                    string            `json:"lastError,omitempty"`
	CreatedAt                    string            `json:"createdAt,omitempty"`
	UpdatedAt                    string            `json:"updatedAt,omitempty"`
}

type ConnectorInstallationsResponse struct {
	Installations []ConnectorInstallation `json:"installations"`
}

type CreateConnectorInstallationRequest struct {
	ProviderID      string                   `json:"providerId"`
	AuthProfileID   string                   `json:"authProfileId"`
	Name            string                   `json:"name"`
	ExternalAccount ConnectorExternalAccount `json:"externalAccount"`
	Secret          string                   `json:"secret"`
	UseID           string                   `json:"useId,omitempty"`
	IdentifierTags  map[string]string        `json:"identifierTags,omitempty"`
}

type RotateConnectorCredentialRequest struct {
	Secret string `json:"secret"`
}

type ConnectorGroup struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name"`
	ProviderID            string   `json:"providerId"`
	UseID                 string   `json:"useId"`
	MemberInstallationIDs []string `json:"memberInstallationIds"`
	SelectionStrategy     string   `json:"selectionStrategy"`
	CreatedAt             string   `json:"createdAt,omitempty"`
	UpdatedAt             string   `json:"updatedAt,omitempty"`
}

type ConnectorGroupsResponse struct {
	Groups []ConnectorGroup `json:"groups"`
}

type CreateConnectorGroupRequest struct {
	Name                  string   `json:"name"`
	ProviderID            string   `json:"providerId"`
	UseID                 string   `json:"useId"`
	MemberInstallationIDs []string `json:"memberInstallationIds"`
	SelectionStrategy     string   `json:"selectionStrategy"`
}

type ConnectorExternalResource struct {
	ID                   string            `json:"id"`
	InstallationID       string            `json:"installationId"`
	ExternalProjectID    string            `json:"externalProjectId,omitempty"`
	ExternalResourceType string            `json:"externalResourceType"`
	ExternalResourceID   string            `json:"externalResourceId"`
	DisplayName          string            `json:"displayName,omitempty"`
	ProviderMetadata     map[string]string `json:"providerMetadata,omitempty"`
	LastSeenAt           string            `json:"lastSeenAt,omitempty"`
}

type ConnectorResourcesResponse struct {
	Resources []ConnectorExternalResource `json:"resources"`
}

func connectorPath(parts ...string) string {
	path := "/v1/connectors"
	for _, part := range parts {
		path += "/" + url.PathEscape(part)
	}
	return path
}

func connectorQuery(path string, providerID, useID string) string {
	query := url.Values{}
	if providerID != "" {
		query.Set("providerId", providerID)
	}
	if useID != "" {
		query.Set("useId", useID)
	}
	if encoded := query.Encode(); encoded != "" {
		return path + "?" + encoded
	}
	return path
}

func (c *Client) ListConnectorProviders() (json.RawMessage, error) {
	var value json.RawMessage
	if err := c.doRequest(http.MethodGet, connectorPath("providers"), nil, &value); err != nil {
		return nil, err
	}
	return value, nil
}

func (c *Client) ListConnectorInstallations(providerID, useID string) ([]ConnectorInstallation, error) {
	var response ConnectorInstallationsResponse
	if err := c.doRequest(http.MethodGet, connectorQuery(connectorPath("installations"), providerID, useID), nil, &response); err != nil {
		return nil, err
	}
	return response.Installations, nil
}

func (c *Client) GetConnectorInstallation(id string) (*ConnectorInstallation, error) {
	var installation ConnectorInstallation
	if err := c.doRequest(http.MethodGet, connectorPath("installations", id), nil, &installation); err != nil {
		return nil, err
	}
	return &installation, nil
}

func (c *Client) CreateConnectorInstallation(req CreateConnectorInstallationRequest) (*ConnectorInstallation, error) {
	var installation ConnectorInstallation
	if err := c.doRequest(http.MethodPost, connectorPath("installations"), req, &installation); err != nil {
		return nil, err
	}
	return &installation, nil
}

func (c *Client) RotateConnectorCredential(id string, req RotateConnectorCredentialRequest) (*ConnectorInstallation, error) {
	var installation ConnectorInstallation
	path := connectorPath("installations", id, "credentials", "rotate")
	if err := c.doRequest(http.MethodPost, path, req, &installation); err != nil {
		return nil, err
	}
	return &installation, nil
}

func (c *Client) DeleteConnectorInstallation(id string) error {
	return c.doRequest(http.MethodDelete, connectorPath("installations", id), nil, nil)
}

func (c *Client) ListConnectorGroups(providerID, useID string) ([]ConnectorGroup, error) {
	var response ConnectorGroupsResponse
	if err := c.doRequest(http.MethodGet, connectorQuery(connectorPath("connection-groups"), providerID, useID), nil, &response); err != nil {
		return nil, err
	}
	return response.Groups, nil
}

func (c *Client) GetConnectorGroup(id string) (*ConnectorGroup, error) {
	groups, err := c.ListConnectorGroups("", "")
	if err != nil {
		return nil, err
	}
	for index := range groups {
		if groups[index].ID == id {
			return &groups[index], nil
		}
	}
	return nil, notFoundError("connector connection group not found")
}

func (c *Client) CreateConnectorGroup(req CreateConnectorGroupRequest) (*ConnectorGroup, error) {
	var group ConnectorGroup
	if err := c.doRequest(http.MethodPost, connectorPath("connection-groups"), req, &group); err != nil {
		return nil, err
	}
	return &group, nil
}

func (c *Client) DeleteConnectorGroup(id string) error {
	return c.doRequest(http.MethodDelete, connectorPath("connection-groups", id), nil, nil)
}

func (c *Client) ListConnectorResources(installationID string) ([]ConnectorExternalResource, error) {
	var response ConnectorResourcesResponse
	if err := c.doRequest(http.MethodGet, connectorPath("installations", installationID, "resources"), nil, &response); err != nil {
		return nil, err
	}
	return response.Resources, nil
}
