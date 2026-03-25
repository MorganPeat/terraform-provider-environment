---
page_title: "environment Provider"
description: "The environment provider reads shell environment variables and makes them available as terraform data sources and provider-defined functions"
---


# environment Provider

The environment provider reads shell environment variables and makes them available as terraform data sources and provider-defined functions.

The data sources work with Terraform 1.0 and later. Provider-defined functions require Terraform 1.8 or later.

Each environment variable is read using the exact name you provide, including any leading or trailing whitespace. This makes it straightforward to use policy-as-code rules in
a language like [HashiCorp Sentinel](https://www.hashicorp.com/sentinel) to control which environment variables are
exposed to terraform state.

Provider-defined function lookups are non-sensitive. For secret use cases, use the `environment_sensitive_variable` data source and Terraform `sensitive(...)` handling.


## Example Usage

```terraform
terraform {
  required_version = ">= 1.8.0"

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