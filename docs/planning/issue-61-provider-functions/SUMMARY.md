# Provider-Defined Functions Implementation Plan - Executive Summary

## Overview

This plan addresses [GitHub Issue #61](https://github.com/MorganPeat/terraform-provider-environment/issues/61), which requests implementing provider-defined functions for Terraform 1.8+ while maintaining backward compatibility with existing data sources.

## What Are Provider-Defined Functions?

Provider-defined functions (introduced in Terraform 1.8) allow simpler, more intuitive syntax for accessing provider functionality:

**Current (Data Source)**:
```hcl
data "environment_variable" "path" {
  name = "PATH"
}
output "path" {
  value = data.environment_variable.path.value
}
```

**Proposed (Function)**:
```hcl
output "path" {
  value = provider::environment::variable("PATH")
}
```

## Key Benefits

1. **Simpler Syntax**: Direct inline usage without resource blocks
2. **No State Overhead**: Functions don't store values in state
3. **Modern Terraform**: Leverages latest Terraform capabilities
4. **Backward Compatible**: Existing data sources remain unchanged

## Proposed Implementation

### Two Functions

1. **`variable(name string) string`**
   - Returns environment variable value
   - Non-sensitive return type
   - Example: `provider::environment::variable("PATH")`

2. **`sensitive_variable(name string) string`**
   - Returns environment variable value marked as sensitive
   - Sensitive return type
   - Example: `provider::environment::sensitive_variable("API_KEY")`

### Why Two Functions?

Terraform function return types must be statically defined. We cannot conditionally mark a value as sensitive at runtime, so we need separate functions for sensitive and non-sensitive values (mirroring the existing data source design).

## Architecture

### Current State
- ✅ `environment_variable` data source (Terraform 1.0+)
- ✅ `environment_sensitive_variable` data source (Terraform 1.0+)

### After Implementation
- ✅ `environment_variable` data source (Terraform 1.0+) - **unchanged**
- ✅ `environment_sensitive_variable` data source (Terraform 1.0+) - **unchanged**
- 🆕 `variable` function (Terraform 1.8+)
- 🆕 `sensitive_variable` function (Terraform 1.8+)

## Implementation Phases

### Phase 1: Core Implementation (2-3 hours)
- Add `Functions()` method to provider
- Implement `variable` function
- Implement `sensitive_variable` function
- Add comprehensive unit tests

### Phase 2: Testing (1-2 hours)
- Add acceptance tests
- Verify backward compatibility
- Test with Terraform 1.8+

### Phase 3: Documentation (1-2 hours)
- Update provider documentation
- Create function examples
- Update README
- Generate docs with `make generate`

### Phase 4: Release (1 hour)
- Update CHANGELOG
- Create release notes
- Tag release

**Total Estimated Time**: 5-8 hours

## Technical Details

### Framework Support
- Using `terraform-plugin-framework v1.18.0` ✅
- Full function support available ✅
- No framework upgrade needed ✅

### File Changes

**New Files** (6):
- `internal/provider/variable_function.go`
- `internal/provider/variable_function_test.go`
- `internal/provider/sensitive_variable_function.go`
- `internal/provider/sensitive_variable_function_test.go`
- `examples/functions/environment_variable/function.tf`
- `examples/functions/environment_sensitive_variable/function.tf`

**Modified Files** (3):
- `internal/provider/provider.go` (add `Functions()` method)
- `templates/index.md.tmpl` (add functions documentation)
- `README.md` (add functions example)

**No Breaking Changes**: All existing code remains unchanged.

## Testing Strategy

### Unit Tests
- ✅ Variable exists and returns correct value
- ✅ Variable doesn't exist, returns error
- ✅ Empty value handling
- ✅ Special characters in names/values
- ✅ Sensitive marking verification

### Acceptance Tests
- ✅ Function in output blocks
- ✅ Function in locals blocks
- ✅ Function in resource configurations
- ✅ Error handling
- ✅ Backward compatibility with data sources

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| Breaking existing functionality | Low | High | Comprehensive testing, no changes to existing code |
| Terraform version compatibility | Low | Medium | Clear documentation of version requirements |
| User confusion about when to use functions vs data sources | Medium | Low | Clear migration guide and comparison documentation |
| Performance differences | Low | Low | Document behavior differences (state vs on-demand) |

## Success Criteria

- ✅ Both functions work correctly with Terraform 1.8+
- ✅ All existing data sources continue to work
- ✅ All tests pass (unit + acceptance)
- ✅ Documentation is clear and comprehensive
- ✅ Examples run successfully
- ✅ No breaking changes introduced

## Documentation Deliverables

1. **[DESIGN_FUNCTIONS.md](DESIGN_FUNCTIONS.md)** (568 lines)
   - Comprehensive design specification
   - Technical analysis and decisions
   - Implementation details
   - Testing strategy

2. **[ARCHITECTURE_DIAGRAM.md](ARCHITECTURE_DIAGRAM.md)** (329 lines)
   - Visual architecture diagrams (Mermaid)
   - Current vs proposed architecture
   - Implementation flow
   - File structure

3. **[IMPLEMENTATION_CHECKLIST.md](IMPLEMENTATION_CHECKLIST.md)** (329 lines)
   - Step-by-step implementation guide
   - Checkboxes for tracking progress
   - Verification steps
   - Post-implementation tasks

## Recommendations

### Immediate Next Steps

1. **Review Documentation**: Review the three planning documents
2. **Ask Questions**: Clarify any concerns or questions
3. **Approve Plan**: Confirm the approach is acceptable
4. **Begin Implementation**: Switch to Code mode to implement

### Future Enhancements (Post-Initial Release)

Consider these features based on user feedback:
- Bulk read function: `variables([]string) map[string]string`
- Default value support: `variable_with_default(name, default string) string`
- Existence check: `variable_exists(name string) bool`
- JSON parsing: `json_variable(name string) dynamic`

## Questions for Review

1. **Naming**: Are you satisfied with `variable` and `sensitive_variable` as function names?
2. **Scope**: Should we implement both functions initially, or start with just `variable`?
3. **Documentation**: Is the level of documentation appropriate?
4. **Timeline**: Does the 5-8 hour estimate seem reasonable?
5. **Testing**: Are there any specific test cases you'd like to see?

## Conclusion

This implementation plan provides a comprehensive, well-tested approach to adding provider-defined functions while maintaining full backward compatibility. The design follows Terraform best practices and the existing codebase patterns.

The implementation is straightforward, low-risk, and provides significant value to users on Terraform 1.8+, while ensuring users on older versions can continue using the existing data sources without any changes.

---

**Ready to proceed?** Review the detailed documents and let me know if you have any questions or would like any changes to the plan.