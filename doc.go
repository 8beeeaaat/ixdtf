// Package ixdtf implements RFC 9557 Internet Extended Date/Time Format (IXDTF).
// IXDTF extends RFC 3339 by adding optional suffix elements for timezone names.
// and additional metadata while maintaining full backward compatibility.
//
// See RFC 9557: https://datatracker.ietf.org/doc/rfc9557/
//
// The package is organized so each file covers one RFC 9557 concern:
//
//   - format.go: serialization (Section 4.1) and critical output rules (Section 3.3)
//   - parse.go: RFC 3339 core, unknown local offset (Section 2.2), and orchestration
//   - suffix.go: suffix grammar (Section 4.1)
//   - timezone.go: time-zone resolution and consistency (Section 3.4)
//   - validate.go: extension semantics (Section 3.3)
//   - calendar.go: the calendar suffix key (Section 5)
//   - extensions.go: the suffix data model (Section 3)
//   - errors.go: error types and sentinels
package ixdtf

import "time"

// Layout represents time layout strings.
type Layout string

const (
	// LayoutRFC3339 represents the RFC 3339 layout.
	LayoutRFC3339 Layout = time.RFC3339

	// LayoutRFC3339Nano represents the RFC 3339 layout with nanoseconds.
	LayoutRFC3339Nano Layout = time.RFC3339Nano

	// LayoutRFC3339Extended represents the RFC 9557(IXDTF) layout with optional timezone and extensions.
	// https://www.rfc-editor.org/rfc/rfc9557.html#section-3
	LayoutRFC3339Extended Layout = time.RFC3339 + "*([time-zone-name][tags])"

	// LayoutRFC3339NanoExtended represents the RFC 9557(IXDTF) layout with nanoseconds and optional timezone and extensions.
	// https://www.rfc-editor.org/rfc/rfc9557.html#section-3
	LayoutRFC3339NanoExtended Layout = time.RFC3339Nano + "*([time-zone-name][tags])"
)
