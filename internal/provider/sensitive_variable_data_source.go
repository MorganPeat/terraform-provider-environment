package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &sensitiveVariableDataSource{}

// NewSensitiveVariableDataSource is a helper function to simplify the provider implementation.
func NewSensitiveVariableDataSource() datasource.DataSource {
	return &sensitiveVariableDataSource{}
}

// sensitiveVariableDataSource is the data source implementation.
type sensitiveVariableDataSource struct{}

// Metadata returns the data source type name.
func (d *sensitiveVariableDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sensitive_variable"
}

// Schema defines the schema for the data source.
func (d *sensitiveVariableDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
The sensitive variable data source exposes a sensitive shell environment variable to terraform.

Any change in the value of the shell environment variable will show up as a change in the terraform plan.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for this resource. This matches the name of the environment variable.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the shell environment variable to read.",
			},
			"value": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "The value of the shell environment variable.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *sensitiveVariableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// NOTE: We can use the read-method for the data source `environment_variable` as-is, because
	// all this data source does, is adding "Sensitive: true" to the schema of the property.
	//
	// The values and the property names are meant to be kept the same between data sources.
	NewVariableDataSource().Read(ctx, req, resp)
}
