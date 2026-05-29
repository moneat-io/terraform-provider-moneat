resource "moneat_feature_flag_environment" "qa" {
  key  = "qa"
  name = "QA"
}

resource "moneat_feature_flag" "checkout_flow" {
  key        = "checkout-flow"
  name       = "Checkout Flow"
  value_type = "BOOLEAN"

  variants_json = jsonencode([
    {
      key   = "off"
      value = false
    },
    {
      key   = "on"
      value = true
    }
  ])
}

resource "moneat_feature_flag_config" "qa" {
  flag_key            = moneat_feature_flag.checkout_flow.key
  environment_key     = moneat_feature_flag_environment.qa.key
  enabled             = true
  default_variant_key = "on"
  off_variant_key     = "off"
  rules_json          = jsonencode({ rules = [] })
}
