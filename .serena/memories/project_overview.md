# Project Overview

## Purpose
IXDTF is a Go implementation of RFC 9557 Internet Extended Date/Time Format (IXDTF).

IXDTF extends RFC 3339 by adding optional suffix elements for timezone names and additional metadata while maintaining full backward compatibility.

## Tech Stack
- **Language**: Go 1.25.1
- **Testing**: Built-in Go testing framework
- **Linting**: golangci-lint with extensive configuration
- **CI/CD**: GitHub Actions with codecov integration
- **Security**: gosec security scanner

## Key Features
- RFC 3339 backward compatibility
- Extended date/time format parsing and formatting
- Timezone name support
- ABNF validation patterns
- Comprehensive error handling with detailed error types

## Repository Structure
- `ixdtf.go` - Main library implementation with parsing, formatting, and validation functions
- `abnf.go` - ABNF (Augmented Backus-Naur Form) patterns and validation
- `ixdtf_test.go` - Comprehensive test suite
- `abnf_test.go` - ABNF validation tests
- `.golangci.yml` - Extensive linting configuration (based on golden config)
- CI/CD configured with GitHub Actions for testing, linting, and security scanning

## Dependencies
- No external dependencies (pure Go standard library)

## License
MIT License (Copyright 2025 sadao komaki)