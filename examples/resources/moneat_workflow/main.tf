resource "moneat_workflow" "critical_alert_email" {
  name         = "Critical alert email"
  trigger_name = "alert.triggered"
  enabled      = true
  published    = true

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

resource "moneat_workflow" "security_signal_triage" {
  name         = "Security signal triage"
  trigger_name = "security.signal"
  enabled      = true
  published    = true

  graph_json = jsonencode({
    nodes = [
      {
        id      = "trigger"
        type    = "trigger"
        trigger = "security.signal"
        position = {
          x = 0
          y = 0
        }
      },
      {
        id     = "notify"
        type   = "action"
        action = "notification.slack"
        params = {
          channel = "#security"
          message = "Security signal: {{signal.title}}"
        }
        position = {
          x = 320
          y = 0
        }
      }
    ]
    edges = [
      {
        from = "trigger"
        to   = "notify"
      }
    ]
  })
}
