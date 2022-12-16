package main

import (
	"context"
	"flag"
	"log"

	"github.com/MorganPeat/terraform-provider-environment/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Ensure examples are correctly formatted
//go:generate terraform fmt -recursive ./examples/

// Auto-generates documentation for the terraform registry
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

var (
	// These will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"
)

func main() {

	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/morganpeat/environment",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
