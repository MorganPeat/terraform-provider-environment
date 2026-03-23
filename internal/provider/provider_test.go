package provider

import (
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate providers for testing.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"environment": providerserver.NewProtocol6WithError(New("test")()),
}

// isTFAccNotSet returns true when TF_ACC is not present in the environment,
// enabling in-process unit test mode.
func isTFAccNotSet() bool {
	_, ok := os.LookupEnv("TF_ACC")
	return !ok
}
