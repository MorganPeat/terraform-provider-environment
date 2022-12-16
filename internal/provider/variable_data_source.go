package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var _ datasource.DataSource = &variableDataSource{}

// variableDataSource is the data source implementation.
type variableDataSource struct{}

// variableDataSourceModel describes the data model for this data source.
type variableDataSourceModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

// Metadata returns the data source type name.
func (d *variableDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable"
}

// GetSchema defines the schema for the data source.

func (d *variableDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
The variable data source exposes a shell environment variable to terraform.

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
				MarkdownDescription: "The value of the shell environment variable.",
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *variableDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data variableDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	v, ok := os.LookupEnv(data.Name.ValueString())
	if !ok {
		resp.Diagnostics.AddError(
			"Not found",
			"The environment variable is not present.",
		)
	}

	data.ID = data.Name
	data.Value = types.StringValue(v)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

// NewVariableDataSource is a helper function to simplify the provider implementation.
func NewVariableDataSource() datasource.DataSource {
	return &variableDataSource{}
}
