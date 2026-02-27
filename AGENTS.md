# AGENTS.md

This file provides guidance to agents when working with code in this repository.

## Project

Terraform provider (Go) using `terraform-plugin-framework` v1. Exposes shell environment variables as Terraform data sources. Provider address: `registry.terraform.io/morganpeat/environment`.

## Commands

```bash
make build       # go build -v ./...
make test        # go test -v -cover -timeout=120s -parallel=4 ./...
make testacc     # TF_ACC=1 go test -v -cover -timeout 120m ./... (acceptance tests)
make lint        # golangci-lint run
make fmt         # gofmt -s -w -e .
make generate    # go generate ./... (formats examples + regenerates docs via tfplugindocs)
```

Run a single test:
```bash
go test -v -run TestAccEnvironmentVariableDataSource ./internal/provider/
```

## Key Patterns

- **Tests use `IsUnitTest: true`** — acceptance tests run without `TF_ACC=1` because of this flag; `make test` runs them all.
- **`sensitiveVariableDataSource.Read` delegates to `NewVariableDataSource().Read()`** — the only difference between the two data sources is `Sensitive: true` on the `value` schema attribute.
- **`id` is always set equal to `name`** — this is the convention for both data sources.
- **`go generate`** runs two things: `terraform fmt -recursive ./examples/` and `tfplugindocs` — run it after changing schemas or examples.
- **Linter excludes `examples/`** — golangci-lint and gofmt skip the `examples$` path.
- **No provider configuration block** — `Configure()` is a no-op; the provider takes no config.
- All code lives in `internal/provider/` as a single package `provider`.