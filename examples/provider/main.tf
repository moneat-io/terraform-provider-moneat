terraform {
  required_providers {
    moneat = {
      source  = "moneat-io/moneat"
      version = "~> 0.1"
    }
  }
}

provider "moneat" {
  # base_url = "https://api.moneat.io"  # or set MONEAT_BASE_URL
  # token    = var.moneat_token          # or set MONEAT_AUTH_TOKEN
}
