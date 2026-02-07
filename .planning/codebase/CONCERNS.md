# Technical Concerns

## Overview

Known issues, technical debt, and areas requiring attention in the codebase.

## Known Issues

### Type Mismatch in check_testdata.go
**File**: `check_testdata.go`
**Issue**: Type mismatches between `*sysml.Model` and `sysml.Model`

```
Line 47: cannot use result.Model (variable of type *sysml.Model) as sysml.Model value
Line 76: cannot use model (variable of struct type sysml.Model) as *sysml.Model value
```

**Impact**: Check testdata tool doesn't compile
**Status**: Pre-existing issue, not critical for library functionality

### Grammar Coverage Gaps
**Current**: ~73% (58 of ~80 elements)
**Missing**: P3 low-priority elements
**Impact**: Some complex SysML files may not parse completely

### Validation Failures
**Current**: 96.4% success rate (54/56 files)
**Failing**: 2 files with complex syntax edge cases
**Impact**: Minor - not core functionality

## Technical Debt

### Planned but Unexecuted Plans
- Phase 4 plans 04-02 and 04-03 were planned but superseded by gap closure approach
- No functional impact, but roadmap history shows incomplete plans

### Hardcoded Paths
- Some workflow files reference hardcoded paths
- Should be parameterized for flexibility

### Documentation Duplication
- Some overlap between docs/README.md and docs/USAGE.md
- Could be consolidated in future cleanup

## Performance Considerations

### Memory Usage
- Parse trees retained by default (can use `WithDiscardTree()`)
- Large models may require streaming mode

### Parallel Parsing
- Speedup plateaus at ~2 workers (1.7x)
- Diminishing returns beyond 2 workers

## Security

### No Known Security Issues
- No external API calls
- No database connections
- File system access limited to parsing input files

### Potential Concerns
- Parser processes arbitrary SysML input
- No input size limits enforced
- Could be vulnerable to resource exhaustion with malicious input

## Maintainability

### ANTLR-Generated Code
- Files in `internal/parser/` should not be modified
- Generated code has `go vet` warnings (accepted as artifacts)
- Updating ANTLR version requires regenerating parser

### Test Coverage
- Some edge cases may lack coverage
- Integration tests depend on external validation files

## Recommendations

### High Priority
1. Fix `check_testdata.go` type mismatches
2. Increase grammar coverage to 85%+
3. Add input size limits for security

### Medium Priority
1. Consolidate documentation
2. Parameterize hardcoded paths
3. Add more edge case tests

### Low Priority
1. Optimize parallel parsing further
2. Reduce memory allocations in hot paths
3. Add more benchmark coverage
