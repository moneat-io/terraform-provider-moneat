resource "moneat_project" "checkout" {
  name     = "Checkout"
  platform = "javascript"
}

resource "moneat_otlp_service_mapping" "checkout_api" {
  service_namespace   = "commerce"
  service_name        = "checkout-api"
  project_resource_id = moneat_project.checkout.id
}
