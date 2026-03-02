# Provider-Defined Functions Design Specification

## Overview

This document outlines the design for implementing provider-defined functions in the terraform-provider-environment, as requested in [Issue #61](https://github.com/MorganPeat/terraform-provider-environment/issues/61).

## Background

Terraform 1.8+ introduced support for provider-defined functions, which offer a simpler, more intuitive syntax compared to data sources. This feature allows users to call functions directly in expressions without requiring resource blocks.

### Current Implementation (Data Sources)
```hcl
data "environment_variable" "path" {
  name = "PATH"
}

output "path" {
  value = data.environment_variable.path.value
}
```

### Proposed Implementation (Functions)
```hcl
output "path" {
  value = provider::environment::variable("PATH")
}
```

## Requirements

### Functional Requirements
1. **Backward Compatibility**: Maintain existing data sources for Terraform < 1.8
2. **Feature Parity**: Provide equivalent functionality to existing data sources
3. **Sensitive Value Support**: Handle sensitive environment variables appropriately
4. **Error Handling**: Provide clear error messages when variables don't exist
5. **Documentation**: Update docs with function examples and migration guidance

### Non-Functional Requirements
1. **Performance**: Functions should be lightweight and fast
2. **Testing**: Comprehensive unit and acceptance tests
3. **Maintainability**: Code should follow existing patterns and be easy to maintain

## Technical Analysis

### Framework Support

The project uses `terraform-plugin-framework v1.18.0`, which has full support for provider-defined functions through the `function` package.

Key framework capabilities:
- **String Parameters**: Functions can accept string parameters
- **String Returns**: Functions can return string values
- **Sensitive Returns**: Functions can mark return values as sensitive
- **Error Handling**: Functions can return diagnostics for error conditions
- **Validation**: Parameter validation is built into the framework

### Function vs Data Source Comparison

| Aspect | Data Source | Function |
|--------|-------------|----------|
| Syntax | `data.environment_variable.name.value` | `provider::environment::variable("NAME")` |
| State | Stored in state file | Not stored in state |
| Lifecycle | Managed by Terraform lifecycle | Evaluated on-demand |
| Complexity | Requires resource block | Direct inline usage |
| Performance | Cached in state | Evaluated each time |
| Sensitive Values | Schema attribute `Sensitive: true` | Return type `function.StringReturn{Sensitive: true}` |

### Design Decisions

#### 1. Function Naming Convention

**Decision**: Implement two functions:
- `variable(name string) string` - Returns environment variable value
- `sensitive_variable(name string) string` - Returns sensitive environment variable value

**Rationale**:
- Mirrors existing data source naming (`environment_variable` and `environment_sensitive_variable`)
- Clear distinction between sensitive and non-sensitive values
- Terraform functions cannot dynamically determine sensitivity, so separate functions are needed
- Maintains consistency with current provider API

**Alternative Considered**: Single function with optional `sensitive` parameter
- **Rejected**: Terraform function return types must be statically defined; cannot conditionally mark as sensitive

#### 2. Error Handling

**Decision**: Return error diagnostic when environment variable doesn't exist

**Behavior**:
```hcl
# If MY_VAR is not set:
output "test" {
  value = provider::environment::variable("MY_VAR")
}
# Error: Environment variable "MY_VAR" not found
```

**Error Message Format**: `Environment variable "NAME" not found`
- Consistent with data source error handling
- Uses quoted variable name for clarity
- Clear and concise

**Rationale**:
- Matches existing data source behavior (fails on missing variable)
- Prevents silent failures
- Clear error messages for debugging

**Edge Cases**:
- Empty string value (`VAR=""`) returns empty string (not an error)
- Unset variable returns error
- This matches `os.LookupEnv()` behavior

**Alternative Considered**: Return empty string for missing variables
- **Rejected**: Could mask configuration errors; explicit failure is better

#### 3. Implementation Architecture

**Decision**: Implement functions in separate files following existing patterns

**Structure**:
```
internal/provider/
├── provider.go                          # Add Functions() method
├── variable_data_source.go              # Existing
├── sensitive_variable_data_source.go    # Existing
├── variable_function.go                 # New
├── sensitive_variable_function.go       # New
├── variable_function_test.go            # New
└── sensitive_variable_function_test.go  # New
```

**Rationale**:
- Follows existing code organization patterns
- Separates concerns (data sources vs functions)
- Makes testing straightforward
- Easy to maintain and extend

## Implementation Specification

### 1. Provider Interface Update

**File**: `internal/provider/provider.go`

Add interface assertion and `Functions()` method to `environmentProvider`:

```go
// Ensure the implementation satisfies the provider.ProviderWithFunctions interface.
var _ provider.ProviderWithFunctions = &environmentProvider{}

// Functions defines the functions implemented in the provider.
func (p *environmentProvider) Functions(ctx context.Context) []func() function.Function {
    return []func() function.Function{
        NewVariableFunction,
        NewSensitiveVariableFunction,
    }
}
```

### 2. Variable Function Implementation

**File**: `internal/provider/variable_function.go`

```go
package provider

import (
    "context"
    "fmt"
    "os"

    "github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = &variableFunction{}

type variableFunction struct{}

func NewVariableFunction() function.Function {
    return &variableFunction{}
}

func (f *variableFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
    resp.Name = "variable"
}

func (f *variableFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
    resp.Definition = function.Definition{
        Summary:     "Read an environment variable",
        Description: "Returns the value of the specified environment variable.",
        Parameters: []function.Parameter{
            function.StringParameter{
                Name:        "name",
                Description: "The name of the environment variable to read.",
            },
        },
        Return: function.StringReturn{},
    }
}

// Run executes the function logic, reading the environment variable
// and returning its value or an error if not found.
func (f *variableFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
    var name string
    resp.Error = req.Arguments.Get(ctx, &name)
    if resp.Error != nil {
        return
    }

    value, ok := os.LookupEnv(name)
    if !ok {
        resp.Error = function.NewArgumentFuncError(
            0,
            fmt.Sprintf("Environment variable %q not found", name),
        )
        return
    }

    resp.Error = resp.Result.Set(ctx, value)
}
```

### 3. Sensitive Variable Function Implementation

**File**: `internal/provider/sensitive_variable_function.go`

```go
package provider

import (
    "context"
    "fmt"
    "os"

    "github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = &sensitiveVariableFunction{}

type sensitiveVariableFunction struct{}

func NewSensitiveVariableFunction() function.Function {
    return &sensitiveVariableFunction{}
}

func (f *sensitiveVariableFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
    resp.Name = "sensitive_variable"
}

func (f *sensitiveVariableFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
    resp.Definition = function.Definition{
        Summary:     "Read a sensitive environment variable",
        Description: "Returns the value of the specified environment variable, marked as sensitive.",
        Parameters: []function.Parameter{
            function.StringParameter{
                Name:        "name",
                Description: "The name of the environment variable to read.",
            },
        },
        Return: function.StringReturn{
            Sensitive: true,
        },
    }
}

// Run executes the function logic, reading the environment variable
// and returning its value marked as sensitive, or an error if not found.
func (f *sensitiveVariableFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
    var name string
    resp.Error = req.Arguments.Get(ctx, &name)
    if resp.Error != nil {
        return
    }

    value, ok := os.LookupEnv(name)
    if !ok {
        resp.Error = function.NewArgumentFuncError(
            0,
            fmt.Sprintf("Environment variable %q not found", name),
        )
        return
    }

    resp.Error = resp.Result.Set(ctx, value)
}
```

### 4. Testing Strategy

#### Unit Tests

**File**: `internal/provider/variable_function_test.go`

Test cases:
1. **Success case**: Variable exists and is returned correctly
2. **Error case**: Variable doesn't exist, returns appropriate error
3. **Empty value case**: Variable exists but is empty string (should return empty string, not error)
4. **Special characters**: Variable name and value with special characters (including unicode, spaces, quotes)
5. **Multiline values**: Variable with newlines in the value
6. **Large values**: Variable with very long value (>1KB)

```go
func TestVariableFunction_Success(t *testing.T) {
    // Set environment variable
    os.Setenv("TEST_VAR", "test_value")
    defer os.Unsetenv("TEST_VAR")

    // Test function execution
    // Assert value matches
}

func TestVariableFunction_NotFound(t *testing.T) {
    // Ensure variable doesn't exist
    os.Unsetenv("NONEXISTENT_VAR")

    // Test function execution
    // Assert error is returned
}
```

**File**: `internal/provider/sensitive_variable_function_test.go`

Test cases:
1. **Sensitivity verification**: Ensure return value is marked as sensitive
2. **Same behavior as regular function**: Variable lookup works identically
3. **Error handling**: Same error behavior as regular function

#### Acceptance Tests

Add acceptance tests to verify functions work in actual Terraform configurations:

```go
func TestAccVariableFunction(t *testing.T) {
    os.Setenv("TF_TEST_VAR", "test_value")
    defer os.Unsetenv("TF_TEST_VAR")

    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: `
                    output "test" {
                        value = provider::environment::variable("TF_TEST_VAR")
                    }
                `,
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckOutput("test", "test_value"),
                ),
            },
        },
    })
}
```

### 5. Documentation Updates

#### Provider Documentation

**File**: `templates/index.md.tmpl`

Add section on functions after the data sources section:

```markdown
## Functions (Terraform 1.8+)

This provider supports provider-defined functions for a more concise syntax when using Terraform 1.8 or later:

### variable

Returns the value of an environment variable.

**Signature**: `provider::environment::variable(name string) string`

**Example**:
```hcl
output "path" {
  value = provider::environment::variable("PATH")
}
```

### sensitive_variable

Returns the value of an environment variable, marked as sensitive.

**Signature**: `provider::environment::sensitive_variable(name string) string`

**Example**:
```hcl
output "api_key" {
  value     = provider::environment::sensitive_variable("API_KEY")
  sensitive = true
}
```

## Migration from Data Sources to Functions

If you're using Terraform 1.8 or later, you can migrate from data sources to functions:

**Before (Data Source)**:
```hcl
data "environment_variable" "path" {
  name = "PATH"
}

output "path" {
  value = data.environment_variable.path.value
}
```

**After (Function)**:
```hcl
output "path" {
  value = provider::environment::variable("PATH")
}
```

**Benefits of Functions**:
- More concise syntax
- No state storage overhead
- Direct inline usage in expressions
- Evaluated on-demand

**When to Use Functions vs Data Sources**:
- Use **functions** (Terraform 1.8+) for simpler, inline access to environment variables
- Use **data sources** (Terraform 1.0+) when you need backward compatibility or prefer explicit resource blocks
- Both approaches can be used together in the same configuration

**Note**: Data sources remain fully supported for backward compatibility with Terraform < 1.8.
```

#### Function-Specific Documentation

Create new documentation files:

**File**: `docs/functions/variable.md`
**File**: `docs/functions/sensitive_variable.md`

These will be auto-generated by `tfplugindocs` from the function definitions.

#### Example Files

**File**: `examples/functions/environment_variable/function.tf`

```hcl
# Example using the variable function
# Requires Terraform 1.8 or later

terraform {
  required_version = ">= 1.8"
  required_providers {
    environment = {
      source = "registry.terraform.io/morganpeat/environment"
    }
  }
}

provider "environment" {}

# Read PATH environment variable using function
output "path" {
  value = provider::environment::variable("PATH")
}

# Use in locals
locals {
  home_dir = provider::environment::variable("HOME")
  user     = provider::environment::variable("USER")
}

# Use in resource configuration
resource "null_resource" "example" {
  triggers = {
    path = provider::environment::variable("PATH")
  }
}
```

**File**: `examples/functions/environment_sensitive_variable/function.tf`

```hcl
# Example using the sensitive_variable function
# Requires Terraform 1.8 or later

terraform {
  required_version = ">= 1.8"
  required_providers {
    environment = {
      source = "registry.terraform.io/morganpeat/environment"
    }
  }
}

provider "environment" {}

# Read sensitive environment variable
output "api_key" {
  value     = provider::environment::sensitive_variable("API_KEY")
  sensitive = true
}

# Use in resource configuration
resource "null_resource" "example" {
  triggers = {
    # Value is marked as sensitive
    secret = provider::environment::sensitive_variable("SECRET_TOKEN")
  }
}
```

### 6. README Updates

**File**: `README.md`

Add section after the data source example:

```markdown
## Functions (Terraform 1.8+)

For Terraform 1.8 and later, you can use provider-defined functions for a more concise syntax:

```hcl
terraform {
  required_version = ">= 1.8"
  required_providers {
    environment = {
      source = "registry.terraform.io/morganpeat/environment"
    }
  }
}

provider "environment" {}

output "path" {
  value = provider::environment::variable("PATH")
}

output "api_key" {
  value     = provider::environment::sensitive_variable("API_KEY")
  sensitive = true
}
```

See the [documentation](https://registry.terraform.io/providers/morganpeat/environment/latest/docs) for more details.
```

## Implementation Phases

### Phase 1: Core Implementation
1. Add `Functions()` method to provider
2. Implement `variable` function
3. Implement `sensitive_variable` function
4. Add unit tests for both functions

### Phase 2: Testing
1. Add acceptance tests
2. Test with actual Terraform 1.8+ configurations
3. Verify backward compatibility (data sources still work)

### Phase 3: Documentation
1. Update provider documentation template
2. Create function example files
3. Update README with function examples
4. Run `make generate` to regenerate docs

### Phase 4: Release
1. Update CHANGELOG
2. Create release notes highlighting new functions
3. Tag release with appropriate version bump (minor version)

## Testing Checklist

- [ ] Unit tests pass for `variable` function
  - [ ] Variable exists and returns correct value
  - [ ] Variable doesn't exist returns error
  - [ ] Empty string value handled correctly
  - [ ] Special characters (unicode, quotes, spaces)
  - [ ] Multiline values
  - [ ] Large values (>1KB)
- [ ] Unit tests pass for `sensitive_variable` function
  - [ ] Same test cases as variable function
  - [ ] Verify sensitivity marking
- [ ] Acceptance tests pass for both functions
- [ ] Functions work with Terraform 1.8+
- [ ] Data sources still work (backward compatibility)
- [ ] Mixed usage (functions + data sources in same config)
- [ ] Error messages are clear and helpful
- [ ] Error messages consistent with data sources
- [ ] Sensitive values are properly marked
- [ ] Documentation generates correctly
- [ ] Examples run successfully
- [ ] Linter passes (`make lint`)
- [ ] All tests pass (`make test` and `make testacc`)

## Backward Compatibility

This implementation maintains full backward compatibility:

1. **Existing data sources remain unchanged**: No modifications to existing data source code
2. **No breaking changes**: All existing Terraform configurations continue to work
3. **Terraform version support**: 
   - Functions require Terraform 1.8+
   - Data sources work with Terraform 1.0+
   - Provider gracefully handles both

## Performance Considerations

**Functions vs Data Sources**:
- **Functions**: Evaluated on-demand, not stored in state
- **Data Sources**: Cached in state file

**Implications**:
- Functions may be evaluated multiple times in a single plan/apply
- For frequently accessed values, data sources may be more efficient
- For one-time reads, functions are simpler and cleaner

**Recommendation**: Document both approaches and let users choose based on their use case.

## Security Considerations

1. **Sensitive Values**: 
   - `sensitive_variable` function marks return as sensitive
   - Terraform will mask these values in output
   - Users must still mark outputs as `sensitive = true`

2. **Environment Variable Access**:
   - Functions have same access as data sources
   - No additional security concerns
   - Policy-as-code rules (Sentinel) can still control access

## Future Enhancements

Potential future improvements (not in scope for initial implementation):

1. **Bulk read function**: `variables(names []string) map[string]string`
2. **Default value support**: `variable_with_default(name string, default string) string`
3. **Validation function**: `variable_exists(name string) bool`
4. **JSON parsing function**: `json_variable(name string) dynamic`

These can be considered in future releases based on user feedback.

## References

- [Terraform 1.8 Release Notes](https://github.com/hashicorp/terraform/blob/v1.8.0/CHANGELOG.md)
- [Provider-Defined Functions Documentation](https://developer.hashicorp.com/terraform/plugin/framework/functions)
- [terraform-plugin-framework Functions Package](https://pkg.go.dev/github.com/hashicorp/terraform-plugin-framework/function)
- [Issue #61](https://github.com/MorganPeat/terraform-provider-environment/issues/61)

## Conclusion

This design provides a comprehensive plan for implementing provider-defined functions while maintaining backward compatibility and following Terraform best practices. The implementation is straightforward, well-tested, and provides clear value to users on Terraform 1.8+.