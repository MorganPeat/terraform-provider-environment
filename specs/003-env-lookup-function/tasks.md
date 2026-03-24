# Tasks: Environment Lookup Function

**Input**: Design documents from `/specs/003-env-lookup-function/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Included. The specification and constitution require acceptance-style parity coverage across all in-scope lookup methods (non-sensitive data source, sensitive data source, function), with compatibility coverage retained for both data sources.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Terraform provider code: `internal/provider/`
- Generated/provider docs: `docs/`
- Provider examples: `examples/`

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create feature scaffolding and test harness files needed by all stories.

- [X] T001 Create shared lookup scaffolding and canonical error constants in `internal/provider/lookup.go`
- [X] T002 [P] Create shared lookup helper test scaffold in `internal/provider/lookup_test.go`
- [X] T003 [P] Create provider-defined function scaffolding and test scaffold in `internal/provider/variable_function.go` and `internal/provider/variable_function_test.go`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core behavior centralization that MUST be complete before user story implementation.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [X] T004 Implement shared raw-input lookup behavior (`os.LookupEnv`, success/error mapping) in `internal/provider/lookup.go`
- [X] T005 [P] Refactor non-sensitive data source read path to call shared helper and preserve `id=name` in `internal/provider/variable_data_source.go`
- [X] T006 [P] Refactor sensitive data source read path to call the shared helper while preserving existing delegation/parity semantics in `internal/provider/sensitive_variable_data_source.go`
- [X] T007 Register provider-defined function factories on the provider in `internal/provider/provider.go`
- [X] T008 Add foundational regression tests for canonical missing-variable message and helper behavior across function and both data sources in `internal/provider/lookup_test.go`, `internal/provider/variable_data_source_test.go`, and `internal/provider/sensitive_variable_data_source_test.go`

**Checkpoint**: Shared lookup rule and provider function wiring are ready for story work.

---

## Phase 3: User Story 1 - Read environment values in expressions (Priority: P1) 🎯 MVP

**Goal**: Deliver a provider-defined function that returns one environment variable value for expression usage.

**Independent Test**: Call `provider::environment::variable("PATH")` in a Terraform configuration and verify the output equals shell `PATH`; call with a missing key and verify lookup failure.

### Tests for User Story 1

- [X] T009 [P] [US1] Add acceptance-style success and missing-variable function scenarios in `internal/provider/variable_function_test.go`

### Implementation for User Story 1

- [X] T010 [US1] Implement function signature/definition for `provider::environment::variable(name)` in `internal/provider/variable_function.go`
- [X] T011 [US1] Wire function evaluation to the shared lookup helper and return string results in `internal/provider/variable_function.go`
- [X] T012 [US1] Add or update provider test config helpers for function scenarios in `internal/provider/provider_test.go` and `internal/provider/variable_function_test.go`

**Checkpoint**: User Story 1 is fully functional and testable independently.

---

## Phase 4: User Story 2 - Get consistent behavior across access methods (Priority: P2)

**Goal**: Ensure non-sensitive data source, sensitive data source, and function outcomes are parity-locked for present, missing, and empty-value cases with identical missing-variable text.

**Independent Test**: Run parity scenarios for present/missing/empty inputs across all in-scope lookup methods and confirm value and error outputs match exactly.

### Tests for User Story 2

- [X] T013 [P] [US2] Add cross-method parity tests for present and empty-value cases across function and both data sources in `internal/provider/lookup_parity_test.go`
- [X] T014 [P] [US2] Add cross-method parity tests for missing variable, empty-name, and whitespace-name with exact diagnostics checks across function and both data sources in `internal/provider/lookup_parity_test.go`

### Implementation for User Story 2

- [X] T015 [US2] Enforce one canonical missing-variable message constant usage in `internal/provider/lookup.go`, `internal/provider/variable_data_source.go`, and `internal/provider/variable_function.go`
- [X] T016 [US2] Ensure non-sensitive data source, sensitive data source, and function delegate all lookup/error mapping to the shared helper in `internal/provider/variable_data_source.go`, `internal/provider/sensitive_variable_data_source.go`, and `internal/provider/variable_function.go`
- [X] T017 [US2] Extend acceptance-style coverage for both data sources to guard compatibility during centralization in `internal/provider/variable_data_source_test.go` and `internal/provider/sensitive_variable_data_source_test.go`

**Checkpoint**: User Stories 1 and 2 behave consistently and are independently testable.

---

## Phase 5: User Story 3 - Keep sensitive handling simple (Priority: P3)

**Goal**: Clearly document that there is no sensitive function and direct secret use cases to the existing sensitive data source and Terraform sensitivity handling.

**Independent Test**: Review generated docs and examples to confirm one unambiguous recommendation path for sensitive use cases.

### Implementation for User Story 3

- [X] T018 [US3] Add explicit sensitive-handling guidance to function Markdown description in `internal/provider/variable_function.go`
- [X] T019 [P] [US3] Update provider-level docs text for function vs sensitive data source guidance in `internal/provider/provider.go` and `docs/index.md`
- [X] T020 [P] [US3] Update examples to show function usage and sensitive handling recommendations in `examples/provider/provider.tf` and `examples/data-sources/environment_sensitive_variable/data-source.tf`
- [X] T021 [US3] Regenerate documentation to produce/update function docs in `docs/functions/variable.md` via `main.go` (`go:generate`)

**Checkpoint**: All three user stories are complete and independently verifiable.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final quality, generation, and validation checks spanning all stories.

- [X] T022 [P] Run formatting and generation for touched files in `internal/provider/*.go`, `examples/**/*.tf`, and `docs/**/*.md` (`make fmt` and `go generate`)
- [X] T023 Run full verification for this feature using `GNUmakefile` targets (`make build` and `make test`) and resolve regressions in touched files
- [X] T024 Perform final quickstart/doc consistency pass in `specs/003-env-lookup-function/quickstart.md` and `docs/functions/variable.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies; starts immediately.
- **Phase 2 (Foundational)**: Depends on Phase 1; blocks all user stories.
- **Phase 3 (US1)**: Depends on Phase 2; delivers MVP.
- **Phase 4 (US2)**: Depends on Phase 3 (function exists) and Phase 2.
- **Phase 5 (US3)**: Depends on Phase 3 (function contract established).
- **Phase 6 (Polish)**: Depends on completion of desired user stories.

### User Story Dependencies

- **US1 (P1)**: First deliverable after foundation; no dependency on other user stories.
- **US2 (P2)**: Requires US1 function surface to exist for parity checks.
- **US3 (P3)**: Can begin after US1 naming/behavior are stable; independent of US2 parity internals.

### Within Each User Story

- Test tasks are defined before implementation tasks.
- Shared behavior tasks must be complete before parity assertions are finalized.
- Story-level checkpoint must pass before moving to next priority.

### Parallel Opportunities

- **Setup**: T002 and T003 can run in parallel after T001.
- **Foundational**: T005 and T006 can run in parallel after T004.
- **US1**: T009 can run while implementation is prepared, but must pass with T010-T011 complete.
- **US2**: T013 and T014 can run in parallel; T015 and T016 can proceed once failing parity tests are established.
- **US3**: T019 and T020 can run in parallel after T018.
- **Polish**: T022 can run in parallel with final doc review preparation for T024.

---

## Parallel Example: User Story 2

```bash
# Run parity test authoring tasks in parallel:
Task: "T013 Add cross-method parity tests for present and empty-value cases in internal/provider/lookup_parity_test.go"
Task: "T014 Add cross-method parity tests for missing/empty-name/whitespace cases in internal/provider/lookup_parity_test.go"

# Then implement shared behavior wiring updates:
Task: "T015 Enforce canonical missing-variable message constant usage across lookup surfaces"
Task: "T016 Delegate both data sources and function lookup/error mapping to shared helper"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1 (Setup).
2. Complete Phase 2 (Foundational centralization and wiring).
3. Complete Phase 3 (US1 function implementation + tests).
4. Validate MVP via US1 independent test and `make test`.

### Incremental Delivery

1. Build foundation once (Phases 1-2).
2. Deliver US1 as MVP.
3. Add US2 parity hardening and exact-message guarantees.
4. Add US3 documentation/sensitive-guidance updates.
5. Run Phase 6 polish and full verification.

### Parallel Team Strategy

1. One engineer handles shared helper and provider wiring (Phase 2).
2. One engineer implements function surface and tests (US1).
3. One engineer drives parity matrix tests/refactors (US2) after US1 baseline lands.
4. One engineer updates docs/examples and generation pipeline (US3).

---

## Notes

- `[P]` tasks touch different files and have no unfinished direct dependency.
- Every user story task includes the `[US#]` label for traceability.
- Keep missing-variable error text byte-for-byte identical across methods.
- Keep lookup inputs raw; do not trim or add extra validation branches.
