resource "moneat_connector_installation" "pagerduty" {
  provider_id           = "pagerduty"
  auth_profile_id       = "workflow_api_key"
  use_id                = "workflow_actions"
  name                  = "PagerDuty response"
  external_account_json = jsonencode({})
  secret                = var.pagerduty_api_key
  identifier_tags       = { region = "us" }
}
