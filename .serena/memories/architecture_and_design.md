# Architecture and Design

## Core Components

### Main Library (`ixdtf.go`)
- **Core Types**:
  - `IXDTFExtensions` struct - Handles extension data parsing and validation
  - `ParseError` struct - Detailed error reporting with position information  
  - `TimezoneConsistencyResult` struct - Timezone validation results

- **Key Functions**:
  - `Parse()` - Main parsing function for IXDTF strings
  - `Format()` / `FormatNano()` - Time formatting with IXDTF extensions
  - `Validate()` - Input validation without full parsing
  - `NewIXDTFExtensions()` - Extension builder function

- **Internal Helpers**:
  - RFC 3339 portion parsing
  - Suffix element parsing and validation
  - Timezone consistency checking
  - Extension tag handling

### ABNF Validation (`abnf.go`)
- **Abnf struct** - Contains compiled regex patterns for ABNF validation
- **Validation patterns** for:
  - Timezone names
  - Suffix keys and values
  - Extended date-time format components

## Error Handling Strategy
- **Structured errors** with specific error types:
  - `ErrCriticalExtension`, `ErrExperimentalExtension`, `ErrPrivateExtension`
  - `ErrInvalidExtension`, `ErrInvalidSuffix`, `ErrInvalidTimezone`
  - `ErrTimezoneOffsetMismatch`
- **ParseError** provides detailed context including position information

## Design Patterns
- **Zero external dependencies** - Pure Go standard library implementation
- **Backward compatibility** - Full RFC 3339 support
- **Extensibility** - Support for custom extensions via suffix elements
- **Validation-first approach** - ABNF patterns ensure format compliance
- **Immutable data structures** - Time-based operations maintain immutability

## Key Architectural Decisions
1. **Two-phase parsing**: RFC 3339 portion first, then suffix elements
2. **ABNF-based validation**: Regex patterns derived from RFC specifications
3. **Extension categorization**: Critical, experimental, and private extension handling
4. **Timezone consistency**: Cross-validation between offset and timezone name