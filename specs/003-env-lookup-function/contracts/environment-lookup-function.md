# Contract: Provider-Defined Environment Lookup Function

## Overview

Defines the behavior contract for a provider-defined function that retrieves one shell environment variable by name and mirrors existing provider lookup semantics.

## Function Signature

- Canonical name: `provider::environment::variable(name)`
- Inputs:
  - `name` (string, required): raw environment variable name.
- Output:
  - string value of the environment variable when found.
- Error:
  - canonical missing-variable user-facing error when not found.

## Behavioral Requirements

1. Input passthrough
   - `name` is used exactly as provided.
   - No trimming, normalization, or extra validation.

2. Lookup execution
   - Uses shared lookup behavior with `os.LookupEnv` semantics.

3. Success mapping
   - If `os.LookupEnv(name)` is found, return value as-is.
   - Empty string value is a valid success result.

4. Failure mapping
   - If variable is missing, emit the same user-facing error text used by existing data source lookup.
   - Error text parity must be byte-for-byte identical across methods.

5. Sensitive behavior
   - No separate sensitive function is exposed.
   - Sensitive use cases are supported through `environment_sensitive_variable` data source and Terraform `sensitive(...)` handling.

## Cross-Method Parity Matrix

| Scenario | `environment_variable` data source | lookup function | Parity Requirement |
|---|---|---|---|
| Variable exists, non-empty value | Return value | Return value | Exact value match |
| Variable exists, empty value | Return empty string | Return empty string | Exact value and success-path match |
| Variable missing | Emit canonical missing error | Emit same canonical missing error | Byte-for-byte identical message |
| Empty input name, not found | Emit canonical missing error | Emit same canonical missing error | Byte-for-byte identical message |
| Whitespace in name | Treated as literal name | Treated as literal name | Exact lookup behavior match |

## Acceptance Contract Tests

Minimum contract assertions:

- Existing variable returns same value for data source and function.
- Missing variable error text matches exactly across data source and function.
- Empty value variable succeeds for both methods with empty result.
- Empty-name lookup follows standard missing-variable path.
- Whitespace-containing names are treated literally.
