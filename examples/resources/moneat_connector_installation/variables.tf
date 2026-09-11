variable "pagerduty_api_key" {
  description = "PagerDuty API key used only when creating or rotating the installation."
  type        = string
  sensitive   = true
}
