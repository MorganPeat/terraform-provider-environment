# Terraform Provider Environment [![release](https://github.com/MorganPeat/terraform-provider-environment/actions/workflows/release.yml/badge.svg)](https://github.com/MorganPeat/terraform-provider-environment/actions/workflows/release.yml)

The `environment` provider reads shell environment variables and makes them available as terraform data sources and functions.


## Documentation

The documentation for this provider is available on the [Terraform Registry](https://registry.terraform.io/providers/morganpeat/environment/latest/docs).

## Examples

### Data Sources (Terraform 1.0+)

```hcl
terraform {
  required_providers {
    environment = {
      source = "registry.terraform.io/morganpeat/environment"
    }
  }
}

provider "environment" {}


data "environment_variable" "path" {
  name = "PATH"
}

output "path" {
  value = data.environment_variable.path.value
}
```

### Functions (Terraform 1.8+)

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

# Read environment variable using function
output "path" {
  value = provider::environment::variable("PATH")
}

# Read sensitive environment variable
output "api_key" {
  value     = provider::environment::sensitive_variable("API_KEY")
  sensitive = true
}
```

See the [documentation](https://registry.terraform.io/providers/morganpeat/environment/latest/docs) for more details on functions and migration guidance.

## Requirements

* [Terraform](https://www.terraform.io/downloads.html) >= 1.0 (>= 1.8 for functions)
* [Go](https://golang.org/doc/install) >= 1.19

## Building the Provider

To build the provider, you'll need to clone the repository and execute the Go
`install` command from inside the repository's directory.

```bash
go install
```

## Using the provider

The provider can be used by adding it to the [provider
requirements](https://developer.hashicorp.com/terraform/language/providers/requirements).

```terraform
terraform {
  required_providers {
    environment = {
      source = "registry.terraform.io/morganpeat/environment"
    }
  }
}
```

If you wish to use a local provider binary instead, it will need to added to the
[development overrides](https://developer.hashicorp.com/terraform/cli/config/config-file#development-overrides-for-provider-developers).

```terraform
provider_installation {
  dev_overrides {
    "morganpeat/environment" = "/home/developer/go/bin/terraform-provider-environment"
  }

  direct {}
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need
[Go](https://www.golang.org) installed on your machine (see
[Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put
the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

## LICENSE

This project is under [MPL-2.0 license](./LICENSE).
