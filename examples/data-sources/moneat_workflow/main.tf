data "moneat_workflow" "critical_alerts" {
  id = "123"
}

output "workflow_graph" {
  value = data.moneat_workflow.critical_alerts.graph_json
}
