resource "moneat_security_detection_rule" "repeated_failed_logins" {
  name            = "Repeated failed logins"
  description     = "Flags repeated authentication failures by account."
  source          = "logs"
  filter          = "message contains \"failed login\""
  group_by        = ["user.email", "host.name"]
  window_seconds  = 300
  type            = "threshold"
  threshold_count = 5
  severity        = "high"
  signal_title    = "Repeated failed logins"
  signal_message  = "{{user.email}} had repeated failed login attempts."
  enabled         = true
  tags            = ["auth", "security"]
}
