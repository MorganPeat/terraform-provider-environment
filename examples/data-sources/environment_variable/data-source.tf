
data "environment_variable" "path" {
  name = "PATH"
}

output "path" {
  value = data.environment_variable.path.value
}