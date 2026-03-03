package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDashboardResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccDashboardResourceConfig("Test Dashboard", "A test dashboard"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("moneat_dashboard.test", "name", "Test Dashboard"),
					resource.TestCheckResourceAttr("moneat_dashboard.test", "description", "A test dashboard"),
					resource.TestCheckResourceAttrSet("moneat_dashboard.test", "id"),
				),
			},
			// ImportState
			{
				ResourceName:      "moneat_dashboard.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update
			{
				Config: testAccDashboardResourceConfig("Updated Dashboard", "Updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("moneat_dashboard.test", "name", "Updated Dashboard"),
					resource.TestCheckResourceAttr("moneat_dashboard.test", "description", "Updated description"),
				),
			},
		},
	})
}

func testAccDashboardResourceConfig(name, description string) string {
	return fmt.Sprintf(`
resource "moneat_dashboard" "test" {
  name        = %[1]q
  description = %[2]q
}
`, name, description)
}
