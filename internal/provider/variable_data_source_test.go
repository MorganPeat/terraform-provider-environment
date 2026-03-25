package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnvironmentVariableDataSource(t *testing.T) {
	const presentVar = "TF_PROVIDER_ENV_DATA_SOURCE_PRESENT"
	const presentValue = "test-value-data-source"
	t.Setenv(presentVar, presentValue)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentVariableDataSourceConfig(presentVar),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.environment_variable.test", "id", presentVar),
					resource.TestCheckResourceAttr("data.environment_variable.test", "value", presentValue),
				),
			},
		},
	})
}

func TestAccEnvironmentVariableDataSource_ErrorPaths(t *testing.T) {
	tests := []struct {
		name        string
		varName     string
		expectError *regexp.Regexp
	}{
		{
			name:        "missing environment variable returns error",
			varName:     "TF_PROVIDER_ENV_TEST_DEFINITELY_NOT_SET_XYZ",
			expectError: missingVariableErrorRegexp("TF_PROVIDER_ENV_TEST_DEFINITELY_NOT_SET_XYZ"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testUnsetEnv(t, tt.varName)
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      testAccEnvironmentVariableDataSourceConfig(tt.varName),
						ExpectError: tt.expectError,
					},
				},
			})
		})
	}
}

func TestAccEnvironmentVariableDataSource_EmptyValue(t *testing.T) {
	t.Setenv("TF_PROVIDER_ENV_EMPTY_DATA_SOURCE", "")

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentVariableDataSourceConfig("TF_PROVIDER_ENV_EMPTY_DATA_SOURCE"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.environment_variable.test", "id", "TF_PROVIDER_ENV_EMPTY_DATA_SOURCE"),
					resource.TestCheckResourceAttr("data.environment_variable.test", "value", ""),
				),
			},
		},
	})
}

func TestAccEnvironmentVariableDataSource_WhitespaceName(t *testing.T) {
	const whitespaceVar = " TF_PROVIDER_ENV_DATA_SOURCE_WHITESPACE "
	const whitespaceValue = "test-value-data-source-whitespace"
	const trimmedVar = "TF_PROVIDER_ENV_DATA_SOURCE_WHITESPACE"
	const trimmedValue = "test-value-data-source-trimmed"
	testSetEnv(t, whitespaceVar, whitespaceValue)
	t.Setenv(trimmedVar, trimmedValue)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentVariableDataSourceConfig(whitespaceVar),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.environment_variable.test", "id", whitespaceVar),
					resource.TestCheckResourceAttr("data.environment_variable.test", "value", whitespaceValue),
				),
			},
			{
				Config: testAccEnvironmentVariableDataSourceConfig(trimmedVar),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.environment_variable.test", "id", trimmedVar),
					resource.TestCheckResourceAttr("data.environment_variable.test", "value", trimmedValue),
				),
			},
		},
	})
}

func TestAccEnvironmentVariableDataSource_MissingMessage(t *testing.T) {
	const missingVar = "TF_PROVIDER_ENV_MISSING_NON_SENSITIVE"
	t.Setenv("TF_PROVIDER_ENV_OTHER", "value")
	testUnsetEnv(t, missingVar)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccEnvironmentVariableDataSourceConfig(missingVar),
				ExpectError: missingVariableErrorRegexp(missingVar),
			},
		},
	})
}

func TestAccEnvironmentVariableDataSource_EmptyName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccEnvironmentVariableDataSourceConfig(""),
				ExpectError: missingVariableErrorRegexp(""),
			},
		},
	})
}
