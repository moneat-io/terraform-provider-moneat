resource "moneat_uptime_monitor" "example" {
  name             = "API Health Check"
  url              = "https://api.example.com/health"
  type             = "http"
  interval_seconds = 60
}
