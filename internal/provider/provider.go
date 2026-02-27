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

// Metadata returns the metadata for the provider, such as
// the type name and version data.
func (p *environmentProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	// TypeName is the prefix used in each data source and resource name.
	resp.TypeName = "environment"

	// Version is not used by the terraform framework yet.
	resp.Version = p.Version
}

// Schema returns the schema for this provider.
func (p *environmentProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
The environment provider reads shell environment variables and makes them available as a terraform data source.  

Each environment variable is read using its explicit name. This makes it straightforward to use policy-as-code rules in
a language like [HashiCorp Sentinel](https://www.hashicorp.com/sentinel) to control which environment variables are
exposed to terraform state.
`,
	}
}

// Configure prepares the provider for use, with the values the user specified
// in the provider configuration block.
func (p *environmentProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

// DataSources defines the data sources implemented in the provider.
func (p *environmentProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewVariableDataSource,
		NewSensitiveVariableDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *environmentProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

// New creates a new environmentProvider.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &environmentProvider{
			Version: version,
		}
	}
}
