package provider

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate providers for testing.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"environment": providerserver.NewProtocol6WithError(New("test")()),
}

func testSetEnv(t testing.TB, name, value string) {
	t.Helper()

	originalValue, originalSet := os.LookupEnv(name)

	if err := os.Setenv(name, value); err != nil {
		t.Skipf("cannot set environment variable %q on this platform: %v", name, err)
	}

	t.Cleanup(func() {
		if originalSet {
			if err := os.Setenv(name, originalValue); err != nil {
				t.Errorf("failed to restore environment variable %q: %v", name, err)
			}
			return
		}

		if err := os.Unsetenv(name); err != nil {
			t.Errorf("failed to restore environment variable %q: %v", name, err)
		}
	})
}

func testUnsetEnv(t testing.TB, name string) {
	t.Helper()

	if name == "" {
		// Empty names are invalid on some platforms, but the provider still treats
		// them as a lookup miss during tests.
		return
	}

	originalValue, originalSet := os.LookupEnv(name)

	if err := os.Unsetenv(name); err != nil {
		t.Fatalf("failed to unset environment variable %q: %v", name, err)
	}

	t.Cleanup(func() {
		if originalSet {
			if err := os.Setenv(name, originalValue); err != nil {
				t.Errorf("failed to restore environment variable %q: %v", name, err)
			}
			return
		}

		if err := os.Unsetenv(name); err != nil {
			t.Errorf("failed to restore environment variable %q: %v", name, err)
		}
	})
}

func testAccEnvironmentVariableDataSourceConfig(name string) string {
	return fmt.Sprintf(`
data "environment_variable" "test" {
  name = %q
}
`, name)
}

func testAccEnvironmentSensitiveVariableDataSourceConfig(name string) string {
	return fmt.Sprintf(`
data "environment_sensitive_variable" "test" {
  name = %q
}
`, name)
}

func testAccVariableFunctionConfig(name string) string {
	return fmt.Sprintf(`
output "value" {
  value = provider::environment::variable(%q)
}
`, name)
}

func missingVariableErrorMessage(name string) string {
	return fmt.Sprintf("environment variable '%s' not found", name)
}

func missingVariableErrorRegexp(name string) *regexp.Regexp {
	parts := strings.Fields(missingVariableErrorMessage(name))
	quotedParts := make([]string, 0, len(parts))

	for _, part := range parts {
		quotedParts = append(quotedParts, regexp.QuoteMeta(part))
	}

	return regexp.MustCompile(strings.Join(quotedParts, `\s+`))
}
