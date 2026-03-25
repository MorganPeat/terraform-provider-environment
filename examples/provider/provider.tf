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
