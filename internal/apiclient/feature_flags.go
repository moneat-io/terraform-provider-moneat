package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// FeatureFlagEnvironment represents a feature flag environment.
type FeatureFlagEnvironment struct {
	ID          int    `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Version     int    `json:"version"`
	CreatedAt   string `json:"createdAt,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
}

// CreateFeatureFlagEnvironmentRequest is the request body for creating an environment.
type CreateFeatureFlagEnvironmentRequest struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type featureFlagEnvironmentsResponse struct {
	Environments []FeatureFlagEnvironment `json:"environments"`
}

// ListFeatureFlagEnvironments retrieves all feature flag environments.
func (c *Client) ListFeatureFlagEnvironments() ([]FeatureFlagEnvironment, error) {
	var response featureFlagEnvironmentsResponse
	err := c.doRequest(http.MethodGet, "/v1/feature-flags/environments", nil, &response)
	if err != nil {
		return nil, err
	}
	return response.Environments, nil
}

// GetFeatureFlagEnvironment retrieves a feature flag environment by key.
func (c *Client) GetFeatureFlagEnvironment(key string) (*FeatureFlagEnvironment, error) {
	environments, err := c.ListFeatureFlagEnvironments()
	if err != nil {
		return nil, err
	}
	for _, environment := range environments {
		if environment.Key == key {
			return &environment, nil
		}
	}
	return nil, notFoundError("Feature flag environment not found")
}

// CreateFeatureFlagEnvironment creates a feature flag environment.
func (c *Client) CreateFeatureFlagEnvironment(
	req CreateFeatureFlagEnvironmentRequest,
) (*FeatureFlagEnvironment, error) {
	var environment FeatureFlagEnvironment
	err := c.doRequest(http.MethodPost, "/v1/feature-flags/environments", req, &environment)
	if err != nil {
		return nil, err
	}
	return &environment, nil
}

// FeatureFlagSegment represents a feature flag segment.
type FeatureFlagSegment struct {
	ID          int             `json:"id"`
	Key         string          `json:"key"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Conditions  json.RawMessage `json:"conditions"`
	CreatedAt   string          `json:"createdAt,omitempty"`
	UpdatedAt   string          `json:"updatedAt,omitempty"`
}

// FeatureFlagSegmentRequest is the request body for creating or updating a segment.
type FeatureFlagSegmentRequest struct {
	Key         string          `json:"key"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Conditions  json.RawMessage `json:"conditions"`
}

type featureFlagSegmentsResponse struct {
	Segments []FeatureFlagSegment `json:"segments"`
}

// ListFeatureFlagSegments retrieves all feature flag segments.
func (c *Client) ListFeatureFlagSegments() ([]FeatureFlagSegment, error) {
	var response featureFlagSegmentsResponse
	err := c.doRequest(http.MethodGet, "/v1/feature-flags/segments", nil, &response)
	if err != nil {
		return nil, err
	}
	return response.Segments, nil
}

// GetFeatureFlagSegment retrieves a feature flag segment by key.
func (c *Client) GetFeatureFlagSegment(key string) (*FeatureFlagSegment, error) {
	segments, err := c.ListFeatureFlagSegments()
	if err != nil {
		return nil, err
	}
	for _, segment := range segments {
		if segment.Key == key {
			return &segment, nil
		}
	}
	return nil, notFoundError("Feature flag segment not found")
}

// UpsertFeatureFlagSegment creates or updates a feature flag segment.
func (c *Client) UpsertFeatureFlagSegment(req FeatureFlagSegmentRequest) (*FeatureFlagSegment, error) {
	var segment FeatureFlagSegment
	err := c.doRequest(http.MethodPost, "/v1/feature-flags/segments", req, &segment)
	if err != nil {
		return nil, err
	}
	return &segment, nil
}

// DeleteFeatureFlagSegment deletes a feature flag segment.
func (c *Client) DeleteFeatureFlagSegment(key string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/feature-flags/segments/%s", key), nil, nil)
}

// FeatureFlagSdkKey represents a feature flag SDK key.
type FeatureFlagSdkKey struct {
	ID             int    `json:"id"`
	EnvironmentKey string `json:"environmentKey"`
	Name           string `json:"name"`
	KeyType        string `json:"keyType"`
	KeyPrefix      string `json:"keyPrefix"`
	Key            string `json:"key,omitempty"`
	CreatedAt      string `json:"createdAt,omitempty"`
	LastUsedAt     string `json:"lastUsedAt,omitempty"`
}

// FeatureFlagSdkKeyRequest is the request body for creating a feature flag SDK key.
type FeatureFlagSdkKeyRequest struct {
	EnvironmentKey string `json:"environmentKey"`
	Name           string `json:"name"`
	KeyType        string `json:"keyType"`
}

type featureFlagSdkKeysResponse struct {
	Keys []FeatureFlagSdkKey `json:"keys"`
}

// ListFeatureFlagSdkKeys retrieves all feature flag SDK keys.
func (c *Client) ListFeatureFlagSdkKeys() ([]FeatureFlagSdkKey, error) {
	var response featureFlagSdkKeysResponse
	err := c.doRequest(http.MethodGet, "/v1/feature-flags/sdk-keys", nil, &response)
	if err != nil {
		return nil, err
	}
	return response.Keys, nil
}

// GetFeatureFlagSdkKey retrieves a feature flag SDK key by ID.
func (c *Client) GetFeatureFlagSdkKey(id int) (*FeatureFlagSdkKey, error) {
	keys, err := c.ListFeatureFlagSdkKeys()
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		if key.ID == id {
			return &key, nil
		}
	}
	return nil, notFoundError("Feature flag SDK key not found")
}

// CreateFeatureFlagSdkKey creates a feature flag SDK key.
func (c *Client) CreateFeatureFlagSdkKey(req FeatureFlagSdkKeyRequest) (*FeatureFlagSdkKey, error) {
	var key FeatureFlagSdkKey
	err := c.doRequest(http.MethodPost, "/v1/feature-flags/sdk-keys", req, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// DeleteFeatureFlagSdkKey revokes a feature flag SDK key.
func (c *Client) DeleteFeatureFlagSdkKey(id int) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/feature-flags/sdk-keys/%d", id), nil, nil)
}

// FeatureFlagVariantRequest is the request body shape for a feature flag variant.
type FeatureFlagVariantRequest struct {
	Key   string          `json:"key"`
	Name  string          `json:"name,omitempty"`
	Value json.RawMessage `json:"value"`
}

// FeatureFlagVariant represents a feature flag variant.
type FeatureFlagVariant struct {
	ID        int             `json:"id"`
	Key       string          `json:"key"`
	Name      string          `json:"name"`
	Value     json.RawMessage `json:"value"`
	SortOrder int             `json:"sortOrder"`
}

// FeatureFlagConfig represents a per-environment feature flag configuration.
type FeatureFlagConfig struct {
	EnvironmentKey    string          `json:"environmentKey"`
	EnvironmentName   string          `json:"environmentName"`
	Enabled           bool            `json:"enabled"`
	DefaultVariantKey string          `json:"defaultVariantKey,omitempty"`
	OffVariantKey     string          `json:"offVariantKey,omitempty"`
	Rules             json.RawMessage `json:"rules"`
	Version           int             `json:"version"`
	UpdatedAt         string          `json:"updatedAt,omitempty"`
}

// FeatureFlag represents a feature flag.
type FeatureFlag struct {
	ID            int                  `json:"id"`
	Key           string               `json:"key"`
	Name          string               `json:"name"`
	Description   string               `json:"description,omitempty"`
	ValueType     string               `json:"valueType"`
	ClientVisible bool                 `json:"clientVisible"`
	Tags          []string             `json:"tags"`
	Variants      []FeatureFlagVariant `json:"variants"`
	Configs       []FeatureFlagConfig  `json:"configs"`
	CreatedAt     string               `json:"createdAt,omitempty"`
	UpdatedAt     string               `json:"updatedAt,omitempty"`
}

// CreateFeatureFlagRequest is the request body for creating a feature flag.
type CreateFeatureFlagRequest struct {
	Key               string                      `json:"key"`
	Name              string                      `json:"name"`
	Description       string                      `json:"description,omitempty"`
	ValueType         string                      `json:"valueType"`
	ClientVisible     bool                        `json:"clientVisible"`
	Tags              []string                    `json:"tags"`
	Variants          []FeatureFlagVariantRequest `json:"variants"`
	DefaultVariantKey string                      `json:"defaultVariantKey,omitempty"`
	OffVariantKey     string                      `json:"offVariantKey,omitempty"`
}

// UpdateFeatureFlagRequest is the request body for updating a feature flag.
type UpdateFeatureFlagRequest struct {
	Name          string                      `json:"name,omitempty"`
	Description   string                      `json:"description,omitempty"`
	ClientVisible *bool                       `json:"clientVisible,omitempty"`
	Tags          []string                    `json:"tags"`
	Variants      []FeatureFlagVariantRequest `json:"variants"`
}

// UpdateFeatureFlagConfigRequest is the request body for updating a flag config.
type UpdateFeatureFlagConfigRequest struct {
	Enabled           *bool           `json:"enabled,omitempty"`
	DefaultVariantKey string          `json:"defaultVariantKey,omitempty"`
	OffVariantKey     string          `json:"offVariantKey,omitempty"`
	Rules             json.RawMessage `json:"rules,omitempty"`
}

// GetFeatureFlag retrieves a feature flag by key.
func (c *Client) GetFeatureFlag(key string) (*FeatureFlag, error) {
	var flag FeatureFlag
	err := c.doRequest(http.MethodGet, fmt.Sprintf("/v1/feature-flags/%s", key), nil, &flag)
	if err != nil {
		return nil, err
	}
	return &flag, nil
}

// CreateFeatureFlag creates a feature flag.
func (c *Client) CreateFeatureFlag(req CreateFeatureFlagRequest) (*FeatureFlag, error) {
	var flag FeatureFlag
	err := c.doRequest(http.MethodPost, "/v1/feature-flags", req, &flag)
	if err != nil {
		return nil, err
	}
	return &flag, nil
}

// UpdateFeatureFlag updates a feature flag.
func (c *Client) UpdateFeatureFlag(key string, req UpdateFeatureFlagRequest) (*FeatureFlag, error) {
	var flag FeatureFlag
	err := c.doRequest(http.MethodPut, fmt.Sprintf("/v1/feature-flags/%s", key), req, &flag)
	if err != nil {
		return nil, err
	}
	return &flag, nil
}

// DeleteFeatureFlag archives a feature flag.
func (c *Client) DeleteFeatureFlag(key string) error {
	return c.doRequest(http.MethodDelete, fmt.Sprintf("/v1/feature-flags/%s", key), nil, nil)
}

// UpdateFeatureFlagConfig updates a feature flag configuration for one environment.
func (c *Client) UpdateFeatureFlagConfig(
	flagKey string,
	environmentKey string,
	req UpdateFeatureFlagConfigRequest,
) (*FeatureFlagConfig, error) {
	var config FeatureFlagConfig
	path := fmt.Sprintf("/v1/feature-flags/%s/config/%s", flagKey, environmentKey)
	err := c.doRequest(http.MethodPut, path, req, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetFeatureFlagConfig retrieves a feature flag configuration from the flag response.
func (c *Client) GetFeatureFlagConfig(flagKey, environmentKey string) (*FeatureFlagConfig, error) {
	flag, err := c.GetFeatureFlag(flagKey)
	if err != nil {
		return nil, err
	}
	for _, config := range flag.Configs {
		if config.EnvironmentKey == environmentKey {
			return &config, nil
		}
	}
	return nil, notFoundError("Feature flag config not found")
}
