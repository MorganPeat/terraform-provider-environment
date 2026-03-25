package provider

import (
	"context"
	"fmt"
	"os"

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

The function looks up the exact name you pass, including any leading or trailing whitespace.

Provider-defined functions require Terraform 1.8 or later.
`,
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "name",
				MarkdownDescription: "The name of the shell environment variable to read. This name is looked up exactly as provided, including any leading or trailing whitespace.",
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

	v, ok := os.LookupEnv(name)
	if !ok {
		resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("environment variable '%s' not found", name))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, basetypes.NewStringValue(v)))
}
