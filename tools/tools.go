//go:build tools

package tools

import (
	// Makes sure that `go mod tidy` does not remove the dependency
	_ "github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs"
)
