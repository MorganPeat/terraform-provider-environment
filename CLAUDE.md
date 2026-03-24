# terraform-provider-environment Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-03-24

## Active Technologies
- Go 1.25.0 + `terraform-plugin-framework` v1.19.0, `terraform-plugin-testing` v1.15.0, `terraform-plugin-docs` v0.24.0 (003-env-lookup-function)
- N/A (reads process environment via `os.LookupEnv`; values flow into Terraform state where applicable) (003-env-lookup-function)

- Go 1.25.0 + `terraform-plugin-framework` v1.19.0, `terraform-plugin-testing` v1.15.0, Go standard library (`os`) (003-env-lookup-function)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Go 1.25.0

## Code Style

Go 1.25.0: Follow standard conventions

## Recent Changes
- 003-env-lookup-function: Added Go 1.25.0 + `terraform-plugin-framework` v1.19.0, `terraform-plugin-testing` v1.15.0, `terraform-plugin-docs` v0.24.0

- 003-env-lookup-function: Added Go 1.25.0 + `terraform-plugin-framework` v1.19.0, `terraform-plugin-testing` v1.15.0, Go standard library (`os`)

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
