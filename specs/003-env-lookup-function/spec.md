# Feature Specification: Environment Lookup Function

**Feature Branch**: `003-env-lookup-function`  
**Created**: 2026-03-24  
**Status**: Draft  
**Input**: User description: "create a spec for a provider-defined function that mimics the current environment variable lookup data source. decide whether a separate sensitive function is needed. aim to centralise the environment variable lookup function between all three methods (non-sensitive data source, sensitive data source, function). keep it simple."

## Clarifications

### Session 2026-03-24

- Q: How should variable names with leading/trailing whitespace be handled? → A: Keep current raw behavior (no extra validation).
- Q: How strict should missing-variable error parity be across methods? → A: Require byte-for-byte identical error messages across non-sensitive data source, sensitive data source, and function.
- Q: How should an empty variable name input be handled? → A: Perform lookup with empty name as-is; if not found, return the standard missing-variable error.
- Q: How should implementation-level parity be enforced across methods? → A: Non-sensitive data source, sensitive data source, and function must all call one shared lookup helper for validation/lookup/error mapping.
- Q: What is the canonical sensitive-handling path for function-based lookups? → A: No sensitive function; documentation must direct sensitive use cases to sensitive data source and Terraform `sensitive(...)` handling.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Read environment values in expressions (Priority: P1)

As a Terraform user, I want to look up one environment variable through a provider-defined function so I can use the value directly in expressions without creating extra configuration objects.

**Why this priority**: This delivers the core user value: simple environment lookup where users need it most.

**Independent Test**: Define a plan that calls the function with an existing variable name and confirm the returned value matches the shell environment.

**Acceptance Scenarios**:

1. **Given** an environment variable exists, **When** the user calls the provider-defined lookup function with that name, **Then** the function returns the variable value.
2. **Given** an environment variable does not exist, **When** the user calls the provider-defined lookup function, **Then** the user receives a clear lookup failure message.

---

### User Story 2 - Get consistent behavior across access methods (Priority: P2)

As a provider maintainer, I want lookup behavior to be consistent across non-sensitive data source, sensitive data source, and function access methods so users get predictable outcomes regardless of entry point.

**Why this priority**: Inconsistent behavior creates confusion, support burden, and migration risk.

**Independent Test**: Evaluate the same set of variable names (present, missing, empty value) across each access method and confirm outcomes are aligned.

**Acceptance Scenarios**:

1. **Given** the same variable name and shell state, **When** a user retrieves it through any supported method, **Then** the success result and failure behavior are equivalent in meaning.
2. **Given** the same variable name and shell state, **When** a user retrieves it through any supported method, **Then** success values and missing-variable error messages are byte-for-byte identical.
3. **Given** lookup behavior changes in the future, **When** maintainers update the provider, **Then** all three in-scope lookup methods reflect that same behavior change.

---

### User Story 3 - Keep sensitive handling simple (Priority: P3)

As a Terraform user handling secrets, I want clear guidance on whether I need a separate sensitive function so I can choose a safe and simple usage pattern.

**Why this priority**: Sensitive data handling is important, but this can remain simple if the decision is explicit and documented.

**Independent Test**: Review provider documentation and examples to confirm there is a single, unambiguous recommendation for sensitive function usage.

**Acceptance Scenarios**:

1. **Given** this feature is released, **When** a user reviews the lookup function documentation, **Then** it clearly states that a separate sensitive function is not provided and explains the intended handling path for sensitive values.

---

### Edge Cases

- Variable exists but has an empty string value.
- Variable name contains mixed case, underscores, or digits.
- Variable name input is empty (lookup is attempted as-is; missing result uses the standard missing-variable error).
- Variable is unset after a previous successful read.
- Users pass a name with leading or trailing whitespace (treated as part of the variable name; no trimming or validation is applied).
- Users expect sensitive redaction but use the non-sensitive lookup path.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The provider MUST offer a provider-defined function that accepts one environment variable name and returns its value.
- **FR-002**: The function MUST match existing lookup semantics for success and missing-variable failure from the current environment variable lookup behavior.
- **FR-003**: The provider MUST maintain one shared lookup behavior definition that is used by non-sensitive data source, sensitive data source, and function methods.
- **FR-004**: For an existing variable with an empty value, all methods MUST return an empty value rather than a not-found error.
- **FR-005**: For a missing variable, all methods MUST return a clear, user-facing error that indicates the variable is not present.
- **FR-006**: The feature MUST NOT introduce a separate sensitive lookup function in this scope.
- **FR-007**: The provider documentation MUST describe how to handle sensitive values when using the function and how this differs from the sensitive data source path, explicitly directing sensitive use cases to the sensitive data source and Terraform `sensitive(...)` handling.
- **FR-008**: Existing data source behavior and compatibility MUST remain unchanged for current users.
- **FR-009**: Lookup inputs MUST be used as provided, without trimming or additional name validation, so parity remains aligned with current behavior.
- **FR-010**: Missing-variable failures MUST use one byte-for-byte identical user-facing error message across non-sensitive data source, sensitive data source, and function methods.
- **FR-011**: An empty variable name MUST follow the same raw lookup path and standard missing-variable error behavior as any other missing lookup, with no dedicated validation branch.
- **FR-012**: Non-sensitive data source, sensitive data source, and function methods MUST delegate lookup input handling, `os.LookupEnv` evaluation, and error/value mapping to one shared helper implementation.

### Key Entities *(include if feature involves data)*

- **Lookup Request**: User-provided variable name used to query the current shell environment.
- **Lookup Result**: Retrieved value or a not-found failure outcome for that request.
- **Lookup Behavior Rule**: Shared policy that defines how success, missing variables, and empty values are interpreted across all access methods.

### Assumptions

- The provider continues to support explicit per-variable lookup rather than bulk environment export.
- Non-sensitive data source, sensitive data source, and function methods are all considered in scope for behavior consistency.
- Sensitive use cases continue to be supported by existing sensitive data source patterns and Terraform-level sensitive handling, rather than adding a second function.

### Dependencies

- Provider documentation updates for function usage and sensitive handling guidance.
- Validation coverage that checks parity across all applicable access methods.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: In parity validation scenarios (present, missing, and empty variables), 100% of outcomes are identical across non-sensitive data source, sensitive data source, and function methods, including missing-variable error text.
- **SC-002**: 100% of missing-variable lookups produce a clear user-facing error indicating the variable is not present.
- **SC-003**: At least 90% of users in documentation walkthrough feedback can correctly choose the recommended path for sensitive vs non-sensitive lookups on first attempt.
- **SC-004**: Feature review sign-off records zero open questions about whether a separate sensitive function exists or is required.
