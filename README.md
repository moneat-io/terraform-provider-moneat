<p align="center">
  <a href="https://moneat.io">
    <img alt="Moneat" src="moneat-logo.svg" width="280">
  </a>
</p>

<p align="center">
  <em>Terraform Provider for <a href="https://github.com/moneat-io/moneat">Moneat</a> — the open-source observability platform</em>
</p>

<p align="center">
  <a href="https://registry.terraform.io/providers/moneat-io/moneat/latest"><img src="https://img.shields.io/badge/Terraform%20Registry-moneat--io%2Fmoneat-blueviolet?style=flat-square&logo=terraform" alt="Terraform Registry"></a>
  <a href="https://github.com/moneat-io/terraform-provider-moneat/releases"><img src="https://img.shields.io/github/v/release/moneat-io/terraform-provider-moneat?style=flat-square" alt="Latest Release"></a>
  <a href="https://github.com/moneat-io/terraform-provider-moneat/actions"><img src="https://img.shields.io/github/actions/workflow/status/moneat-io/terraform-provider-moneat/ci.yml?style=flat-square&label=CI" alt="CI Status"></a>
  <a href="LICENSE"><img src="https://img.shields.io/badge/License-MPL--2.0-blue.svg?style=flat-square" alt="License: MPL-2.0"></a>
</p>

---

# Terraform Provider for Moneat

The official [Terraform](https://www.terraform.io/) provider for [Moneat](https://github.com/moneat-io/moneat), enabling infrastructure-as-code management of your observability platform.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23 (to build the provider plugin)

## Usage

```hcl
terraform {
  required_providers {
    moneat = {
      source  = "moneat-io/moneat"
      version = "~> 0.1"
    }
  }
}

provider "moneat" {
  # base_url = "https://api.moneat.io"  # or MONEAT_BASE_URL env var
  # token    = var.moneat_token          # or MONEAT_AUTH_TOKEN env var
}

resource "moneat_project" "backend" {
  name     = "backend-api"
  platform = "python"
}

resource "moneat_uptime_monitor" "api_health" {
  name             = "API Health Check"
  url              = "https://api.example.com/health"
  type             = "http"
  interval_seconds = 60
}
```

## Authentication

The provider authenticates using a bearer token. You can provide it in two ways:

1. **Provider attribute**: Set `token` in the provider block
2. **Environment variable**: Set `MONEAT_AUTH_TOKEN`

Generate an auth token in Moneat under **Settings → Auth Tokens**, or via the API:

```bash
curl -X POST https://api.moneat.io/v1/auth-tokens \
  -H "Authorization: Bearer <your-session-token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "terraform"}'
```

## Resources

### Phase 1 (Core)

| Resource | Description |
|----------|-------------|
| `moneat_project` | Manage projects |
| `moneat_uptime_monitor` | HTTP, TCP, ping, and push uptime monitors |
| `moneat_system_alert` | Threshold-based host metric alerts |
| `moneat_status_page` | Public/private status pages |
| `moneat_dashboard` | Custom dashboards |
| `moneat_notification_preferences` | Global notification settings |
| `moneat_alert_notification_channels` | Per-source notification channel preferences |

### Additional API Coverage

| Resource | Description |
|----------|-------------|
| `moneat_feature_flag` | Feature flag definitions and variants |
| `moneat_feature_flag_config` | Per-environment feature flag rollout config |
| `moneat_feature_flag_environment` | Feature flag environments |
| `moneat_feature_flag_segment` | Feature flag targeting segments |
| `moneat_feature_flag_sdk_key` | Feature flag SDK keys |
| `moneat_workflow` | Alert and automation workflows |
| `moneat_workflow_connection` | Vaulted workflow connector credentials |
| `moneat_workflow_connection_group` | Workflow connector routing groups |
| `moneat_security_detection_rule` | Security detection rules |
| `moneat_synthetic_variable` | Reusable synthetic test variables |
| `moneat_mcp_api_key` | MCP API keys and tool/resource permissions |
| `moneat_otlp_service_mapping` | Telemetry service-to-project routing |
| `moneat_project_target` | Additional project DSN targets |

### Data Sources

| Data Source | Description |
|-------------|-------------|
| `moneat_project` | Look up a project by ID |
| `moneat_projects` | List all projects |
| `moneat_workflow` | Look up a workflow by ID |
| `moneat_workflows` | List all workflows |
| `moneat_mcp_tool_catalog` | Read the MCP tool and resource catalog |
| `moneat_security_signals` | Read security signal triage data as JSON |
| `moneat_security_detection_coverage` | Read security detection coverage as JSON |
| `moneat_security_vulnerability_summary` | Read vulnerability summary data as JSON |
| `moneat_security_vulnerability_findings` | Read vulnerability findings as JSON |
| `moneat_security_vulnerability_inventory` | Read vulnerability package inventory as JSON |

## Configuration Reference

| Attribute | Description | Default | Environment Variable |
|-----------|-------------|---------|---------------------|
| `base_url` | Moneat API base URL | `https://api.moneat.io` | `MONEAT_BASE_URL` |
| `token` | API authentication token | — | `MONEAT_AUTH_TOKEN` |

## Development

### Building

```bash
make build
```

### Testing

Acceptance tests run against a real Moneat instance:

```bash
export MONEAT_AUTH_TOKEN="your-token"
export MONEAT_BASE_URL="http://localhost:8080"
make testacc
```

### Installing Locally

```bash
make install
```

### Generating Documentation

```bash
make docs
```

## License

[Mozilla Public License v2.0](LICENSE)

---

<p align="center">
  <a href="https://github.com/moneat-io/moneat">
    <img alt="Moneat" src="moneat-icon.svg" width="32">
  </a>
  <br>
  <sub>Built for <a href="https://github.com/moneat-io/moneat">Moneat</a> — the open-source observability platform</sub>
</p>
