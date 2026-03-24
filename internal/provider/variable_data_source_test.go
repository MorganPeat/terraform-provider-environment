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
			testUnsetEnv(t, tt.varName)
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      testAccEnvironmentVariableDataSourceConfig(tt.varName),
						ExpectError: regexp.MustCompile(tt.expectError),
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
				ExpectError: regexp.MustCompile(regexp.QuoteMeta(canonicalMissingVariableError)),
			},
		},
	})
}

func TestAccEnvironmentVariableDataSource_InvalidName(t *testing.T) {
	testCases := []struct {
		name        string
		varName     string
		expectError *regexp.Regexp
	}{
		{
			name:        "empty variable name returns validation error",
			varName:     "",
			expectError: canonicalInvalidVariableErrorRegexp(),
		},
		{
			name:        "whitespace variable name returns validation error",
			varName:     " TF_PROVIDER_ENV_WHITESPACE ",
			expectError: canonicalInvalidVariableErrorRegexp(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      testAccEnvironmentVariableDataSourceConfig(testCase.varName),
						ExpectError: testCase.expectError,
					},
				},
			})
		})
	}
}
