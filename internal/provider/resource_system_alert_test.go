package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSystemAlertResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccSystemAlertResourceConfig("test-system-id", "cpu", "gt", 90, 300, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("moneat_system_alert.test", "system_id", "test-system-id"),
					resource.TestCheckResourceAttr("moneat_system_alert.test", "metric", "cpu"),
					resource.TestCheckResourceAttr("moneat_system_alert.test", "condition", "gt"),
					resource.TestCheckResourceAttr("moneat_system_alert.test", "threshold", "90"),
					resource.TestCheckResourceAttr("moneat_system_alert.test", "duration_seconds", "300"),
					resource.TestCheckResourceAttr("moneat_system_alert.test", "enabled", "true"),
					resource.TestCheckResourceAttrSet("moneat_system_alert.test", "id"),
				),
			},
			// ImportState
			{
				ResourceName:      "moneat_system_alert.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "test-system-id/will-be-replaced",
			},
			// Update
			{
				Config: testAccSystemAlertResourceConfig("test-system-id", "memory", "gt", 85, 600, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("moneat_system_alert.test", "metric", "memory"),
					resource.TestCheckResourceAttr("moneat_system_alert.test", "threshold", "85"),
					resource.TestCheckResourceAttr("moneat_system_alert.test", "duration_seconds", "600"),
				),
			},
		},
	})
}

func testAccSystemAlertResourceConfig(systemID, metric, condition string, threshold, duration int, enabled bool) string {
	return fmt.Sprintf(`
resource "moneat_system_alert" "test" {
  system_id        = %[1]q
  metric           = %[2]q
  condition        = %[3]q
  threshold        = %[4]d
  duration_seconds = %[5]d
  enabled          = %[6]t
}
`, systemID, metric, condition, threshold, duration, enabled)
}
