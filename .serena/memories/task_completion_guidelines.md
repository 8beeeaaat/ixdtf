# Task Completion Guidelines

## Required Steps After Code Changes

### 1. Code Quality Checks
Always run the following commands after making code changes:

```bash
# Format code
go fmt ./...

# Vet for common issues
go vet ./...

# Lint with golangci-lint (if available)
golangci-lint run
```

### 2. Testing
Comprehensive testing is required:

```bash
# Run all tests
go test ./...

# Run with race detection
go test -race ./...

# Check coverage
go test -coverprofile=coverage.out
```

### 3. Security Check
Run security analysis:

```bash
# Security scanning (if gosec is available)
gosec ./...
```

### 4. Dependency Management
Keep dependencies clean:

```bash
# Tidy dependencies
go mod tidy

# Verify integrity
go mod verify
```

## Definition of "Complete"
A task is only complete when:
- ✅ All tests pass (`go test ./...`)
- ✅ No linting errors (`golangci-lint run`)
- ✅ No vet issues (`go vet ./...`)
- ✅ Code is properly formatted (`go fmt ./...`)
- ✅ No security issues (`gosec ./...`)
- ✅ Dependencies are tidy (`go mod tidy`)

## CI/CD Compliance
All changes must pass the same checks that run in CI:
- Test suite with race detection
- golangci-lint validation
- Security scanning
- Format verification
- Coverage reporting

## Before Committing
1. Run the full test suite
2. Verify linting passes
3. Check that coverage hasn't decreased significantly
4. Ensure no security warnings
5. Review changes for adherence to project conventions