package provider

import (
	"context"
	"regexp"
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

func TestAccEnvironmentSensitiveVariableDataSource_EmptyAndWhitespaceNames(t *testing.T) {
	testCases := []struct {
		name        string
		varName     string
		expectError *regexp.Regexp
	}{
		{
			name:        "empty variable name returns not found error",
			varName:     "",
			expectError: missingVariableErrorRegexp(""),
		},
		{
			name:        "whitespace variable name returns not found error",
			varName:     " TF_PROVIDER_ENV_SENSITIVE_WHITESPACE ",
			expectError: missingVariableErrorRegexp(" TF_PROVIDER_ENV_SENSITIVE_WHITESPACE "),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      testAccEnvironmentSensitiveVariableDataSourceConfig(testCase.varName),
						ExpectError: testCase.expectError,
					},
				},
			})
		})
	}
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
