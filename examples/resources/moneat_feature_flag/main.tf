resource "moneat_feature_flag" "checkout_flow" {
  key            = "checkout-flow"
  name           = "Checkout Flow"
  description    = "Controls the new checkout flow"
  value_type     = "BOOLEAN"
  client_visible = false
  tags           = ["checkout", "release"]

  variants_json = jsonencode([
    {
      key   = "off"
      name  = "Off"
      value = false
    },
    {
      key   = "on"
      name  = "On"
      value = true
    }
  ])
}
