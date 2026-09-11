resource "moneat_connection_group" "pagerduty" {
  name                    = "PagerDuty by region"
  provider_id             = "pagerduty"
  use_id                  = "workflow_actions"
  member_installation_ids = var.member_installation_ids
  selection_strategy      = "first_match"
}
