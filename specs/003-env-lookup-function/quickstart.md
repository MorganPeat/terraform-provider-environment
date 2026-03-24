# Quickstart: Environment Lookup Function

## Prerequisites

- Go 1.25+
- Terraform CLI 1.0+
- Local provider build available via development overrides
- Shell environment variables set for test scenarios

## 1) Build provider

```bash
make build
```

## 2) Configure Terraform to use local provider

Use the existing dev override pattern from repository documentation.

## 3) Example: function lookup for existing variable

```hcl
terraform {
  required_providers {
    environment = {
      source = "registry.terraform.io/morganpeat/environment"
    }
  }
}

provider "environment" {}

output "path_via_function" {
  value = provider::environment::variable("PATH")
}
```

Expected: Output equals shell `PATH` value.

## 4) Example: missing variable parity check

Change the function input to a guaranteed-missing key:

```hcl
output "missing_via_function" {
  value = provider::environment::variable("TF_PROVIDER_ENV_NOT_SET")
}
```

Expected: Missing-variable error message must be byte-for-byte identical to the existing data source path.

## 5) Example: empty value behavior

Set an environment variable to empty string and evaluate function + data source:

```bash
export TF_PROVIDER_ENV_EMPTY=""
```

Expected: Both methods return empty value (not not-found).

## 6) Sensitive handling guidance

- There is no separate sensitive function.
- For secret use cases, use `data "environment_sensitive_variable" ...`.
- Where expression flow needs explicit sensitivity, use Terraform `sensitive(...)` handling.

## 7) Run validation tests

```bash
make test
```

Expected: acceptance-style tests covering present/missing/empty parity pass for function and data source methods.
