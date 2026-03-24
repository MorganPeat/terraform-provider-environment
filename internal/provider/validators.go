package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type environmentVariableNameValidator struct{}

func (v environmentVariableNameValidator) Description(ctx context.Context) string {
	return canonicalInvalidVariableError
}

func (v environmentVariableNameValidator) MarkdownDescription(ctx context.Context) string {
	return canonicalInvalidVariableError
}

func (v environmentVariableNameValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if err := validateEnvironmentVariableName(req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, invalidVariableLookupErrorSummary, canonicalInvalidVariableError)
	}
}

func (v environmentVariableNameValidator) ValidateParameterString(ctx context.Context, req function.StringParameterValidatorRequest, resp *function.StringParameterValidatorResponse) {
	if req.Value.IsNull() || req.Value.IsUnknown() {
		return
	}

	if err := validateEnvironmentVariableName(req.Value.ValueString()); err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewArgumentFuncError(req.ArgumentPosition, canonicalInvalidVariableError))
	}
}
