resource "moneat_workflow" "critical_alert_email" {
  name         = "Critical alert email"
  trigger_name = "alert.triggered"
  enabled      = true

  conditions_json = jsonencode([
    {
      reference = "alert.severity"
      operation = "eq"
      value     = "CRITICAL"
    }
  ])

  steps_json = jsonencode([
    {
      name = "notification.email_org"
      params = {
        subject = "Critical alert: {{alert.title}}"
        body    = "{{alert.description}}"
      }
    }
  ])

  once_for_template = ["alert.deduplication_key"]
}
