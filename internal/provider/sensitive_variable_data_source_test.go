package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnvironmentSensitiveVariableDataSource(t *testing.T) {
	const presentVar = "TF_PROVIDER_ENV_SENSITIVE_PRESENT"
	const presentValue = "test-value-sensitive"
	t.Setenv(presentVar, presentValue)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentSensitiveVariableDataSourceConfig(presentVar),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.environment_sensitive_variable.test", "id", presentVar),
					resource.TestCheckResourceAttr("data.environment_sensitive_variable.test", "value", presentValue),
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

func TestAccEnvironmentSensitiveVariableDataSource_EmptyValue(t *testing.T) {
	t.Setenv("TF_PROVIDER_ENV_EMPTY_SENSITIVE", "")

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentSensitiveVariableDataSourceConfig("TF_PROVIDER_ENV_EMPTY_SENSITIVE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.environment_sensitive_variable.test", "id", "TF_PROVIDER_ENV_EMPTY_SENSITIVE"),
					resource.TestCheckResourceAttr("data.environment_sensitive_variable.test", "value", ""),
				),
			},
		},
	})
}

func TestAccEnvironmentSensitiveVariableDataSource_WhitespaceName(t *testing.T) {
	const whitespaceVar = " TF_PROVIDER_ENV_SENSITIVE_WHITESPACE "
	const whitespaceValue = "test-value-sensitive-whitespace"
	const trimmedVar = "TF_PROVIDER_ENV_SENSITIVE_WHITESPACE"
	const trimmedValue = "test-value-sensitive-trimmed"
	testSetEnv(t, whitespaceVar, whitespaceValue)
	t.Setenv(trimmedVar, trimmedValue)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentSensitiveVariableDataSourceConfig(whitespaceVar),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.environment_sensitive_variable.test", "id", whitespaceVar),
					resource.TestCheckResourceAttr("data.environment_sensitive_variable.test", "value", whitespaceValue),
				),
			},
			{
				Config: testAccEnvironmentSensitiveVariableDataSourceConfig(trimmedVar),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.environment_sensitive_variable.test", "id", trimmedVar),
					resource.TestCheckResourceAttr("data.environment_sensitive_variable.test", "value", trimmedValue),
				),
			},
		},
	})
}

func TestAccEnvironmentSensitiveVariableDataSource_EmptyName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccEnvironmentSensitiveVariableDataSourceConfig(""),
				ExpectError: missingVariableErrorRegexp(""),
			},
		},
	})
}

func TestAccEnvironmentSensitiveVariableDataSource_MissingMessage(t *testing.T) {
	const missingVar = "TF_PROVIDER_ENV_MISSING_SENSITIVE"
	testUnsetEnv(t, missingVar)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccEnvironmentSensitiveVariableDataSourceConfig(missingVar),
				ExpectError: missingVariableErrorRegexp(missingVar),
			},
		},
	})
}
