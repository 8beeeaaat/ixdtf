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
 strictMode := false
 // Parse an IXDTF string
 t, extensions, err := ixdtf.Parse("2023-08-07T14:30:00Z[America/New_York][u-ca=gregorian]", strictMode)
 if err != nil {
  panic(err)
 }

 fmt.Printf("Time: %v\n", t)                 // Time: 2023-08-07 14:30:00 +0000 UTC
 fmt.Printf("Extensions: %+v\n", extensions) // &ixdtf.IXDTFExtensions{Location:America/New_York Tags:map[u-ca:gregorian] Critical:map[]}

 // Format a time with extensions
 now := time.Now()
 formatted, err := ixdtf.Format(now, extensions)
 if err != nil {
  panic(err)
 }
 fmt.Printf("Formatted: %s\n", formatted) // 2025-09-13T23:26:20+09:00[America/New_York][u-ca=gregorian]

 t, extensions, err = ixdtf.Parse(formatted, strictMode)
 if err != nil {
  panic(err)
 }

 fmt.Printf("Time: %v\n", t)                 // Time: 2023-08-07 14:30:00 +0000 UTC
 fmt.Printf("Extensions: %+v\n", extensions) // &ixdtf.IXDTFExtensions{Location:America/New_York Tags:map[u-ca:gregorian] Critical:map[]}

}
```

### Validation Only

```go
strictMode := true
err := ixdtf.Validate("2023-08-07T14:30:00Z[America/New_York]", strictMode)
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
- `2023-08-07T14:30:00Z[!u-ca=gregorian]` - Critical extension (must be processed)

## Extension Tag Validation

IXDTF supports extension tags with specific validation rules following RFC 9557:

### Validation Flow

Extension tags undergo multi-layer validation:

1. **ABNF Syntax Validation**: Keys and values must conform to RFC-defined patterns
2. **Extension Type Validation**:
   - **Private extensions** (`x-*`, `X-*`): Rejected per [BCP 178](https://www.rfc-editor.org/info/bcp178)
   - **Experimental extensions** (`_*`): Rejected unless specifically configured
   - **Standard extensions** (`u-*`, `t-*`): Subject to format validation
3. **Critical Extension Processing**: Extensions marked with `!` must be processable or rejected

ref: <https://www.rfc-editor.org/rfc/rfc9557.html#section-3.2>

#### Rejection Cases

The following extension patterns are automatically rejected:

```go
// Private extensions
"2023-08-07T14:30:00Z[x-custom=value]"       // Error: private extension cannot be processed

// Experimental extensions
"2023-08-07T14:30:00Z[_experimental=value]"  // Error: experimental extension cannot be processed

// Critical private/experimental extensions also
"2023-08-07T14:30:00Z[!x-custom=value]"      // Error: private extension cannot be processed
"2023-08-07T14:30:00Z[!_experimental=value]" // Error: experimental extension cannot be processed
```

### Supported Extensions

- **Unicode Extensions** (`u-*`): Calendar, locale, and formatting preferences
  - Example: `u-ca=gregorian` (Gregorian calendar)
- **Transform Extensions** (`t-*`): Text transformation specifications
- **Custom Extensions**: Application-specific extensions following RFC patterns

## API Reference

### Core Functions

- `Parse(s string, strict bool) (time.Time, *IXDTFExtensions, error)` - Parse IXDTF string
- `Format(t time.Time, ext *IXDTFExtensions) (string, error)` - Format time with extensions
- `FormatNano(t time.Time, ext *IXDTFExtensions) (string, error)` - Format with nanosecond precision
- `Validate(s string, strict bool) error` - Validate format without parsing

#### Strict flag

The second argument `strict` in `Parse` / `Validate` controls how strictly the library enforces consistency between:

1) the UTC offset embedded in the RFC 3339 portion, and
2) the IANA time zone name supplied in a suffix.

| Mode | Behavior | Example |
|------|----------|---------|
| `true` | If the zone-derived offset for that instant differs from the RFC 3339 numeric offset, an **error (ErrTimezoneOffsetMismatch)** is returned. | `2025-01-01T12:00:00+09:00[America/New_York]` → New York at that instant is `-05:00`, so mismatch → error |
| `false` | Mismatches do NOT produce an error. The original timestamp (its instant + numeric offset) is kept; the location is only applied if offsets match. | Same example above: no error; the provided time value is kept as-is (location not applied) |

Notes:

- An invalid or unresolvable time zone name always yields an error (regardless of mode).
- Time zones with `Etc/GMT±X` naming are skipped for consistency checking (POSIX inverted offset semantics would cause false positives).
- `Validate` follows the same policy: with `strict=false` an offset mismatch is considered acceptable.
- Extension tag syntax and critical tag handling are independent of `strict`.
- Recommended usage: accept loosely formed inputs with `strict=false` at system boundaries (ingest phase), then re-normalize if needed; enforce `strict=true` where data integrity or audit requirements apply.

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

- Go 1.24 or later

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
