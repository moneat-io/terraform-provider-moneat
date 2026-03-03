package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationPreferencesResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: `
resource "moneat_notification_preferences" "test" {
  issue_alerts            = true
  error_alerts            = false
  weekly_summary          = true
  alert_frequency_minutes = 30
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("moneat_notification_preferences.test", "issue_alerts", "true"),
					resource.TestCheckResourceAttr("moneat_notification_preferences.test", "error_alerts", "false"),
					resource.TestCheckResourceAttr("moneat_notification_preferences.test", "weekly_summary", "true"),
					resource.TestCheckResourceAttr("moneat_notification_preferences.test", "alert_frequency_minutes", "30"),
				),
			},
			// Update
			{
				Config: `
resource "moneat_notification_preferences" "test" {
  issue_alerts            = false
  error_alerts            = true
  weekly_summary          = false
  alert_frequency_minutes = 60
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("moneat_notification_preferences.test", "issue_alerts", "false"),
					resource.TestCheckResourceAttr("moneat_notification_preferences.test", "error_alerts", "true"),
					resource.TestCheckResourceAttr("moneat_notification_preferences.test", "alert_frequency_minutes", "60"),
				),
			},
		},
	})
}

func TestAccAlertNotificationChannelsResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: `
resource "moneat_alert_notification_channels" "test" {
  alert_source    = "system_alert"
  email_enabled   = true
  slack_enabled   = true
  discord_enabled = false
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("moneat_alert_notification_channels.test", "alert_source", "system_alert"),
					resource.TestCheckResourceAttr("moneat_alert_notification_channels.test", "email_enabled", "true"),
					resource.TestCheckResourceAttr("moneat_alert_notification_channels.test", "slack_enabled", "true"),
					resource.TestCheckResourceAttr("moneat_alert_notification_channels.test", "discord_enabled", "false"),
				),
			},
			// Update
			{
				Config: `
resource "moneat_alert_notification_channels" "test" {
  alert_source    = "system_alert"
  email_enabled   = false
  slack_enabled   = false
  discord_enabled = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("moneat_alert_notification_channels.test", "email_enabled", "false"),
					resource.TestCheckResourceAttr("moneat_alert_notification_channels.test", "discord_enabled", "true"),
				),
			},
		},
	})
}
