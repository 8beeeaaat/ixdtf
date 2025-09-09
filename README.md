# IXDTF - Internet Extended Date/Time Format for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/8beeeaaat/ixdtf.svg)](https://pkg.go.dev/github.com/8beeeaaat/ixdtf)
[![Go Report Card](https://goreportcard.com/badge/github.com/8beeeaaat/ixdtf)](https://goreportcard.com/report/github.com/8beeeaaat/ixdtf)
[![Coverage Status](https://coveralls.io/repos/github/8beeeaaat/ixdtf/badge.svg)](https://coveralls.io/github/8beeeaaat/ixdtf)

A Go implementation of [RFC 9557 Internet Extended Date/Time Format (IXDTF)](https://datatracker.ietf.org/doc/rfc9557/).

IXDTF extends RFC 3339 by adding optional suffix elements for timezone names and additional metadata while maintaining full backward compatibility.

## Features

- **RFC 3339 Compatible**: Fully backward compatible with existing RFC 3339 implementations
- **Timezone Suffixes**: Support for IANA timezone names like `[Asia/Tokyo]`
- **Extension Tags**: Support for extension metadata like `[u-ca=japanese]`
- **Critical Flags**: Support for critical extensions that must be processed `[!u-ca=japanese]`
- **High Performance**: Built on Go's standard `time` package for optimal performance
- **Comprehensive Testing**: 98%+ test coverage with extensive edge case testing

## Installation

```bash
go get github.com/8beeeaaat/ixdtf
```

## Quick Start

```go
package main

import (
    "fmt"
    "time"

    "github.com/8beeeaaat/ixdtf"
)

func main() {
    // Parse IXDTF with timezone suffix
    t, ext, err := ixdtf.ParseRFC9557("2006-01-02T15:04:05+09:00[Asia/Tokyo]")
    if err != nil {
        panic(err)
    }

    fmt.Printf("Time: %v\n", t)                    // 2006-01-02 15:04:05 +0900 JST
    fmt.Printf("Timezone: %s\n", ext.TimeZone)    // Asia/Tokyo

    // Format time with extensions
    ext.Tags = map[string]string{"u-ca": "japanese"}
    formatted := ixdtf.FormatRFC9557(t, ext)
    fmt.Printf("Formatted: %s\n", formatted)      // 2006-01-02T15:04:05+09:00[Asia/Tokyo][u-ca=japanese]
}
```

## Supported Formats

### Basic RFC 3339 (unchanged)

```
2006-01-02T15:04:05Z
2006-01-02T15:04:05.999999999-07:00
```

### IXDTF with Timezone Suffix

```
2006-01-02T15:04:05Z[UTC]
2006-01-02T15:04:05+09:00[Asia/Tokyo]
```

### IXDTF with Extension Tags

```
2006-01-02T15:04:05Z[u-ca=japanese]
2006-01-02T15:04:05Z[!u-ca=japanese]  // Critical extension
```

### Complex IXDTF

```
2006-01-02T15:04:05+09:00[Asia/Tokyo][u-ca=japanese][x-example=value]
```

## API Reference

### Core Functions

#### `ParseRFC9557(s string) (time.Time, IXDTFExtensions, error)`

Parses an IXDTF string and returns the parsed time and extensions.

```go
t, ext, err := ixdtf.ParseRFC9557("2006-01-02T15:04:05Z[UTC]")
```

#### `FormatRFC9557(t time.Time, ext IXDTFExtensions) string`

Formats a time with IXDTF extensions using RFC 3339 format.

```go
formatted := ixdtf.FormatRFC9557(time.Now(), ext)
```

#### `FormatRFC9557Nano(t time.Time, ext IXDTFExtensions) string`

Formats a time with IXDTF extensions using RFC 3339 format with nanoseconds.

```go
formatted := ixdtf.FormatRFC9557Nano(time.Now(), ext)
```

### Data Structures

#### `IXDTFExtensions`

```go
type IXDTFExtensions struct {
    TimeZone string            // IANA timezone name
    Tags     map[string]string // Extension tags
    Critical map[string]bool   // Critical flags for tags
}
```

#### `NewIXDTFExtensions() IXDTFExtensions`

Creates a new IXDTFExtensions with initialized maps.

```go
ext := ixdtf.NewIXDTFExtensions()
ext.TimeZone = "Asia/Tokyo"
ext.Tags["u-ca"] = "japanese"
ext.Critical["u-ca"] = true
```

## Examples

### Parsing Examples

```go
// Basic RFC 3339
t, ext, _ := ixdtf.ParseRFC9557("2006-01-02T15:04:05Z")

// With timezone
t, ext, _ := ixdtf.ParseRFC9557("2006-01-02T15:04:05Z[UTC]")

// With extension tag
t, ext, _ := ixdtf.ParseRFC9557("2006-01-02T15:04:05Z[u-ca=japanese]")

// With critical extension
t, ext, _ := ixdtf.ParseRFC9557("2006-01-02T15:04:05Z[!u-ca=japanese]")

// Complex example
t, ext, _ := ixdtf.ParseRFC9557("2006-01-02T15:04:05+09:00[Asia/Tokyo][u-ca=japanese]")
```

### Formatting Examples

```go
t := time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)

// Basic formatting
ext := ixdtf.NewIXDTFExtensions()
fmt.Println(ixdtf.FormatRFC9557(t, ext))  // 2006-01-02T15:04:05Z

// With timezone
ext.TimeZone = "UTC"
fmt.Println(ixdtf.FormatRFC9557(t, ext))  // 2006-01-02T15:04:05Z[UTC]

// With extension
ext.Tags["u-ca"] = "japanese"
fmt.Println(ixdtf.FormatRFC9557(t, ext))  // 2006-01-02T15:04:05Z[UTC][u-ca=japanese]

// With critical extension
ext.Critical["u-ca"] = true
fmt.Println(ixdtf.FormatRFC9557(t, ext))  // 2006-01-02T15:04:05Z[UTC][!u-ca=japanese]
```

## Error Handling

The package provides detailed error information for invalid formats:

```go
t, ext, err := ixdtf.ParseRFC9557("invalid-format")
if err != nil {
    fmt.Printf("Parse error: %v\n", err)
}
```

## Performance

This implementation is built on Go's standard `time` package and optimized for performance:

- Zero external dependencies
- Efficient string parsing and formatting
- Minimal memory allocations
- Reuses existing RFC 3339 parsing where possible

## Testing

Run the test suite:

```bash
go test -v
```

Run tests with coverage:

```bash
go test -cover
```

Generate coverage report:

```bash
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the BSD 3-Clause License - see the [LICENSE](LICENSE) file for details.

## References

- [RFC 9557: Date and Time on the Internet: Timestamps with Additional Information](https://datatracker.ietf.org/doc/rfc9557/)
- [RFC 3339: Date and Time on the Internet: Timestamps](https://datatracker.ietf.org/doc/rfc3339/)
- [Go time package](https://pkg.go.dev/time)

## Related Projects

- [GitHub Issue #75320](https://github.com/golang/go/issues/75320) - Original proposal for Go standard library integration

---

Made with ❤️ in Go
