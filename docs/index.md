---
page_title: "environment Provider"
description: "The environment provider reads shell environment variables and makes them available as a terraform data source"
---


# environment Provider

The environment provider reads shell environment variables and makes them available as a terraform data source.  

Each environment variable is read using its explicit name. This makes it straightforward to use policy-as-code rules in
a language like [HashiCorp Sentinel](https://www.hashicorp.com/sentinel) to control which environment variables are
exposed to terraform state.


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