package provider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccEnvironmentSensitiveVariableDataSourceConfig string = `
data "environment_sensitive_variable" "path" {
	name = "PATH"
  }
`

func TestAccEnvironmentSensitiveVariableDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentSensitiveVariableDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.environment_sensitive_variable.path", "id", "PATH"),
					resource.TestCheckResourceAttr("data.environment_sensitive_variable.path", "value", os.Getenv("PATH")),
				),
			},
		},
	})
}

func TestAccEnvironmentSensitiveVariableDataSourceCheckSensitiveAttributes(t *testing.T) {
	dataSource := NewSensitiveVariableDataSource()
	schemaResponse := datasource.SchemaResponse{}

	dataSource.Schema(context.Background(), datasource.SchemaRequest{}, &schemaResponse)
	if !schemaResponse.Schema.Attributes["value"].IsSensitive() {
		t.Errorf("attribute 'value' should be marked as 'Sensitive'")
	}
}
