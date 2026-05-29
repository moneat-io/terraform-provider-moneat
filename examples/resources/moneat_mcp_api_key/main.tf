resource "moneat_mcp_api_key" "project_reader" {
  name          = "Project reader"
  enabled_tools = ["list_projects", "get_project"]

  enabled_resources = []
  expires_in_days   = 90
}
