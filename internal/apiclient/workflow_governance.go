package apiclient

import "net/http"

type WorkflowGrant struct {
	ID             string   `json:"id"`
	WorkflowID     string   `json:"workflow_id"`
	UserID         string   `json:"user_id"`
	Role           string   `json:"role"`
	AllowedActions []string `json:"allowed_actions"`
	GrantedBy      string   `json:"granted_by,omitempty"`
	CreatedAt      string   `json:"created_at,omitempty"`
	UpdatedAt      string   `json:"updated_at,omitempty"`
}

type WorkflowGrantRequest struct {
	UserID         string   `json:"user_id"`
	Role           string   `json:"role"`
	AllowedActions []string `json:"allowed_actions"`
}

func workflowGrantPath(workflowID string) string {
	return workflowPath(workflowID) + "/grants"
}

func (c *Client) ListWorkflowGrants(workflowID string) ([]WorkflowGrant, error) {
	var grants []WorkflowGrant
	if err := c.doRequest(http.MethodGet, workflowGrantPath(workflowID), nil, &grants); err != nil {
		return nil, err
	}
	return grants, nil
}

func (c *Client) UpsertWorkflowGrant(workflowID string, req WorkflowGrantRequest) (*WorkflowGrant, error) {
	var grant WorkflowGrant
	if err := c.doRequest(http.MethodPost, workflowGrantPath(workflowID), req, &grant); err != nil {
		return nil, err
	}
	return &grant, nil
}

func (c *Client) UpdateWorkflowGrant(workflowID, grantID string, req WorkflowGrantRequest) (*WorkflowGrant, error) {
	var grant WorkflowGrant
	path := workflowGrantPath(workflowID) + "/" + grantID
	if err := c.doRequest(http.MethodPut, path, req, &grant); err != nil {
		return nil, err
	}
	return &grant, nil
}

func (c *Client) DeleteWorkflowGrant(workflowID, grantID string) error {
	return c.doRequest(http.MethodDelete, workflowGrantPath(workflowID)+"/"+grantID, nil, nil)
}

type WorkflowServicePrincipal struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	TokenPrefix string `json:"token_prefix"`
	CreatedBy   string `json:"created_by,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	RotatedAt   string `json:"rotated_at,omitempty"`
}

type CreateWorkflowServicePrincipalRequest struct {
	Name string `json:"name"`
}

type WorkflowServicePrincipalSecretResponse struct {
	Principal WorkflowServicePrincipal `json:"principal"`
	Token     string                   `json:"token"`
}

func (c *Client) ListWorkflowServicePrincipals() ([]WorkflowServicePrincipal, error) {
	var principals []WorkflowServicePrincipal
	if err := c.doRequest(http.MethodGet, "/v1/workflows/service-principals", nil, &principals); err != nil {
		return nil, err
	}
	return principals, nil
}

func (c *Client) CreateWorkflowServicePrincipal(req CreateWorkflowServicePrincipalRequest) (*WorkflowServicePrincipalSecretResponse, error) {
	var response WorkflowServicePrincipalSecretResponse
	if err := c.doRequest(http.MethodPost, "/v1/workflows/service-principals", req, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) RotateWorkflowServicePrincipal(id string) (*WorkflowServicePrincipalSecretResponse, error) {
	var response WorkflowServicePrincipalSecretResponse
	path := "/v1/workflows/service-principals/" + id + "/rotate"
	if err := c.doRequest(http.MethodPost, path, nil, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *Client) DeleteWorkflowServicePrincipal(id string) error {
	return c.doRequest(http.MethodDelete, "/v1/workflows/service-principals/"+id, nil, nil)
}

type WorkflowOrganizationPolicy struct {
	ID                     string   `json:"id"`
	OrganizationID         string   `json:"organization_id"`
	DefaultRole            string   `json:"default_role"`
	RequireApproval        bool     `json:"require_approval"`
	AllowServicePrincipals bool     `json:"allow_service_principals"`
	MaxConcurrentRuns      *int     `json:"max_concurrent_runs,omitempty"`
	AllowedTriggers        []string `json:"allowed_triggers"`
	CreatedAt              string   `json:"created_at,omitempty"`
	UpdatedAt              string   `json:"updated_at,omitempty"`
}

type WorkflowOrganizationPolicyRequest struct {
	DefaultRole            string   `json:"default_role"`
	RequireApproval        bool     `json:"require_approval"`
	AllowServicePrincipals bool     `json:"allow_service_principals"`
	MaxConcurrentRuns      *int     `json:"max_concurrent_runs,omitempty"`
	AllowedTriggers        []string `json:"allowed_triggers"`
}

func (c *Client) GetWorkflowOrganizationPolicy() (*WorkflowOrganizationPolicy, error) {
	var policy WorkflowOrganizationPolicy
	if err := c.doRequest(http.MethodGet, "/v1/workflows/organization-policy", nil, &policy); err != nil {
		return nil, err
	}
	return &policy, nil
}

func (c *Client) SetWorkflowOrganizationPolicy(req WorkflowOrganizationPolicyRequest) (*WorkflowOrganizationPolicy, error) {
	var policy WorkflowOrganizationPolicy
	if err := c.doRequest(http.MethodPut, "/v1/workflows/organization-policy", req, &policy); err != nil {
		return nil, err
	}
	return &policy, nil
}

type WorkflowBlueprint struct {
	Key string `json:"key"`
}
