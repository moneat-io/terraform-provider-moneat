package apiclient

import (
	"fmt"
	"net/http"
)

// OnCallSchedule represents a Moneat on-call schedule.
type OnCallSchedule struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Timezone     string `json:"timezone"`
	RotationType string `json:"rotationType"`
	CreatedAt    string `json:"createdAt,omitempty"`
}

// CreateOnCallScheduleRequest is the request body for creating an on-call schedule.
type CreateOnCallScheduleRequest struct {
	Name         string `json:"name"`
	Timezone     string `json:"timezone"`
	RotationType string `json:"rotationType"`
}

// UpdateOnCallScheduleRequest is the request body for updating an on-call schedule.
type UpdateOnCallScheduleRequest struct {
	Name         string `json:"name"`
	Timezone     string `json:"timezone"`
	RotationType string `json:"rotationType"`
}

// GetOnCallSchedule retrieves an on-call schedule by ID.
func (c *Client) GetOnCallSchedule(id string) (*OnCallSchedule, error) {
	var schedule OnCallSchedule
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/on-call/schedules/%s", id), nil, &schedule)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

// ListOnCallSchedules retrieves all on-call schedules.
func (c *Client) ListOnCallSchedules() ([]OnCallSchedule, error) {
	var schedules []OnCallSchedule
	err := c.doRequest(http.MethodGet, "/v1/on-call/schedules", nil, &schedules)
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

// CreateOnCallSchedule creates a new on-call schedule.
func (c *Client) CreateOnCallSchedule(req CreateOnCallScheduleRequest) (*OnCallSchedule, error) {
	var schedule OnCallSchedule
	err := c.doRequest(http.MethodPost, "/v1/on-call/schedules", req, &schedule)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

// UpdateOnCallSchedule updates an existing on-call schedule.
func (c *Client) UpdateOnCallSchedule(id string, req UpdateOnCallScheduleRequest) (*OnCallSchedule, error) {
	var schedule OnCallSchedule
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/on-call/schedules/%s", id), req, &schedule)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

// DeleteOnCallSchedule deletes an on-call schedule by ID.
func (c *Client) DeleteOnCallSchedule(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/on-call/schedules/%s", id), nil, nil)
}

// OnCallOverride represents an on-call schedule override.
type OnCallOverride struct {
	ID         string `json:"id"`
	ScheduleID string `json:"scheduleId"`
	UserID     string `json:"userId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	Reason     string `json:"reason,omitempty"`
}

// CreateOnCallOverrideRequest is the request body for creating an on-call override.
type CreateOnCallOverrideRequest struct {
	UserID    string `json:"userId"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Reason    string `json:"reason,omitempty"`
}

// UpdateOnCallOverrideRequest is the request body for updating an on-call override.
type UpdateOnCallOverrideRequest struct {
	UserID    string `json:"userId"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Reason    string `json:"reason,omitempty"`
}

// GetOnCallOverride retrieves an on-call override by ID.
func (c *Client) GetOnCallOverride(scheduleID, overrideID string) (*OnCallOverride, error) {
	var override OnCallOverride
	path := fmt.Sprintf("/v1/on-call/schedules/%s/overrides/%s", scheduleID, overrideID)
	err := c.doRequest(http.MethodGet, path, nil, &override)
	if err != nil {
		return nil, err
	}
	return &override, nil
}

// ListOnCallOverrides retrieves all overrides for an on-call schedule.
func (c *Client) ListOnCallOverrides(scheduleID string) ([]OnCallOverride, error) {
	var overrides []OnCallOverride
	path := fmt.Sprintf("/v1/on-call/schedules/%s/overrides", scheduleID)
	err := c.doRequest(http.MethodGet, path, nil, &overrides)
	if err != nil {
		return nil, err
	}
	return overrides, nil
}

// CreateOnCallOverride creates a new on-call override.
func (c *Client) CreateOnCallOverride(scheduleID string, req CreateOnCallOverrideRequest) (*OnCallOverride, error) {
	var override OnCallOverride
	path := fmt.Sprintf("/v1/on-call/schedules/%s/overrides", scheduleID)
	err := c.doRequest(http.MethodPost, path, req, &override)
	if err != nil {
		return nil, err
	}
	return &override, nil
}

// UpdateOnCallOverride updates an existing on-call override.
func (c *Client) UpdateOnCallOverride(scheduleID, overrideID string, req UpdateOnCallOverrideRequest) (*OnCallOverride, error) {
	var override OnCallOverride
	path := fmt.Sprintf("/v1/on-call/schedules/%s/overrides/%s", scheduleID, overrideID)
	err := c.doRequest(http.MethodPut, path, req, &override)
	if err != nil {
		return nil, err
	}
	return &override, nil
}

// DeleteOnCallOverride deletes an on-call override by ID.
func (c *Client) DeleteOnCallOverride(scheduleID, overrideID string) error {
	path := fmt.Sprintf("/v1/on-call/schedules/%s/overrides/%s", scheduleID, overrideID)
	return c.doRequest(http.MethodDelete, path, nil, nil)
}

// EscalationPolicy represents a Moneat escalation policy.
type EscalationPolicy struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"createdAt,omitempty"`
}

// CreateEscalationPolicyRequest is the request body for creating an escalation policy.
type CreateEscalationPolicyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// UpdateEscalationPolicyRequest is the request body for updating an escalation policy.
type UpdateEscalationPolicyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// GetEscalationPolicy retrieves an escalation policy by ID.
func (c *Client) GetEscalationPolicy(id string) (*EscalationPolicy, error) {
	var policy EscalationPolicy
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/on-call/escalation-policies/%s", id), nil, &policy)
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// ListEscalationPolicies retrieves all escalation policies.
func (c *Client) ListEscalationPolicies() ([]EscalationPolicy, error) {
	var policies []EscalationPolicy
	err := c.doRequest(http.MethodGet, "/v1/on-call/escalation-policies", nil, &policies)
	if err != nil {
		return nil, err
	}
	return policies, nil
}

// CreateEscalationPolicy creates a new escalation policy.
func (c *Client) CreateEscalationPolicy(req CreateEscalationPolicyRequest) (*EscalationPolicy, error) {
	var policy EscalationPolicy
	err := c.doRequest(http.MethodPost, "/v1/on-call/escalation-policies", req, &policy)
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// UpdateEscalationPolicy updates an existing escalation policy.
func (c *Client) UpdateEscalationPolicy(id string, req UpdateEscalationPolicyRequest) (*EscalationPolicy, error) {
	var policy EscalationPolicy
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/on-call/escalation-policies/%s", id), req, &policy)
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// DeleteEscalationPolicy deletes an escalation policy by ID.
func (c *Client) DeleteEscalationPolicy(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/on-call/escalation-policies/%s", id), nil, nil)
}

// OnCallPriorities represents on-call priority configuration.
type OnCallPriorities struct {
	Priorities []OnCallPriority `json:"priorities"`
}

// OnCallPriority represents a single priority level.
type OnCallPriority struct {
	Level    string `json:"level"`
	Label    string `json:"label"`
	Pageable bool   `json:"pageable"`
}

// UpdateOnCallPrioritiesRequest is the request body for updating on-call priorities.
type UpdateOnCallPrioritiesRequest struct {
	Priorities []OnCallPriority `json:"priorities"`
}

// GetOnCallPriorities retrieves on-call priority configuration.
func (c *Client) GetOnCallPriorities() (*OnCallPriorities, error) {
	var priorities OnCallPriorities
	err := c.doRequest(http.MethodGet, "/v1/priorities", nil, &priorities)
	if err != nil {
		return nil, err
	}
	return &priorities, nil
}

// UpdateOnCallPriorities updates on-call priority configuration.
func (c *Client) UpdateOnCallPriorities(req UpdateOnCallPrioritiesRequest) (*OnCallPriorities, error) {
	var priorities OnCallPriorities
	err := c.doRequest(http.MethodPut, "/v1/priorities", req, &priorities)
	if err != nil {
		return nil, err
	}
	return &priorities, nil
}

// BusinessHours represents business hours configuration.
type BusinessHours struct {
	Timezone string              `json:"timezone"`
	Enabled  bool                `json:"enabled"`
	Windows  []BusinessHoursSlot `json:"windows"`
}

// BusinessHoursSlot represents a business hours time window.
type BusinessHoursSlot struct {
	Day       string `json:"day"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// UpdateBusinessHoursRequest is the request body for updating business hours.
type UpdateBusinessHoursRequest struct {
	Timezone string              `json:"timezone"`
	Enabled  bool                `json:"enabled"`
	Windows  []BusinessHoursSlot `json:"windows"`
}

// GetBusinessHours retrieves business hours configuration.
func (c *Client) GetBusinessHours() (*BusinessHours, error) {
	var hours BusinessHours
	err := c.doRequest(http.MethodGet, "/v1/business-hours", nil, &hours)
	if err != nil {
		return nil, err
	}
	return &hours, nil
}

// UpdateBusinessHours updates business hours configuration.
func (c *Client) UpdateBusinessHours(req UpdateBusinessHoursRequest) (*BusinessHours, error) {
	var hours BusinessHours
	err := c.doRequest(http.MethodPut, "/v1/business-hours", req, &hours)
	if err != nil {
		return nil, err
	}
	return &hours, nil
}

// OnCallScheduleSlackUsergroup represents a Slack usergroup mapping for an on-call schedule.
type OnCallScheduleSlackUsergroup struct {
	ScheduleID      string `json:"scheduleId"`
	UsergroupID     string `json:"usergroupId"`
	UsergroupHandle string `json:"usergroupHandle"`
}

// SetOnCallScheduleSlackUsergroupRequest is the request body for setting a Slack usergroup mapping.
type SetOnCallScheduleSlackUsergroupRequest struct {
	UsergroupID     string `json:"usergroupId"`
	UsergroupHandle string `json:"usergroupHandle"`
}

// GetOnCallScheduleSlackUsergroup retrieves the Slack usergroup mapping for a schedule.
func (c *Client) GetOnCallScheduleSlackUsergroup(scheduleID string) (*OnCallScheduleSlackUsergroup, error) {
	var mapping OnCallScheduleSlackUsergroup
	path := fmt.Sprintf("/v1/on-call/schedules/%s/slack-usergroup", scheduleID)
	err := c.doRequest(http.MethodGet, path, nil, &mapping)
	if err != nil {
		return nil, err
	}
	return &mapping, nil
}

// SetOnCallScheduleSlackUsergroup sets the Slack usergroup mapping for a schedule.
func (c *Client) SetOnCallScheduleSlackUsergroup(
	scheduleID string, req SetOnCallScheduleSlackUsergroupRequest,
) (*OnCallScheduleSlackUsergroup, error) {
	var mapping OnCallScheduleSlackUsergroup
	path := fmt.Sprintf("/v1/on-call/schedules/%s/slack-usergroup", scheduleID)
	err := c.doRequest(http.MethodPut, path, req, &mapping)
	if err != nil {
		return nil, err
	}
	return &mapping, nil
}

// DeleteOnCallScheduleSlackUsergroup removes the Slack usergroup mapping for a schedule.
func (c *Client) DeleteOnCallScheduleSlackUsergroup(scheduleID string) error {
	path := fmt.Sprintf("/v1/on-call/schedules/%s/slack-usergroup", scheduleID)
	return c.doRequest(http.MethodDelete, path, nil, nil)
}
