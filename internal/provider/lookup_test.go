package provider

import (
	"testing"
)

func TestLookupEnvironmentVariable(t *testing.T) {
	t.Setenv("TF_PROVIDER_ENV_LOOKUP_FOUND", "found-value")
	t.Setenv("TF_PROVIDER_ENV_LOOKUP_EMPTY", "")

	testCases := []struct {
		name          string
		lookupName    string
		expectedValue string
		expectedError string
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
			name:          "missing variable returns not found error",
			lookupName:    "TF_PROVIDER_ENV_LOOKUP_MISSING",
			expectedValue: "",
			expectedError: missingVariableErrorMessage("TF_PROVIDER_ENV_LOOKUP_MISSING"),
		},
		{
			name:          "empty variable name returns not found error",
			lookupName:    "",
			expectedValue: "",
			expectedError: missingVariableErrorMessage(""),
		},
		{
			name:          "whitespace padded variable name returns not found error",
			lookupName:    " TF_PROVIDER_ENV_LOOKUP_FOUND ",
			expectedValue: "",
			expectedError: missingVariableErrorMessage(" TF_PROVIDER_ENV_LOOKUP_FOUND "),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			actualValue, actualErr := lookupEnvironmentVariable(testCase.lookupName)

			if actualValue != testCase.expectedValue {
				t.Fatalf("expected value %q, got %q", testCase.expectedValue, actualValue)
			}

			if testCase.expectedError == "" && actualErr != nil {
				t.Fatalf("expected no error, got %#v", actualErr)
			}

			if testCase.expectedError != "" {
				if actualErr == nil {
					t.Fatalf("expected error %q, got nil", testCase.expectedError)
				}

				if actualErr.Error() != testCase.expectedError {
					t.Fatalf("expected error %q, got %q", testCase.expectedError, actualErr.Error())
				}
			}
		})
	}
}
