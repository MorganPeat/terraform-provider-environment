package provider

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
)

func TestAccVariableFunction_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::environment::variable("PATH")
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(os.Getenv("PATH"))),
				},
			},
		},
	})
}

func TestAccVariableFunction_NotFound(t *testing.T) {
	os.Unsetenv("TF_TEST_NONEXISTENT_VAR_XYZ")
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::environment::variable("TF_TEST_NONEXISTENT_VAR_XYZ")
					}
				`,
				ExpectError: regexp.MustCompile(`Environment variable\s+"TF_TEST_NONEXISTENT_VAR_XYZ"\s+not found`),
			},
		},
	})
}

func TestAccVariableFunction_EmptyValue(t *testing.T) {
	os.Setenv("TF_TEST_EMPTY_VAR", "")
	defer os.Unsetenv("TF_TEST_EMPTY_VAR")

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::environment::variable("TF_TEST_EMPTY_VAR")
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact("")),
				},
			},
		},
	})
}

func TestAccVariableFunction_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "unicode characters",
			value: "Hello 世界 🌍",
		},
		{
			name:  "spaces and quotes",
			value: `value with "quotes" and 'apostrophes'`,
		},
		{
			name:  "special symbols",
			value: "!@#$%^&*()_+-=[]{}|;:,.<>?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			varName := "TF_TEST_SPECIAL_CHARS"
			os.Setenv(varName, tt.value)
			defer os.Unsetenv(varName)

			resource.Test(t, resource.TestCase{
				IsUnitTest:               true,
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`
							output "test" {
								value = provider::environment::variable(%q)
							}
						`, varName),
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(tt.value)),
						},
					},
				},
			})
		})
	}
}

func TestAccVariableFunction_MultilineValue(t *testing.T) {
	multilineValue := "line1\nline2\nline3"
	os.Setenv("TF_TEST_MULTILINE", multilineValue)
	defer os.Unsetenv("TF_TEST_MULTILINE")

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::environment::variable("TF_TEST_MULTILINE")
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(multilineValue)),
				},
			},
		},
	})
}

func TestAccVariableFunction_LargeValue(t *testing.T) {
	// Create a value larger than 1KB
	largeValue := strings.Repeat("a", 2000)
	os.Setenv("TF_TEST_LARGE_VALUE", largeValue)
	defer os.Unsetenv("TF_TEST_LARGE_VALUE")

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = provider::environment::variable("TF_TEST_LARGE_VALUE")
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(largeValue)),
				},
			},
		},
	})
}

func TestAccVariableFunction_InLocals(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					locals {
						path = provider::environment::variable("PATH")
					}
					
					output "test" {
						value = local.path
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact(os.Getenv("PATH"))),
				},
			},
		},
	})
}

func TestAccVariableFunction_MultipleVariables(t *testing.T) {
	os.Setenv("TF_TEST_VAR1", "value1")
	os.Setenv("TF_TEST_VAR2", "value2")
	defer os.Unsetenv("TF_TEST_VAR1")
	defer os.Unsetenv("TF_TEST_VAR2")

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test1" {
						value = provider::environment::variable("TF_TEST_VAR1")
					}
					
					output "test2" {
						value = provider::environment::variable("TF_TEST_VAR2")
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test1", knownvalue.StringExact("value1")),
					statecheck.ExpectKnownOutputValue("test2", knownvalue.StringExact("value2")),
				},
			},
		},
	})
}

func TestAccVariableFunction_InStringInterpolation(t *testing.T) {
	os.Setenv("TF_TEST_USER", "testuser")
	defer os.Unsetenv("TF_TEST_USER")

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					output "test" {
						value = "Hello, ${provider::environment::variable("TF_TEST_USER")}!"
					}
				`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue("test", knownvalue.StringExact("Hello, testuser!")),
				},
			},
		},
	})
}

// Made with Bob
