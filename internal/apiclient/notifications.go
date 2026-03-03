package apiclient

import "net/http"

// NotificationPreferences represents global notification preferences.
type NotificationPreferences struct {
	IssueAlerts           bool  `json:"issueAlerts"`
	ErrorAlerts           bool  `json:"errorAlerts"`
	WeeklySummary         bool  `json:"weeklySummary"`
	AlertFrequencyMinutes int64 `json:"alertFrequencyMinutes"`
}

// UpdateNotificationPreferencesRequest is the request body for updating notification preferences.
type UpdateNotificationPreferencesRequest struct {
	IssueAlerts           bool  `json:"issueAlerts"`
	ErrorAlerts           bool  `json:"errorAlerts"`
	WeeklySummary         bool  `json:"weeklySummary"`
	AlertFrequencyMinutes int64 `json:"alertFrequencyMinutes"`
}

// AlertNotificationChannels represents per-source notification channel preferences.
type AlertNotificationChannels struct {
	AlertSource    string `json:"alertSource"`
	EmailEnabled   bool   `json:"emailEnabled"`
	SlackEnabled   bool   `json:"slackEnabled"`
	DiscordEnabled bool   `json:"discordEnabled"`
}

// UpdateAlertNotificationChannelsRequest is the request body for updating channel preferences.
type UpdateAlertNotificationChannelsRequest struct {
	AlertSource    string `json:"alertSource"`
	EmailEnabled   bool   `json:"emailEnabled"`
	SlackEnabled   bool   `json:"slackEnabled"`
	DiscordEnabled bool   `json:"discordEnabled"`
}

// GetNotificationPreferences retrieves global notification preferences.
func (c *Client) GetNotificationPreferences() (*NotificationPreferences, error) {
	var prefs NotificationPreferences
	err := c.doRequest(http.MethodGet, "/v1/notification-preferences", nil, &prefs)
	if err != nil {
		return nil, err
	}
	return &prefs, nil
}

// UpdateNotificationPreferences updates global notification preferences.
func (c *Client) UpdateNotificationPreferences(req UpdateNotificationPreferencesRequest) (*NotificationPreferences, error) {
	var prefs NotificationPreferences
	err := c.doRequest(http.MethodPut, "/v1/notification-preferences", req, &prefs)
	if err != nil {
		return nil, err
	}
	return &prefs, nil
}

// GetAlertNotificationChannels retrieves alert notification channel preferences.
func (c *Client) GetAlertNotificationChannels() ([]AlertNotificationChannels, error) {
	var channels []AlertNotificationChannels
	err := c.doRequest(http.MethodGet, "/v1/alert-notification-preferences", nil, &channels)
	if err != nil {
		return nil, err
	}
	return channels, nil
}

// UpdateAlertNotificationChannels updates alert notification channel preferences.
func (c *Client) UpdateAlertNotificationChannels(req UpdateAlertNotificationChannelsRequest) (*AlertNotificationChannels, error) {
	var channels AlertNotificationChannels
	err := c.doRequest(http.MethodPut, "/v1/alert-notification-preferences", req, &channels)
	if err != nil {
		return nil, err
	}
	return &channels, nil
}
