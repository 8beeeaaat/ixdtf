# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2026-01-18

### Added

- Added validation and support for Unicode calendar identifiers in `u-ca` tags using known CLDR identifiers

### Changed

- Enforced critical `u-ca` tag validation with clearer errors for invalid calendar identifiers
- Updated README examples and parse error usage for clarity (PR #12)

### Technical

- Bumped GitHub Actions dependencies: actions/checkout v6, actions/cache v5, golangci-lint-action v9 (PR #13, #14, #15)

## [0.2.0] - 2025-09-21

### Changed

- Refactored timezone handling with improved error messages and consistency (PR #10)
- Enhanced parsing logic to support numeric timezone offsets and improved validation for timezone names
- Modified suffix parsing to handle multiple tags and critical flags more effectively
- Improved error handling in strict and non-strict modes for better RFC 9557 compliance

### Technical

- Refactored E2E tests for timezone validation with improved variable naming and error handling (PR #9)
- Enhanced test coverage with new test cases covering unknown time zones and timezone inconsistencies
- Improved benchmark testing with more comprehensive scenarios
- Added comprehensive internal testing for timezone validation logic
- Enhanced ABNF validation patterns with better regex-based format checking

### Fixed

- Better timezone offset validation and error reporting
- Improved handling of timezone name parsing edge cases
- Enhanced error message clarity for timezone-related validation failures

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
