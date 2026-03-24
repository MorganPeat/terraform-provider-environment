package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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

// Schema defines the schema for the data source.
func (d *variableDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `
The variable data source exposes a shell environment variable to terraform.

Any change in the value of the shell environment variable will show up as a change in the terraform plan.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for this data source instance. This matches the name of the environment variable.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the shell environment variable to read. This must not be empty and must not include leading or trailing whitespace.",
				Validators: []validator.String{
					environmentVariableNameValidator{},
				},
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

	value, ok := readEnvironmentVariableDataSourceValue(data.Name, &resp.Diagnostics)
	if !ok {
		return
	}

	data.ID = data.Name
	data.Value = types.StringValue(value)
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func readEnvironmentVariableDataSourceValue(name types.String, diags *diag.Diagnostics) (string, bool) {
	if name.IsNull() || name.IsUnknown() {
		diags.AddAttributeError(
			path.Root("name"),
			invalidVariableLookupErrorSummary,
			canonicalInvalidVariableError,
		)
		return "", false
	}

	value, err := lookupEnvironmentVariable(name.ValueString())
	if err != nil {
		summary, detail := lookupErrorSummaryAndDetail(err)
		diags.AddAttributeError(path.Root("name"), summary, detail)
		return "", false
	}

	return value, true
}

// NewVariableDataSource is a helper function to simplify the provider implementation.
func NewVariableDataSource() datasource.DataSource {
	return &variableDataSource{}
}
