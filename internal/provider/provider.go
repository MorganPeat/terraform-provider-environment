package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the provider.Provider interface.
var _ provider.Provider = &environmentProvider{}

// environmentProvider implements the provider.Provider interface and is the
// root structure used to interface with the terraform plugin framework.
// See https://developer.hashicorp.com/terraform/plugin/framework/providers
type environmentProvider struct {
	Version string
}

// Metadata returns the provider type name.
func (p *environmentProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "environment"
}

// Schema defines the provider-level schema for configuration data.
func (p *environmentProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// Provider specific implementation.
		},
	}
}

// Configure prepares the provider for data sources and resources.
func (p *environmentProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Provider specific implementation.
}

// DataSources defines the data sources implemented in the provider.
func (p *environmentProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Provider specific implementation
	}
}

// Resources defines the resources implemented in the provider.
func (p *environmentProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Provider specific implementation
	}
}

// New creates a new environmentProvider
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &environmentProvider{
			Version: version,
		}
	}
}
