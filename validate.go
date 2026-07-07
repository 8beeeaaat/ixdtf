package ixdtf

import (
	"time"

	"github.com/8beeeaaat/ixdtf/abnf"
)

// validateExtensionsStrict validates IXDTF extensions for correctness and
// processes critical extensions (RFC 9557 Section 3.3). In strict mode,
// registered tag values are also validated.
func validateExtensionsStrict(ext *IXDTFExtensions, strict bool) error {
	if ext == nil {
		return nil
	}

	if err := validateLocationStrict(ext.Location, strict); err != nil {
		return err
	}

	if err := validateTagKeys(ext.Tags); err != nil {
		return err
	}

	if err := validateCriticalTags(ext.Tags, ext.Critical); err != nil {
		return err
	}

	if strict {
		if err := validateTagValuesStrict(ext.Tags); err != nil {
			return err
		}
	}

	return nil
}

func validateLocationStrict(location *time.Location, strict bool) error {
	if location == nil {
		return nil
	}
	// An offset-derived FixedZone (e.g. from "[+09:00]") resolves to itself;
	// unknown named zones are ignored in non-strict mode per RFC 9557.
	if _, err := resolveLocation(location); err != nil && strict {
		return ErrInvalidTimezone
	}
	return nil
}

func validateTagKeys(tags map[string]string) error {
	// Basic tag key validation (syntactic). Value validation is already handled when creating tags.
	for key := range tags {
		if err := abnf.AbnfSuffixKey.ValidateSuffixKey(key); err != nil {
			return err
		}
	}
	return nil
}

func validateCriticalTags(tags map[string]string, critical map[string]bool) error {
	// Critical tag processing:
	// * If a key is marked critical but missing in Tags -> error.
	// * If present but value is empty -> error.
	for key, isCritical := range critical {
		if !isCritical {
			continue
		}
		value, exists := tags[key]
		if !exists { // missing critical tag
			return ErrCriticalExtension
		}
		if err := validateCriticalExtension(key, value); err != nil {
			return err
		}
	}
	return nil
}

// validateCriticalExtension enforces critical extension processing rules.
func validateCriticalExtension(key, value string) error {
	if value == "" { // empty value not allowed for critical
		return ErrCriticalExtension
	}
	return validateTagValue(key, value)
}

func validateTagValuesStrict(tags map[string]string) error {
	for key, value := range tags {
		if err := validateTagValue(key, value); err != nil {
			return err
		}
	}
	return nil
}
