data "moneat_workflows" "all" {}

output "workflow_names" {
  value = [for workflow in data.moneat_workflows.all.workflows : workflow.name]
}
