---
page_title: "environment Provider"
description: "The environment provider reads shell environment variables and makes them available as terraform data sources and functions"
---


# environment Provider

The environment provider reads shell environment variables and makes them available as terraform data sources and functions.

Each environment variable is read using its explicit name. This makes it straightforward to use policy-as-code rules in
a language like [HashiCorp Sentinel](https://www.hashicorp.com/sentinel) to control which environment variables are
exposed to terraform state.

For Terraform 1.8+, provider-defined functions offer a more concise syntax for accessing environment variables.


## Example Usage

```terraform
terraform {
  required_providers {
    environment = {
      source = "registry.terraform.io/morganpeat/environment"
    }
  }
}

provider "environment" {}
```

## Functions (Terraform 1.8+)

This provider supports provider-defined functions for a more concise syntax when using Terraform 1.8 or later.

### variable

Returns the value of an environment variable.

**Signature**: `provider::environment::variable(name string) string`

**Example**:
```hcl
output "path" {
  value = provider::environment::variable("PATH")
}
```

See the [variable function documentation](functions/variable.md) for more details.

### sensitive_variable

Returns the value of an environment variable. The caller should mark outputs using this function as sensitive.

**Signature**: `provider::environment::sensitive_variable(name string) string`

**Example**:
```hcl
output "api_key" {
  value     = provider::environment::sensitive_variable("API_KEY")
  sensitive = true
}
```

See the [sensitive_variable function documentation](functions/sensitive_variable.md) for more details.

## Migration from Data Sources to Functions

If you're using Terraform 1.8 or later, you can migrate from data sources to functions for a more concise syntax:

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

### Benefits of Functions

- **More concise syntax**: Direct inline usage without resource blocks
- **No state storage overhead**: Functions are evaluated on-demand
- **Direct inline usage**: Can be used anywhere in expressions
- **Evaluated on-demand**: No caching in state file

### When to Use Functions vs Data Sources

- Use **functions** (Terraform 1.8+) for simpler, inline access to environment variables
- Use **data sources** (Terraform 1.0+) when you need backward compatibility or prefer explicit resource blocks
- Both approaches can be used together in the same configuration

**Note**: Data sources remain fully supported for backward compatibility with Terraform < 1.8.