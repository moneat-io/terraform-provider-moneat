resource "moneat_feature_flag_environment" "qa" {
  key  = "qa"
  name = "QA"
}

resource "moneat_feature_flag_sdk_key" "server" {
  environment_key = moneat_feature_flag_environment.qa.key
  name            = "QA server SDK key"
  key_type        = "server"
}
