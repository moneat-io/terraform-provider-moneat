package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: testAccProjectResourceConfig("test-project", "python"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("moneat_project.test", "name", "test-project"),
					resource.TestCheckResourceAttr("moneat_project.test", "platform", "python"),
					resource.TestCheckResourceAttrSet("moneat_project.test", "id"),
				),
			},
			// ImportState
			{
				ResourceName:      "moneat_project.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update
			{
				Config: testAccProjectResourceConfig("test-project-updated", "javascript"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("moneat_project.test", "name", "test-project-updated"),
					resource.TestCheckResourceAttr("moneat_project.test", "platform", "javascript"),
				),
			},
		},
	})
}

func testAccProjectResourceConfig(name, platform string) string {
	return fmt.Sprintf(`
resource "moneat_project" "test" {
  name     = %[1]q
  platform = %[2]q
}
`, name, platform)
}

func TestAccProjectDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "moneat_project" "test" {
  name     = "ds-test-project"
  platform = "go"
}

data "moneat_project" "test" {
  id = moneat_project.test.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.moneat_project.test", "name", "ds-test-project"),
					resource.TestCheckResourceAttr("data.moneat_project.test", "platform", "go"),
				),
			},
		},
	})
}

func TestAccProjectsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
resource "moneat_project" "test" {
  name     = "list-test-project"
  platform = "python"
}

data "moneat_projects" "all" {
  depends_on = [moneat_project.test]
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.moneat_projects.all", "projects.#"),
				),
			},
		},
	})
}
