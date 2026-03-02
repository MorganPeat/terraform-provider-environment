# Implementation Checklist for Provider-Defined Functions

This checklist provides a step-by-step guide for implementing the provider-defined functions feature as specified in [`DESIGN_FUNCTIONS.md`](DESIGN_FUNCTIONS.md).

## Prerequisites

- [ ] Review [`DESIGN_FUNCTIONS.md`](DESIGN_FUNCTIONS.md) thoroughly
- [ ] Review [`ARCHITECTURE_DIAGRAM.md`](ARCHITECTURE_DIAGRAM.md) for visual understanding
- [ ] Ensure development environment is set up (Go 1.25+, Terraform 1.8+)
- [ ] Create feature branch: `git checkout -b feature/provider-functions`

## Phase 1: Core Implementation

### 1.1 Update Provider Interface

- [ ] Open [`internal/provider/provider.go`](internal/provider/provider.go)
- [ ] Add import: `"github.com/hashicorp/terraform-plugin-framework/function"`
- [ ] Add `Functions()` method to `environmentProvider`:
  ```go
  func (p *environmentProvider) Functions(ctx context.Context) []func() function.Function {
      return []func() function.Function{
          NewVariableFunction,
          NewSensitiveVariableFunction,
      }
  }
  ```
- [ ] Verify code compiles: `go build ./...`

### 1.2 Implement Variable Function

- [ ] Create [`internal/provider/variable_function.go`](internal/provider/variable_function.go)
- [ ] Implement `variableFunction` struct
- [ ] Implement `NewVariableFunction()` constructor
- [ ] Implement `Metadata()` method (name: "variable")
- [ ] Implement `Definition()` method with:
  - [ ] Summary and description
  - [ ] String parameter "name"
  - [ ] String return type (non-sensitive)
- [ ] Implement `Run()` method with:
  - [ ] Argument parsing
  - [ ] `os.LookupEnv()` call
  - [ ] Error handling for missing variables
  - [ ] Result setting
- [ ] Verify code compiles: `go build ./...`

### 1.3 Implement Sensitive Variable Function

- [ ] Create [`internal/provider/sensitive_variable_function.go`](internal/provider/sensitive_variable_function.go)
- [ ] Implement `sensitiveVariableFunction` struct
- [ ] Implement `NewSensitiveVariableFunction()` constructor
- [ ] Implement `Metadata()` method (name: "sensitive_variable")
- [ ] Implement `Definition()` method with:
  - [ ] Summary and description
  - [ ] String parameter "name"
  - [ ] String return type with `Sensitive: true`
- [ ] Implement `Run()` method (same logic as variable function)
- [ ] Verify code compiles: `go build ./...`

### 1.4 Unit Tests - Variable Function

- [ ] Create [`internal/provider/variable_function_test.go`](internal/provider/variable_function_test.go)
- [ ] Implement `TestVariableFunction_Success`:
  - [ ] Set test environment variable
  - [ ] Call function
  - [ ] Assert correct value returned
  - [ ] Clean up environment variable
- [ ] Implement `TestVariableFunction_NotFound`:
  - [ ] Ensure variable doesn't exist
  - [ ] Call function
  - [ ] Assert error is returned
  - [ ] Verify error message
- [ ] Implement `TestVariableFunction_EmptyValue`:
  - [ ] Set variable to empty string
  - [ ] Call function
  - [ ] Assert empty string returned (not error)
- [ ] Implement `TestVariableFunction_SpecialCharacters`:
  - [ ] Test with special characters in name and value
  - [ ] Assert correct handling
- [ ] Run tests: `go test -v ./internal/provider/ -run TestVariableFunction`
- [ ] Verify all tests pass

### 1.5 Unit Tests - Sensitive Variable Function

- [ ] Create [`internal/provider/sensitive_variable_function_test.go`](internal/provider/sensitive_variable_function_test.go)
- [ ] Implement `TestSensitiveVariableFunction_Success`:
  - [ ] Set test environment variable
  - [ ] Call function
  - [ ] Assert correct value returned
  - [ ] Verify value is marked as sensitive
  - [ ] Clean up environment variable
- [ ] Implement `TestSensitiveVariableFunction_NotFound`:
  - [ ] Ensure variable doesn't exist
  - [ ] Call function
  - [ ] Assert error is returned
- [ ] Implement `TestSensitiveVariableFunction_SensitivityMarking`:
  - [ ] Verify return type has `Sensitive: true`
- [ ] Run tests: `go test -v ./internal/provider/ -run TestSensitiveVariableFunction`
- [ ] Verify all tests pass

### 1.6 Verify Phase 1

- [ ] Run all unit tests: `make test`
- [ ] Run linter: `make lint`
- [ ] Run formatter: `make fmt`
- [ ] Verify no compilation errors: `make build`

## Phase 2: Acceptance Testing

### 2.1 Variable Function Acceptance Tests

- [ ] Add test to [`internal/provider/variable_function_test.go`](internal/provider/variable_function_test.go)
- [ ] Implement `TestAccVariableFunction_Basic`:
  - [ ] Set test environment variable
  - [ ] Create Terraform config using function in output
  - [ ] Verify output value matches
- [ ] Implement `TestAccVariableFunction_InLocals`:
  - [ ] Test function usage in locals block
- [ ] Implement `TestAccVariableFunction_InResource`:
  - [ ] Test function usage in resource configuration
- [ ] Implement `TestAccVariableFunction_NotFound`:
  - [ ] Test error handling with missing variable
  - [ ] Verify appropriate error message

### 2.2 Sensitive Variable Function Acceptance Tests

- [ ] Add test to [`internal/provider/sensitive_variable_function_test.go`](internal/provider/sensitive_variable_function_test.go)
- [ ] Implement `TestAccSensitiveVariableFunction_Basic`:
  - [ ] Set test environment variable
  - [ ] Create Terraform config using sensitive function
  - [ ] Verify value is marked as sensitive
- [ ] Implement `TestAccSensitiveVariableFunction_InOutput`:
  - [ ] Test with sensitive output
  - [ ] Verify masking behavior

### 2.3 Backward Compatibility Tests

- [ ] Verify existing data source tests still pass
- [ ] Test mixed usage (data sources + functions in same config)
- [ ] Verify data sources work independently

### 2.4 Verify Phase 2

- [ ] Run acceptance tests: `make testacc`
- [ ] Verify all tests pass
- [ ] Test with actual Terraform 1.8+ installation

## Phase 3: Documentation

### 3.1 Update Provider Documentation Template

- [ ] Open [`templates/index.md.tmpl`](templates/index.md.tmpl)
- [ ] Add "Functions" section after data sources section
- [ ] Document `variable` function with:
  - [ ] Signature
  - [ ] Description
  - [ ] Example usage
- [ ] Document `sensitive_variable` function with:
  - [ ] Signature
  - [ ] Description
  - [ ] Example usage
- [ ] Add "Migration from Data Sources to Functions" section
- [ ] Include comparison table (data sources vs functions)

### 3.2 Create Function Examples

- [ ] Create directory: `examples/functions/environment_variable/`
- [ ] Create [`examples/functions/environment_variable/function.tf`](examples/functions/environment_variable/function.tf):
  - [ ] Basic usage example
  - [ ] Usage in locals
  - [ ] Usage in resource
  - [ ] Add comments explaining Terraform 1.8+ requirement
- [ ] Create directory: `examples/functions/environment_sensitive_variable/`
- [ ] Create [`examples/functions/environment_sensitive_variable/function.tf`](examples/functions/environment_sensitive_variable/function.tf):
  - [ ] Sensitive variable example
  - [ ] Usage in sensitive output
  - [ ] Add comments about sensitivity

### 3.3 Update README

- [ ] Open [`README.md`](README.md)
- [ ] Add "Functions (Terraform 1.8+)" section after data source example
- [ ] Include basic function usage example
- [ ] Add note about Terraform version requirements
- [ ] Update requirements section if needed

### 3.4 Generate Documentation

- [ ] Run: `make generate`
- [ ] Verify docs are generated in `docs/functions/`:
  - [ ] `docs/functions/variable.md`
  - [ ] `docs/functions/sensitive_variable.md`
- [ ] Review generated documentation for accuracy
- [ ] Verify examples are formatted correctly

### 3.5 Verify Phase 3

- [ ] Review all documentation for clarity
- [ ] Test example configurations manually
- [ ] Verify links work correctly
- [ ] Check for typos and formatting issues

## Phase 4: Final Verification & Release Preparation

### 4.1 Comprehensive Testing

- [ ] Run full test suite: `make test`
- [ ] Run acceptance tests: `make testacc`
- [ ] Run linter: `make lint`
- [ ] Run formatter: `make fmt`
- [ ] Build provider: `make build`
- [ ] Test with local Terraform configuration
- [ ] Test both functions and data sources

### 4.2 Update CHANGELOG

- [ ] Open `CHANGELOG.md` (or create if doesn't exist)
- [ ] Add new version section
- [ ] Document new features:
  - [ ] Provider-defined functions support
  - [ ] `variable` function
  - [ ] `sensitive_variable` function
- [ ] Note backward compatibility maintained
- [ ] Reference issue #61

### 4.3 Create Release Notes

- [ ] Create release notes document
- [ ] Highlight new function features
- [ ] Include migration examples
- [ ] Note Terraform version requirements
- [ ] Emphasize backward compatibility

### 4.4 Code Review Preparation

- [ ] Review all changes
- [ ] Ensure code follows project conventions
- [ ] Verify all tests pass
- [ ] Check documentation completeness
- [ ] Prepare PR description with:
  - [ ] Summary of changes
  - [ ] Link to design document
  - [ ] Testing performed
  - [ ] Breaking changes (none expected)

### 4.5 Final Checks

- [ ] All unit tests pass
- [ ] All acceptance tests pass
- [ ] Linter passes with no warnings
- [ ] Code is formatted correctly
- [ ] Documentation is complete and accurate
- [ ] Examples run successfully
- [ ] Backward compatibility verified
- [ ] No breaking changes introduced

## Post-Implementation

### After Merge

- [ ] Tag release with appropriate version (minor bump)
- [ ] Publish release notes
- [ ] Update Terraform Registry (if applicable)
- [ ] Monitor for issues
- [ ] Respond to community feedback

### Future Enhancements (Optional)

- [ ] Consider bulk read function
- [ ] Consider default value support
- [ ] Consider validation functions
- [ ] Gather user feedback for improvements

## Notes

- **Terraform Version**: Functions require Terraform 1.8+
- **Framework Version**: Using terraform-plugin-framework v1.18.0
- **Backward Compatibility**: All existing functionality must continue to work
- **Testing**: Both unit and acceptance tests are required
- **Documentation**: Must be comprehensive and include examples

## References

- Design Document: [`DESIGN_FUNCTIONS.md`](DESIGN_FUNCTIONS.md)
- Architecture Diagrams: [`ARCHITECTURE_DIAGRAM.md`](ARCHITECTURE_DIAGRAM.md)
- Issue: [#61](https://github.com/MorganPeat/terraform-provider-environment/issues/61)
- Framework Docs: https://developer.hashicorp.com/terraform/plugin/framework/functions