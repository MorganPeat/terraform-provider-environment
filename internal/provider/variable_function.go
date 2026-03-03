package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = &variableFunction{}

type variableFunction struct{}

func NewVariableFunction() function.Function {
	return &variableFunction{}
}

func (f *variableFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "variable"
}

func (f *variableFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Read an environment variable",
		Description: "Returns the value of the specified environment variable.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "name",
				Description: "The name of the environment variable to read.",
			},
		},
		Return: function.StringReturn{},
	}
}

// Run executes the function logic, reading the environment variable
// and returning its value or an error if not found.
func (f *variableFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var name string
	resp.Error = req.Arguments.Get(ctx, &name)
	if resp.Error != nil {
		return
	}

	value, ok := os.LookupEnv(name)
	if !ok {
		resp.Error = function.NewArgumentFuncError(
			0,
			fmt.Sprintf("Environment variable %q not found", name),
		)
		return
	}

	resp.Error = resp.Result.Set(ctx, value)
}

// Made with Bob
