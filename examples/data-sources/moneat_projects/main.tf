data "moneat_projects" "all" {}

output "project_count" {
  value = length(data.moneat_projects.all.projects)
}
