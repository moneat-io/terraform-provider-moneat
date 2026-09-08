package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ResponseConfigurationState is the durable, organization-scoped configuration snapshot
// used by Terraform and other configuration-as-code clients.
type ResponseConfigurationState struct {
	ID            string          `json:"id,omitempty"`
	Revision      int             `json:"revision"`
	Checksum      string          `json:"checksum"`
	Source        string          `json:"source"`
	Configuration json.RawMessage `json:"configuration"`
	CreatedAt     string          `json:"createdAt,omitempty"`
}

type ResponseConfigurationChange struct {
	Action      string `json:"action"`
	Type        string `json:"type"`
	ID          string `json:"id"`
	Destructive bool   `json:"destructive"`
}

type ResponseConfigurationIssue struct {
	Code    string `json:"code"`
	Path    string `json:"path"`
	Message string `json:"message"`
}

type ResponseConfigurationPlan struct {
	BaseRevision     int                           `json:"baseRevision"`
	ProposedRevision int                           `json:"proposedRevision"`
	Checksum         string                        `json:"checksum"`
	Changes          []ResponseConfigurationChange `json:"changes"`
	Issues           []ResponseConfigurationIssue  `json:"issues"`
	CanApply         bool                          `json:"canApply"`
	Configuration    json.RawMessage               `json:"configuration"`
}

type responseConfigurationPlanRequest struct {
	Configuration    json.RawMessage `json:"configuration"`
	ExpectedRevision *int            `json:"expectedRevision,omitempty"`
	AllowDestructive bool            `json:"allowDestructive"`
}

type responseConfigurationApplyRequest struct {
	Plan             ResponseConfigurationPlan `json:"plan"`
	IdempotencyKey   string                    `json:"idempotencyKey"`
	AllowDestructive bool                      `json:"allowDestructive"`
}

func (c *Client) GetResponseConfigurationState() (*ResponseConfigurationState, error) {
	var state ResponseConfigurationState
	err := c.doRequest(http.MethodGet, "/api/v1/response/configuration/state", nil, &state)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (c *Client) PlanResponseConfiguration(
	configuration json.RawMessage,
	expectedRevision *int,
	allowDestructive bool,
) (*ResponseConfigurationPlan, error) {
	var plan ResponseConfigurationPlan
	err := c.doRequest(
		http.MethodPost,
		"/api/v1/response/configuration/plan",
		responseConfigurationPlanRequest{
			Configuration:    configuration,
			ExpectedRevision: expectedRevision,
			AllowDestructive: allowDestructive,
		},
		&plan,
	)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (c *Client) ApplyResponseConfiguration(
	plan *ResponseConfigurationPlan,
	idempotencyKey string,
	allowDestructive bool,
) (*ResponseConfigurationState, error) {
	if plan == nil {
		return nil, fmt.Errorf("response configuration plan cannot be nil")
	}
	var state ResponseConfigurationState
	err := c.doRequestWithHeaders(
		http.MethodPost,
		"/api/v1/response/configuration/apply",
		responseConfigurationApplyRequest{
			Plan:             *plan,
			IdempotencyKey:   idempotencyKey,
			AllowDestructive: allowDestructive,
		},
		&state,
		map[string]string{"Idempotency-Key": idempotencyKey},
	)
	if err != nil {
		return nil, err
	}
	return &state, nil
}
