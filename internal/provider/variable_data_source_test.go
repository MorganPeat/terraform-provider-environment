package provider

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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

func TestAccEnvironmentVariableDataSource_ErrorPaths(t *testing.T) {
	tests := []struct {
		name        string
		varName     string
		expectError string
	}{
		{
			name:        "missing environment variable returns error",
			varName:     "TF_PROVIDER_ENV_TEST_DEFINITELY_NOT_SET_XYZ",
			expectError: "Not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Unsetenv(tt.varName)
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`
data "environment_variable" "test" {
  name = %q
}
`, tt.varName),
						ExpectError: regexp.MustCompile(tt.expectError),
					},
				},
			})
		})
	}
}

// Made with Bob
