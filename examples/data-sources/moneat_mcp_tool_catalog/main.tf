data "moneat_mcp_tool_catalog" "current" {}

output "mcp_tool_catalog_json" {
  value = data.moneat_mcp_tool_catalog.current.json
}
