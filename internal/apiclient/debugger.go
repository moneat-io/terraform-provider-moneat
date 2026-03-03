package apiclient

import (
	"fmt"
	"net/http"
)

// DebuggerProbe represents a Moneat debugger probe.
type DebuggerProbe struct {
	ID          string `json:"id"`
	ServiceName string `json:"serviceName"`
	FilePath    string `json:"filePath"`
	LineNumber  int64  `json:"lineNumber"`
	Expression  string `json:"expression"`
	Enabled     bool   `json:"enabled"`
}

// CreateDebuggerProbeRequest is the request body for creating a debugger probe.
type CreateDebuggerProbeRequest struct {
	ServiceName string `json:"serviceName"`
	FilePath    string `json:"filePath"`
	LineNumber  int64  `json:"lineNumber"`
	Expression  string `json:"expression"`
	Enabled     bool   `json:"enabled"`
}

// UpdateDebuggerProbeRequest is the request body for updating a debugger probe.
type UpdateDebuggerProbeRequest struct {
	ServiceName string `json:"serviceName"`
	FilePath    string `json:"filePath"`
	LineNumber  int64  `json:"lineNumber"`
	Expression  string `json:"expression"`
	Enabled     bool   `json:"enabled"`
}

// GetDebuggerProbe retrieves a debugger probe by ID.
func (c *Client) GetDebuggerProbe(id string) (*DebuggerProbe, error) {
	var probe DebuggerProbe
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/infra/debugger/probes/%s", id), nil, &probe)
	if err != nil {
		return nil, err
	}
	return &probe, nil
}

// ListDebuggerProbes retrieves all debugger probes.
func (c *Client) ListDebuggerProbes() ([]DebuggerProbe, error) {
	var probes []DebuggerProbe
	err := c.doRequest(http.MethodGet, "/v1/infra/debugger/probes", nil, &probes)
	if err != nil {
		return nil, err
	}
	return probes, nil
}

// CreateDebuggerProbe creates a new debugger probe.
func (c *Client) CreateDebuggerProbe(req CreateDebuggerProbeRequest) (*DebuggerProbe, error) {
	var probe DebuggerProbe
	err := c.doRequest(http.MethodPost, "/v1/infra/debugger/probes", req, &probe)
	if err != nil {
		return nil, err
	}
	return &probe, nil
}

// UpdateDebuggerProbe updates an existing debugger probe.
func (c *Client) UpdateDebuggerProbe(id string, req UpdateDebuggerProbeRequest) (*DebuggerProbe, error) {
	var probe DebuggerProbe
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/infra/debugger/probes/%s", id), req, &probe)
	if err != nil {
		return nil, err
	}
	return &probe, nil
}

// DeleteDebuggerProbe deletes a debugger probe by ID.
func (c *Client) DeleteDebuggerProbe(id string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/infra/debugger/probes/%s", id), nil, nil)
}
