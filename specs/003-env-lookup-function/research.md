# Research: Environment Lookup Function

## Decision 1: Single provider-defined function without a sensitive variant

- Decision: Introduce one provider-defined lookup function and do not add a second sensitive function.
- Rationale: The spec requires a single function and directs sensitive use cases to the existing sensitive data source plus Terraform `sensitive(...)` handling. This avoids API duplication and user confusion.
- Alternatives considered:
  - Add `environment_sensitive_variable(...)` function: rejected because it duplicates behavior and creates two function paths for one lookup operation.
  - Add sensitive-mode parameter to the function: rejected because this still expands function surface and blurs schema-level sensitivity vs expression-level sensitivity.

## Decision 2: Centralize lookup semantics in one shared helper

- Decision: Create one shared internal lookup helper for name input passthrough, `os.LookupEnv` evaluation, value mapping, and missing-variable error mapping.
- Rationale: FR-012 requires one implementation path for all methods. This ensures byte-for-byte parity and reduces future drift.
- Alternatives considered:
  - Keep duplicated lookup logic in each surface: rejected due to drift risk and parity failures.
  - Shared constant for errors but separate lookup code: rejected because semantic drift can still occur in edge cases.

## Decision 3: Keep raw-name behavior exactly as current

- Decision: Use input names exactly as provided (including whitespace and empty string) and perform lookup directly.
- Rationale: Clarifications and FR-009/FR-011 require strict parity with current behavior and no additional validation branch.
- Alternatives considered:
  - Trim names before lookup: rejected because this changes current behavior and can produce mismatches.
  - Reject empty names early: rejected because spec requires standard missing-variable path.

## Decision 4: Enforce parity with acceptance-style tests across methods

- Decision: Add/extend acceptance-style tests (`IsUnitTest: true`) covering present, missing, and empty-value cases for data source and new function, including exact missing-variable message parity.
- Rationale: Constitution requires acceptance testing. The spec requires byte-for-byte missing-variable consistency and full behavior parity.
- Alternatives considered:
  - Unit tests only for helper function: rejected because provider integration and diagnostic behavior would be under-tested.
  - Snapshot docs-only verification: rejected as non-executable and insufficient for regression prevention.

## Decision 5: Document sensitive path explicitly in function docs

- Decision: Update function docs/examples and provider docs to explicitly direct sensitive use cases to the sensitive data source and Terraform `sensitive(...)` wrappers where appropriate.
- Rationale: FR-007 requires unambiguous guidance and explicit differentiation from sensitive data source behavior.
- Alternatives considered:
  - Mention sensitivity briefly without explicit guidance: rejected as ambiguous.
  - Defer sensitive guidance to external docs only: rejected because local provider docs should remain self-contained and actionable.
