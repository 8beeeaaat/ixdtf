# Suggested Commands

## Development Commands

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with race detection
go test -race -v ./...

# Run tests with coverage
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Code Quality
```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run golangci-lint (requires golangci-lint installation)
golangci-lint run

# Run gosec security scanner
gosec ./...
```

### Dependencies
```bash
# Download dependencies
go mod download

# Verify dependencies
go mod verify

# Tidy dependencies
go mod tidy
```

### Build
```bash
# Build the package
go build ./...

# Install the package
go install ./...
```

## CI/CD Commands (as used in GitHub Actions)
- All tests are run automatically on push/PR
- golangci-lint validation
- Security scanning with gosec
- Coverage reporting to codecov

## Go Version
- Requires Go 1.25.1 (as specified in go.mod and .go-version)

## macOS Specific Notes
- Commands should work identically on macOS as on Linux
- Use standard Unix tools: `grep`, `find`, `ls`, `git`, etc.
- No special macOS-specific requirements