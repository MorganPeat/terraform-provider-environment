package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ function.Function = &variableFunction{}

type variableFunction struct{}

func NewVariableFunction() function.Function {
	return &variableFunction{}
}

func (f *variableFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "variable"
}

func (f *variableFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Read one shell environment variable by name.",
		MarkdownDescription: `
Read one shell environment variable by name and return its value.

Provider-defined functions require Terraform 1.8 or later.

This function is non-sensitive. For secret use cases, use the ` + "`environment_sensitive_variable` data source" + ` and Terraform's ` + "`sensitive(...)` handling" + `.
`,
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "name",
				MarkdownDescription: "The name of the shell environment variable to read. This must not be empty and must not include leading or trailing whitespace.",
				Validators: []function.StringParameterValidator{
					environmentVariableNameValidator{},
				},
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *variableFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var name string

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &name))
	if resp.Error != nil {
		return
	}

	value, err := lookupEnvironmentVariable(name)
	if err != nil {
		if typedErr, ok := err.(*lookupError); ok {
			resp.Error = function.ConcatFuncErrors(resp.Error, function.NewArgumentFuncError(0, typedErr.Detail))
			return
		}

		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(err.Error()))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, basetypes.NewStringValue(value)))
}
