
terraform {
  required_version = ">= 1.8.0"
}

data "environment_sensitive_variable" "path" {
  name = "PATH"
}

output "path" {
  value     = data.environment_sensitive_variable.path.value
  sensitive = true
}

output "path_via_function_with_explicit_sensitivity" {
  value     = sensitive(provider::environment::variable("PATH"))
  sensitive = true
}
