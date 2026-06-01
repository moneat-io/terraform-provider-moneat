resource "moneat_workflow_connection" "jira" {
  type   = "jira"
  name   = "Jira production"
  secret = var.jira_api_token

  identifier_tags = {
    workspace = "production"
    region    = "us"
  }
}
