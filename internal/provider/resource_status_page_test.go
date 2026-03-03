package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStatusPageResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccStatusPageResourceConfig("Test Status", "test-status", "A test status page", true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("moneat_status_page.test", "name", "Test Status"),
					resource.TestCheckResourceAttr("moneat_status_page.test", "slug", "test-status"),
					resource.TestCheckResourceAttr("moneat_status_page.test", "description", "A test status page"),
					resource.TestCheckResourceAttr("moneat_status_page.test", "is_public", "true"),
					resource.TestCheckResourceAttrSet("moneat_status_page.test", "id"),
				),
			},
			// ImportState
			{
				ResourceName:      "moneat_status_page.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update
			{
				Config: testAccStatusPageResourceConfig("Updated Status", "test-status-v2", "Updated description", false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("moneat_status_page.test", "name", "Updated Status"),
					resource.TestCheckResourceAttr("moneat_status_page.test", "slug", "test-status-v2"),
					resource.TestCheckResourceAttr("moneat_status_page.test", "is_public", "false"),
				),
			},
		},
	})
}

func testAccStatusPageResourceConfig(name, slug, description string, isPublic bool) string {
	return fmt.Sprintf(`
resource "moneat_status_page" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[3]q
  is_public   = %[4]t
}
`, name, slug, description, isPublic)
}
