package provider

import (
	"errors"
	"os"
	"strings"
)

const (
	missingVariableLookupErrorSummary = "Not found"
	canonicalMissingVariableError     = "The environment variable is not present"
	invalidVariableLookupErrorSummary = "Invalid name"
	canonicalInvalidVariableError     = "The environment variable name must not be empty and must not include leading or trailing whitespace"
	lookupVariableErrorSummary        = "Lookup error"
)

type lookupError struct {
	Summary string
	Detail  string
}

func (e *lookupError) Error() string {
	return e.Detail
}

var missingVariableLookupError = &lookupError{
	Summary: missingVariableLookupErrorSummary,
	Detail:  canonicalMissingVariableError,
}

var invalidVariableLookupError = &lookupError{
	Summary: invalidVariableLookupErrorSummary,
	Detail:  canonicalInvalidVariableError,
}

func validateEnvironmentVariableName(name string) error {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" || trimmedName != name {
		return invalidVariableLookupError
	}

	return nil
}

func lookupEnvironmentVariable(name string) (string, error) {
	if err := validateEnvironmentVariableName(name); err != nil {
		return "", err
	}

	value, found := os.LookupEnv(name)
	if !found {
		return "", missingVariableLookupError
	}

	return value, nil
}

func lookupErrorSummaryAndDetail(err error) (string, string) {
	var typedErr *lookupError
	if errors.As(err, &typedErr) {
		return typedErr.Summary, typedErr.Detail
	}

	return lookupVariableErrorSummary, err.Error()
}
