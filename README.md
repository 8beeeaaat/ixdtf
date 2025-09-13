# IXDTF - RFC 9557 Internet Extended Date/Time Format (IXDTF) Support for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/8beeeaaat/ixdtf.svg)](https://pkg.go.dev/github.com/8beeeaaat/ixdtf)
[![codecov](https://codecov.io/gh/8beeeaaat/ixdtf/graph/badge.svg?token=B7ECLHCZS1)](https://codecov.io/gh/8beeeaaat/ixdtf)
[![CI](https://github.com/8beeeaaat/ixdtf/workflows/CI/badge.svg)](https://github.com/8beeeaaat/ixdtf/actions)

A Go implementation of [RFC 9557 Internet Extended Date/Time Format (IXDTF)](https://datatracker.ietf.org/doc/rfc9557/).

IXDTF extends RFC 3339 by adding optional suffix elements for timezone names and additional metadata while maintaining full backward compatibility.

## Features

- **RFC 3339 Compatible**: Full backward compatibility with existing RFC 3339 date/time strings
- **Extended Format Support**: Handles timezone names and additional metadata via suffix elements
- **Zero Dependencies**: Pure Go implementation using only the standard library
- **Comprehensive Validation**: ABNF-based validation ensuring format compliance
- **Detailed Error Reporting**: Structured error types with position information
- **High Performance**: Optimized parsing and formatting operations

## Installation

```bash
go get github.com/8beeeaaat/ixdtf
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "time"

    "github.com/8beeeaaat/ixdtf"
)

func main() {
    // Parse an IXDTF string
    t, extensions, err := ixdtf.Parse("2023-08-07T14:30:00Z[America/New_York][u-ca=gregorian]")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Time: %v\n", t)
    fmt.Printf("Extensions: %+v\n", extensions)

    // Format a time with extensions
    now := time.Now()
    formatted := ixdtf.Format(now, extensions)
    fmt.Printf("Formatted: %s\n", formatted)
}
```

### Validation Only

```go
// Validate without full parsing
err := ixdtf.Validate("2023-08-07T14:30:00Z[America/New_York]")
if err != nil {
    fmt.Printf("Invalid format: %v\n", err)
}
```

## IXDTF Format

IXDTF extends RFC 3339 with optional suffix elements:

```text
<RFC3339-date-time>[<timezone-name>][<extension>...]
```

### Examples

- `2023-08-07T14:30:00Z` - Standard RFC 3339
- `2023-08-07T14:30:00Z[America/New_York]` - With timezone name
- `2023-08-07T14:30:00Z[u-ca=gregorian]` - With Unicode calendar extension
- `2023-08-07T14:30:00Z[America/New_York][u-ca=gregorian]` - Multiple suffixes

## API Reference

### Core Functions

- `Parse(s string) (time.Time, *IXDTFExtensions, error)` - Parse IXDTF string
- `Format(t time.Time, ext *IXDTFExtensions) string` - Format time with extensions
- `FormatNano(t time.Time, ext *IXDTFExtensions) string` - Format with nanosecond precision
- `Validate(s string) error` - Validate format without parsing

### Types

- `IXDTFExtensions` - Container for extension data
- `ParseError` - Detailed parsing error with position information
- `TimezoneConsistencyResult` - Timezone validation results

## Error Handling

The library provides structured error types for different failure scenarios:

```go
t, ext, err := ixdtf.Parse("invalid-date")
if err != nil {
    if parseErr, ok := err.(*ixdtf.ParseError); ok {
        fmt.Printf("Parse error at position %d: %s\n", parseErr.Position, parseErr.Message)
    }
}
```

## Development

### Prerequisites

- Go 1.25.1 or later

### Building and Testing

```bash
# Run tests
go test ./...

# Run tests with race detection
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# Lint code (requires golangci-lint)
golangci-lint run

# Security scan (requires gosec)
gosec ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes and add tests
4. Ensure all tests pass and code is properly formatted
5. Commit your changes (`git commit -am 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Code Quality

This project maintains high code quality standards:

- All tests must pass
- Code coverage should not decrease
- golangci-lint must pass with zero warnings
- gosec security scan must pass
- Code must be properly formatted with `go fmt`

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## References

- [RFC 9557 - Internet Extended Date/Time Format (IXDTF)](https://datatracker.ietf.org/doc/rfc9557/)
- [RFC 3339 - Date and Time on the Internet: Timestamps](https://datatracker.ietf.org/doc/rfc3339/)
