# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.1] - 2025-09-15

### Changed

- Improved README.md with more comprehensive examples and clearer documentation
- Updated README title format to align with RFC 9557 specification
- Enhanced usage examples with real-world timezone scenarios and error handling demonstrations
- Added detailed strict vs non-strict mode examples showing timezone offset validation behavior

### Technical

- Enhanced end-to-end test coverage with comprehensive round-trip testing scenarios
- Added timezone conversion testing with America/New_York location
- Implemented strict mode validation tests for timezone offset consistency
- Updated CI workflow to exclude main branch from push triggers (PR #6, #7)
- Improved test complexity with realistic timezone offset validation scenarios

### Fixed

- Better error handling examples in documentation showing timezone offset mismatch scenarios

## [0.1.0] - 2025-09-14

### Added

- Initial implementation of RFC 9557 Internet Extended Date/Time Format (IXDTF)
- Full RFC 3339 backward compatibility
- Support for timezone names and Unicode extension tags
- ABNF-based validation with structured error reporting
- Zero external dependencies implementation using pure Go standard library
- Comprehensive test suite with race condition detection
- CI/CD pipeline with multiple linters and security scanning
