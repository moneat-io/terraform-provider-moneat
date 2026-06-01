data "moneat_security_signals" "open_critical" {
  status   = "open"
  severity = "critical"
  limit    = 25
}

output "open_critical_security_signals_json" {
  value = data.moneat_security_signals.open_critical.json
}
