resource "moneat_feature_flag_segment" "enterprise_accounts" {
  key         = "enterprise-accounts"
  name        = "Enterprise Accounts"
  description = "Accounts with enterprise entitlements"

  conditions_json = jsonencode({
    all = [
      {
        attribute = "plan"
        op        = "equals"
        value     = "enterprise"
      }
    ]
  })
}
