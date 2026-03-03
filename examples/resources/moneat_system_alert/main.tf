resource "moneat_system_alert" "example" {
  system_id        = "abc-123"
  metric           = "cpu"
  condition        = "gt"
  threshold        = 90
  duration_seconds = 300
  enabled          = true
}
