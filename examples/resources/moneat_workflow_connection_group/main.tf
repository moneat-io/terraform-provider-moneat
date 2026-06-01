resource "moneat_workflow_connection" "jira_us" {
  type   = "jira"
  name   = "Jira US"
  secret = var.jira_us_api_token

  identifier_tags = {
    region = "us"
  }
}

resource "moneat_workflow_connection" "jira_eu" {
  type   = "jira"
  name   = "Jira EU"
  secret = var.jira_eu_api_token

  identifier_tags = {
    region = "eu"
  }
}

resource "moneat_workflow_connection_group" "jira_by_region" {
  name            = "Jira by region"
  connection_type = "jira"
  member_connection_ids = [
    tonumber(moneat_workflow_connection.jira_us.id),
    tonumber(moneat_workflow_connection.jira_eu.id),
  ]
  selection_strategy = "first_match"
}
