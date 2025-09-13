# Code Style and Conventions

## General Go Conventions
- **Package naming**: Single lowercase word package name (`ixdtf`)
- **Function naming**: CamelCase for exported functions, camelCase for internal
- **Struct naming**: PascalCase for exported types
- **Error naming**: Errors prefixed with `Err` (e.g., `ErrInvalidExtension`)

## Specific Project Patterns

### Error Handling
- **Custom error types**: Define specific error variables for different failure modes
- **ParseError struct**: Includes position information and context for parsing failures
- **Error methods**: All custom errors implement the `Error() string` method

### Naming Conventions
- **Constants**: PascalCase with descriptive prefixes (e.g., `LayoutRFC3339Extended`)
- **Private functions**: camelCase starting with lowercase (e.g., `parseRFC3339Portion`)
- **Struct fields**: PascalCase for exported, camelCase for internal

### Documentation
- **Function comments**: Start with function name, describe purpose and parameters
- **Struct comments**: Describe the purpose and usage of the struct
- **Error comments**: Document when and why specific errors occur

## Code Quality Standards
- **golangci-lint configuration**: Extremely strict linting with 80+ enabled linters
- **Line length**: Maximum 120 characters (configured in golines)
- **Cyclomatic complexity**: Max 30 for functions, 10.0 package average
- **Function length**: Max 100 lines, 50 statements

## Testing Conventions
- **Test file naming**: `*_test.go`
- **Test function naming**: `TestFunctionName` pattern
- **Helper functions**: camelCase (e.g., `containsString`, `extensionsEqual`)
- **Test data**: Structured test cases with descriptive names

## Import Organization
- **Standard library**: First group
- **Third-party**: Second group (none in this project)
- **Local imports**: Third group
- **goimports**: Automatic formatting and grouping

## Performance Considerations
- **Zero allocations**: Where possible, avoid unnecessary memory allocations
- **Regex compilation**: Pre-compile regex patterns in ABNF struct
- **String building**: Use efficient string construction methods
- **Time parsing**: Leverage Go's built-in time parsing for RFC 3339 portion