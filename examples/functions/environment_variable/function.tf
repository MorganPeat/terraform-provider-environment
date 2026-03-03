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
  description = "The PATH environment variable"
  value       = provider::environment::variable("PATH")
}

# Use in locals
locals {
  home_dir = provider::environment::variable("HOME")
  user     = provider::environment::variable("USER")
}

output "home_directory" {
  description = "User's home directory"
  value       = local.home_dir
}

output "username" {
  description = "Current username"
  value       = local.user
}

# Use in string interpolation
output "greeting" {
  description = "Personalized greeting"
  value       = "Hello, ${provider::environment::variable("USER")}!"
}

# Multiple variables in a single expression
output "user_info" {
  description = "Combined user information"
  value       = "${provider::environment::variable("USER")} at ${provider::environment::variable("HOME")}"
}
