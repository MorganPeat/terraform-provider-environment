package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const testAccEnvironmentVariableDataSourceConfig string = `
data "environment_variable" "path" {
	name = "PATH"
  }
`

func TestAccEnvironmentVariableDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentVariableDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.environment_variable.path", "id", "PATH"),
					resource.TestCheckResourceAttr("data.environment_variable.path", "value", os.Getenv("PATH")),
				),
			},
		},
	})
}
