package provider

import (
	"errors"
	"testing"
)

func TestLookupEnvironmentVariable(t *testing.T) {
	t.Setenv("TF_PROVIDER_ENV_LOOKUP_FOUND", "found-value")
	t.Setenv("TF_PROVIDER_ENV_LOOKUP_EMPTY", "")

	testCases := []struct {
		name          string
		lookupName    string
		expectedValue string
		expectedError error
	}{
		{
			name:          "found variable returns value",
			lookupName:    "TF_PROVIDER_ENV_LOOKUP_FOUND",
			expectedValue: "found-value",
		},
		{
			name:          "found empty variable returns empty value",
			lookupName:    "TF_PROVIDER_ENV_LOOKUP_EMPTY",
			expectedValue: "",
		},
		{
			name:          "missing variable returns canonical error",
			lookupName:    "TF_PROVIDER_ENV_LOOKUP_MISSING",
			expectedValue: "",
			expectedError: missingVariableLookupError,
		},
		{
			name:          "empty variable name returns validation error",
			lookupName:    "",
			expectedValue: "",
			expectedError: invalidVariableLookupError,
		},
		{
			name:          "whitespace padded variable name returns validation error",
			lookupName:    " TF_PROVIDER_ENV_LOOKUP_FOUND ",
			expectedValue: "",
			expectedError: invalidVariableLookupError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualValue, actualErr := lookupEnvironmentVariable(testCase.lookupName)

			if actualValue != testCase.expectedValue {
				t.Fatalf("expected value %q, got %q", testCase.expectedValue, actualValue)
			}

			if testCase.expectedError == nil && actualErr != nil {
				t.Fatalf("expected no error, got %#v", actualErr)
			}

			if testCase.expectedError != nil && !errors.Is(actualErr, testCase.expectedError) {
				t.Fatalf("expected error %#v, got %#v", testCase.expectedError, actualErr)
			}
		})
	}
}
