<!--
SYNC IMPACT REPORT
Version change: 0.0.0 -> 1.0.0
Modified principles:
- Initial ratification of constitution
Added sections:
- Core Principles (Framework First, Acceptance Testing, Documentation Quality, Idiomatic Go, Security & Sensitivity)
- Provider Design
- Development Standards
Templates requiring updates:
- .specify/templates/plan-template.md (⚠ pending)
- .specify/templates/spec-template.md (⚠ pending)
- .specify/templates/tasks-template.md (⚠ pending)
Follow-up: None
-->
# Terraform Provider Environment Constitution

## Core Principles

### I. Framework First
The provider MUST be built using the `terraform-plugin-framework`. All resources and data sources must implement the `resource.Resource` or `datasource.DataSource` interfaces. Avoid using the legacy SDK unless strictly necessary for backward compatibility with pre-existing resources (not applicable for new providers).

### II. Acceptance Testing (Non-Negotiable)
Every resource and data source MUST have acceptance tests (`TestAcc...`) that run against real Terraform execution. Tests must verify the full lifecycle: create, read, update, and delete. Mocking is discouraged for acceptance tests; use real environment interactions where possible.

### III. Documentation Quality
All schemas (attributes, blocks) MUST have clear, user-centric `MarkdownDescription`s. Documentation is generated via `tfplugindocs`. Examples MUST be functional and copy-pasteable. If a field is sensitive, it MUST be documented as such.

### IV. Idiomatic Go
Code MUST follow standard Go conventions. Contexts (`context.Context`) MUST be propagated through all request pipelines. Errors MUST be wrapped with context (e.g., `fmt.Errorf("reading environment variable: %w", err)`). Avoid global state.

### V. Security & Sensitivity
Since this provider exposes environment variables, sensitive data handling is critical. Any attribute that might contain secrets MUST be marked `Sensitive: true`. The provider must never log sensitive values in plain text.

## Provider Design

Provider configuration should be minimal. Where possible, configuration should be sourced from the environment or standard Terraform conventions. The provider must support `terraform-registry-manifest.json` for proper registry publishing.

## Development Standards

Code must pass `golangci-lint` with standard presets. `go fmt` is mandatory. Commit messages should follow conventional commits (feat, fix, docs, chore).

## Governance

This constitution governs all development on the Terraform Provider Environment. Amendments require a Pull Request with justification. All code reviews must verify compliance with these principles.

**Version**: 1.0.0 | **Ratified**: 2026-03-23 | **Last Amended**: 2026-03-23
