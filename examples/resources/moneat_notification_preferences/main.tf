resource "moneat_notification_preferences" "example" {
  issue_alerts            = true
  error_alerts            = true
  weekly_summary          = true
  alert_frequency_minutes = 15
}
