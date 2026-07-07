package ixdtf

import "time"

// IXDTFExtensions holds IXDTF suffix information that extends RFC 3339.
//
//nolint:revive // Keeping existing public API name for compatibility
type IXDTFExtensions struct {
	Location *time.Location

	// CriticalLocation reports whether the time-zone annotation carried a
	// critical "!" flag (e.g. "[!Europe/London]"). Per RFC 9557 Section 3.3 a
	// critical annotation MUST be processable, and per Section 3.4 an
	// application MUST act on any inconsistency it introduces.
	CriticalLocation bool

	// Tags contains extension tags as key-value pairs.
	// Example: map[ExtensionUnicodeCalendar]"japanese".
	Tags map[string]string

	// Critical indicates which tags are marked as critical (must be processed).
	// Critical tags are marked with "!" prefix in the IXDTF string.
	Critical map[string]bool
}

// NewIXDTFExtensionsArgs contains the arguments for creating IXDTFExtensions.
type NewIXDTFExtensionsArgs struct {
	Location *time.Location
	// CriticalLocation marks the time-zone annotation as critical ("!");
	// see IXDTFExtensions.CriticalLocation.
	CriticalLocation bool
	Tags             map[string]string
	Critical         map[string]bool
}

// NewIXDTFExtensions creates a new IXDTFExtensions with initialized maps.
func NewIXDTFExtensions(args *NewIXDTFExtensionsArgs) *IXDTFExtensions {
	if args == nil {
		args = &NewIXDTFExtensionsArgs{}
	}
	ext := &IXDTFExtensions{
		Location:         args.Location,
		CriticalLocation: args.CriticalLocation,
		Tags:             args.Tags,
		Critical:         args.Critical,
	}
	if ext.Tags == nil {
		ext.Tags = make(map[string]string)
	}
	if ext.Critical == nil {
		ext.Critical = make(map[string]bool)
	}
	return ext
}
