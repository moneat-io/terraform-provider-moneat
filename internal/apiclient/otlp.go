package apiclient

import (
	"fmt"
	"net/http"
)

// OtlpServiceMapping represents a telemetry service-to-project mapping.
type OtlpServiceMapping struct {
	ID                int    `json:"id"`
	ServiceNamespace  string `json:"service_namespace"`
	ServiceName       string `json:"service_name"`
	ProjectID         int64  `json:"project_id"`
	ProjectResourceID string `json:"project_resource_id"`
	ProjectName       string `json:"project_name"`
	UpdatedAt         string `json:"updated_at,omitempty"`
}

// CreateOtlpServiceMappingRequest is the request body for upserting a service mapping.
type CreateOtlpServiceMappingRequest struct {
	ServiceNamespace  string `json:"service_namespace,omitempty"`
	ServiceName       string `json:"service_name"`
	ProjectID         int64  `json:"project_id,omitempty"`
	ProjectResourceID string `json:"project_resource_id,omitempty"`
}

// OtlpObservedService represents an observed telemetry service.
type OtlpObservedService struct {
	ID                int    `json:"id"`
	MappingID         int    `json:"mapping_id,omitempty"`
	ServiceNamespace  string `json:"service_namespace"`
	ServiceName       string `json:"service_name"`
	ProjectID         int64  `json:"project_id,omitempty"`
	ProjectResourceID string `json:"project_resource_id,omitempty"`
	ProjectName       string `json:"project_name,omitempty"`
	SeenLogs          bool   `json:"seen_logs"`
	SeenTraces        bool   `json:"seen_traces"`
	SeenMetrics       bool   `json:"seen_metrics"`
	LastEnvironment   string `json:"last_environment,omitempty"`
	FirstSeenAt       string `json:"first_seen_at,omitempty"`
	LastSeenAt        string `json:"last_seen_at,omitempty"`
}

type otlpObservedServicesResponse struct {
	Services []OtlpObservedService `json:"services"`
}

// ListOtlpObservedServices retrieves observed OTLP services.
func (c *Client) ListOtlpObservedServices() ([]OtlpObservedService, error) {
	var response otlpObservedServicesResponse
	err := c.doRequest(http.MethodGet, "/v1/otlp/services", nil, &response)
	if err != nil {
		return nil, err
	}
	return response.Services, nil
}

// FindOtlpServiceMappingFromObserved finds a mapping when its service has been observed.
func (c *Client) FindOtlpServiceMappingFromObserved(id int) (*OtlpServiceMapping, error) {
	services, err := c.ListOtlpObservedServices()
	if err != nil {
		return nil, err
	}
	for _, service := range services {
		if service.MappingID == id {
			return &OtlpServiceMapping{
				ID:                service.MappingID,
				ServiceNamespace:  service.ServiceNamespace,
				ServiceName:       service.ServiceName,
				ProjectID:         service.ProjectID,
				ProjectResourceID: service.ProjectResourceID,
				ProjectName:       service.ProjectName,
			}, nil
		}
	}
	return nil, notFoundError("OTLP service mapping not found")
}

// UpsertOtlpServiceMapping creates or updates a service mapping.
func (c *Client) UpsertOtlpServiceMapping(req CreateOtlpServiceMappingRequest) (*OtlpServiceMapping, error) {
	var mapping OtlpServiceMapping
	err := c.doRequest(http.MethodPost, "/v1/otlp/service-mappings", req, &mapping)
	if err != nil {
		return nil, err
	}
	return &mapping, nil
}

// DeleteOtlpServiceMapping deletes a service mapping.
func (c *Client) DeleteOtlpServiceMapping(id int) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/otlp/service-mappings/%d", id), nil, nil)
}
