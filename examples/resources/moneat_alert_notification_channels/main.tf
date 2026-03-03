resource "moneat_alert_notification_channels" "example" {
  alert_source    = "system_alert"
  email_enabled   = true
  slack_enabled   = true
  discord_enabled = false
}
