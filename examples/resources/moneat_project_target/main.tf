resource "moneat_project" "backend" {
  name     = "Backend"
  platform = "kotlin"
}

resource "moneat_project_target" "api" {
  project_id = moneat_project.backend.id
  target     = "api"
}
