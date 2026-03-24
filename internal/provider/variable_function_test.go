package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccEnvironmentVariableFunction(t *testing.T) {
	const presentVar = "TF_PROVIDER_ENV_FUNCTION_PRESENT"
	const presentValue = "test-value-function"
	t.Setenv(presentVar, presentValue)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVariableFunctionConfig(presentVar),
				Check:  resource.TestCheckOutput("value", presentValue),
			},
		},
	})
}

func TestAccEnvironmentVariableFunction_MissingVariable(t *testing.T) {
	const missingVar = "TF_PROVIDER_ENV_FUNCTION_MISSING"
	testUnsetEnv(t, missingVar)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccVariableFunctionConfig(missingVar),
				ExpectError: canonicalMissingVariableErrorRegexp(),
			},
		},
	})
}

func TestAccEnvironmentVariableFunction_InvalidName(t *testing.T) {
	testCases := []struct {
		name    string
		varName string
	}{
		{
			name:    "empty variable name returns validation error",
			varName: "",
		},
		{
			name:    "whitespace variable name returns validation error",
			varName: " TF_PROVIDER_ENV_FUNCTION_WHITESPACE ",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      testAccVariableFunctionConfig(testCase.varName),
						ExpectError: canonicalInvalidVariableErrorRegexp(),
					},
				},
			})
		})
	}
}

func TestAccEnvironmentVariableFunction_EmptyValue(t *testing.T) {
	t.Setenv("TF_PROVIDER_ENV_FUNCTION_EMPTY", "")

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVariableFunctionConfig("TF_PROVIDER_ENV_FUNCTION_EMPTY"),
				Check:  resource.TestCheckOutput("value", ""),
			},
		},
	})
}
