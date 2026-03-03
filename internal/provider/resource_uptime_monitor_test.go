package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUptimeMonitorResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccUptimeMonitorResourceConfig("API Health", "https://api.example.com/health", "http", 60),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("moneat_uptime_monitor.test", "name", "API Health"),
					resource.TestCheckResourceAttr("moneat_uptime_monitor.test", "url", "https://api.example.com/health"),
					resource.TestCheckResourceAttr("moneat_uptime_monitor.test", "type", "http"),
					resource.TestCheckResourceAttr("moneat_uptime_monitor.test", "interval_seconds", "60"),
					resource.TestCheckResourceAttr("moneat_uptime_monitor.test", "paused", "false"),
					resource.TestCheckResourceAttrSet("moneat_uptime_monitor.test", "id"),
				),
			},
			// ImportState
			{
				ResourceName:      "moneat_uptime_monitor.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update
			{
				Config: testAccUptimeMonitorResourceConfig("API Health v2", "https://api.example.com/v2/health", "http", 30),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("moneat_uptime_monitor.test", "name", "API Health v2"),
					resource.TestCheckResourceAttr("moneat_uptime_monitor.test", "interval_seconds", "30"),
				),
			},
		},
	})
}

func testAccUptimeMonitorResourceConfig(name, url, monitorType string, interval int) string {
	return fmt.Sprintf(`
resource "moneat_uptime_monitor" "test" {
  name             = %[1]q
  url              = %[2]q
  type             = %[3]q
  interval_seconds = %[4]d
}
`, name, url, monitorType, interval)
}
