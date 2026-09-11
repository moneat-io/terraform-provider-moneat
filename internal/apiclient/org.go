package apiclient

import (
	"fmt"
	"net/http"
	"net/url"
)

// OrgMember represents a Moneat organization member.
type OrgMember struct {
	// The current API exposes the user's opaque resource identifier as userId.
	// Keep the provider model's field name stable while making the UUID contract
	// explicit at the wire boundary.
	ID       string `json:"userId"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	JoinedAt string `json:"joinedAt,omitempty"`
}

// OrgMembersResponse is the current organization-members response envelope.
// Pending invitations are intentionally not modelled by the provider member
// data source, but the envelope must be consumed before returning members.
type OrgMembersResponse struct {
	Members []OrgMember `json:"members"`
}

// UpdateOrgMemberRequest is the request body for updating an organization member's role.
type UpdateOrgMemberRequest struct {
	Role string `json:"role"`
}

// ListOrgMembers retrieves all organization members.
func (c *Client) ListOrgMembers() ([]OrgMember, error) {
	var response OrgMembersResponse
	err := c.doRequest(http.MethodGet, "/v1/org/members", nil, &response)
	if err != nil {
		return nil, err
	}
	return response.Members, nil
}

// UpdateOrgMember updates an organization member's role.
//
// The current API returns an acknowledgement rather than a member resource;
// callers should reload the organization member envelope when they need the
// canonical email and role values.
func (c *Client) UpdateOrgMember(id string, req UpdateOrgMemberRequest) error {
	return c.doRequest(http.MethodPut, fmt.Sprintf("/v1/org/members/%s/role", url.PathEscape(id)), req, nil)
}

// DeleteOrgMember removes a member from the organization.
func (c *Client) DeleteOrgMember(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/org/members/%s", url.PathEscape(id)), nil, nil)
}

// OrgInvitation represents a Moneat organization invitation.
type OrgInvitation struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt,omitempty"`
}

// CreateOrgInvitationRequest is the request body for creating an organization invitation.
type CreateOrgInvitationRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// UpdateOrgInvitationRequest is the request body for updating an organization invitation.
type UpdateOrgInvitationRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

// GetOrgInvitation retrieves an organization invitation by ID.
func (c *Client) GetOrgInvitation(id string) (*OrgInvitation, error) {
	var invitation OrgInvitation
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/org/invitations/%s", id), nil, &invitation)
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

// ListOrgInvitations retrieves all organization invitations.
func (c *Client) ListOrgInvitations() ([]OrgInvitation, error) {
	var invitations []OrgInvitation
	err := c.doRequest(http.MethodGet, "/v1/org/invitations", nil, &invitations)
	if err != nil {
		return nil, err
	}
	return invitations, nil
}

// CreateOrgInvitation creates a new organization invitation.
func (c *Client) CreateOrgInvitation(req CreateOrgInvitationRequest) (*OrgInvitation, error) {
	var invitation OrgInvitation
	err := c.doRequest(http.MethodPost, "/v1/org/invitations", req, &invitation)
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

// UpdateOrgInvitation updates an existing organization invitation.
func (c *Client) UpdateOrgInvitation(id string, req UpdateOrgInvitationRequest) (*OrgInvitation, error) {
	var invitation OrgInvitation
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/org/invitations/%s", id), req, &invitation)
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

// DeleteOrgInvitation deletes an organization invitation by ID.
func (c *Client) DeleteOrgInvitation(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/org/invitations/%s", id), nil, nil)
}
