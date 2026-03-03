package apiclient

import (
	"fmt"
	"net/http"
)

// Host represents a monitored host in Moneat.
type Host struct {
	ID       string `json:"id"`
	Hostname string `json:"hostname"`
	OS       string `json:"os,omitempty"`
	Platform string `json:"platform,omitempty"`
	Status   string `json:"status,omitempty"`
}

// GetHost retrieves a host by ID.
func (c *Client) GetHost(id string) (*Host, error) {
	var host Host
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/monitor/hosts/%s", id), nil, &host)
	if err != nil {
		return nil, err
	}
	return &host, nil
}

// ListHosts retrieves all monitored hosts.
func (c *Client) ListHosts() ([]Host, error) {
	var hosts []Host
	err := c.doRequest(http.MethodGet, "/v1/monitor/hosts", nil, &hosts)
	if err != nil {
		return nil, err
	}
	return hosts, nil
}
