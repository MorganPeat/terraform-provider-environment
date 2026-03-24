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

func testUnsetEnv(t testing.TB, name string) {
	t.Helper()

	originalValue, originalSet := os.LookupEnv(name)

	if err := os.Unsetenv(name); err != nil {
		t.Fatalf("failed to unset environment variable %q: %v", name, err)
	}

	t.Cleanup(func() {
		var err error

		if originalSet {
			err = os.Setenv(name, originalValue)
		} else {
			err = os.Unsetenv(name)
		}

		if err != nil {
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

func canonicalMissingVariableErrorRegexp() *regexp.Regexp {
	parts := strings.Fields(canonicalMissingVariableError)
	quotedParts := make([]string, 0, len(parts))

	for _, part := range parts {
		quotedParts = append(quotedParts, regexp.QuoteMeta(part))
	}

	return regexp.MustCompile(strings.Join(quotedParts, `\s+`))
}

func canonicalInvalidVariableErrorRegexp() *regexp.Regexp {
	parts := strings.Fields(canonicalInvalidVariableError)
	quotedParts := make([]string, 0, len(parts))

	for _, part := range parts {
		quotedParts = append(quotedParts, regexp.QuoteMeta(part))
	}

	return regexp.MustCompile(strings.Join(quotedParts, `\s+`))
}
