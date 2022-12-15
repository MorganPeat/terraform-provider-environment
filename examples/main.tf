terraform {
  required_providers {
    environment = {
      source = "registry.terraform.io/morganpeat/environment"
    }
  }
}

provider "environment" {

}


data "environment_variable" "path" {
  id = "PATH"
}

output "path" {
  value = data.environment_variable.path.value
}