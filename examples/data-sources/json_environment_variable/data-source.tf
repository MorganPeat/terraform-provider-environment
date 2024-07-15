#
# Shows how this provider can read a properly escaped json document, stored in
# an environment variable, into a HCL object.
#
# Set the environment variable before running this example:
#   export JSON_VALUE="{\"an-integer\":1,\"an-object\":{\"a-string\":\"hello, world\"},\"an-array\":[2,4,6,8]}"
#

data "environment_variable" "json" {
  name = "JSON_VALUE"
}

output "json_value" {
  value = jsondecode(data.environment_variable.json.value)
}


locals {
  json_value = jsondecode(data.environment_variable.json.value)
}

output "string_value" {
  value = local.json_value["an-object"]["a-string"]
}