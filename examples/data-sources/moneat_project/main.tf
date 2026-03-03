data "moneat_project" "example" {
  id = "project-id-here"
}

output "project_name" {
  value = data.moneat_project.example.name
}
