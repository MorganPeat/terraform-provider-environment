package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLookupParity_PresentAndEmptyValues(t *testing.T) {
	t.Setenv("TF_PROVIDER_ENV_PARITY_PRESENT", "present-value")
	t.Setenv("TF_PROVIDER_ENV_PARITY_EMPTY", "")

	testCases := []struct {
		name          string
		variableName  string
		expectedValue string
	}{
		{
			name:          "present value",
			variableName:  "TF_PROVIDER_ENV_PARITY_PRESENT",
			expectedValue: "present-value",
		},
		{
			name:          "empty value",
			variableName:  "TF_PROVIDER_ENV_PARITY_EMPTY",
			expectedValue: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: testAccLookupParitySuccessConfig(testCase.variableName),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("data.environment_variable.test", "id", testCase.variableName),
							resource.TestCheckResourceAttr("data.environment_variable.test", "value", testCase.expectedValue),
							resource.TestCheckResourceAttr("data.environment_sensitive_variable.test", "id", testCase.variableName),
							resource.TestCheckResourceAttr("data.environment_sensitive_variable.test", "value", testCase.expectedValue),
							resource.TestCheckOutput("function_value", testCase.expectedValue),
							resource.TestCheckOutput("non_sensitive_value", testCase.expectedValue),
							resource.TestCheckOutput("sensitive_value", testCase.expectedValue),
						),
					},
				},
			})
		})
	}
}

func TestAccLookupParity_MissingEmptyAndWhitespaceNames(t *testing.T) {
	testCases := []struct {
		name         string
		variableName string
		expectError  *regexp.Regexp
	}{
		{
			name:         "missing variable name",
			variableName: "TF_PROVIDER_ENV_PARITY_MISSING",
			expectError:  missingVariableErrorRegexp("TF_PROVIDER_ENV_PARITY_MISSING"),
		},
		{
			name:         "empty variable name",
			variableName: "",
			expectError:  missingVariableErrorRegexp(""),
		},
		{
			name:         "whitespace variable name",
			variableName: " TF_PROVIDER_ENV_PARITY_WHITESPACE ",
			expectError:  missingVariableErrorRegexp(" TF_PROVIDER_ENV_PARITY_WHITESPACE "),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testUnsetEnv(t, testCase.variableName)

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      testAccEnvironmentVariableDataSourceConfig(testCase.variableName),
						ExpectError: testCase.expectError,
					},
					{
						Config:      testAccEnvironmentSensitiveVariableDataSourceConfig(testCase.variableName),
						ExpectError: testCase.expectError,
					},
					{
						Config:      testAccVariableFunctionConfig(testCase.variableName),
						ExpectError: testCase.expectError,
					},
				},
			})
		})
	}
}

func testAccLookupParitySuccessConfig(name string) string {
	return fmt.Sprintf(`
data "environment_variable" "test" {
  name = %q
}

data "environment_sensitive_variable" "test" {
  name = %q
}

output "non_sensitive_value" {
  value = data.environment_variable.test.value
}

output "sensitive_value" {
  value = data.environment_sensitive_variable.test.value
  sensitive = true
}

output "function_value" {
  value = provider::environment::variable(%q)
}
`, name, name, name)
}
