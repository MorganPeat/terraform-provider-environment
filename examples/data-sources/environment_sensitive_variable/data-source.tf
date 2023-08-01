
data "environment_sensitive_variable" "path" {
  name = "PATH"
}

output "path" {
  value     = data.environment_sensitive_variable.path.value
  sensitive = true
}