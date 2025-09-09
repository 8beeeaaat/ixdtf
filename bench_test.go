// Copyright 2024 8beeeaaat. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ixdtf

import (
	"testing"
	"time"
)

var (
	benchmarkTimeUTC     = time.Date(2006, 1, 2, 15, 4, 5, 123456789, time.UTC)
	benchmarkTimeJST     = time.Date(2006, 1, 2, 15, 4, 5, 123456789, time.FixedZone("JST", 9*3600))
	benchmarkExtEmpty    = NewIXDTFExtensions()
	benchmarkExtTimezone = IXDTFExtensions{
		TimeZone: "Asia/Tokyo",
		Tags:     map[string]string{},
		Critical: map[string]bool{},
	}
	benchmarkExtFull = IXDTFExtensions{
		TimeZone: "Asia/Tokyo",
		Tags: map[string]string{
			"u-ca": "japanese",
			"u-nu": "latn",
			"x-tz": "jst",
		},
		Critical: map[string]bool{
			"u-ca": true,
		},
	}
)

// Benchmark parsing simple RFC 3339 datetime
func BenchmarkParse_RFC3339(b *testing.B) {
	input := "2006-01-02T15:04:05Z"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := Parse(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark parsing IXDTF with timezone suffix
func BenchmarkParse_WithTimezone(b *testing.B) {
	input := "2006-01-02T15:04:05+09:00[Asia/Tokyo]"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := Parse(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark parsing IXDTF with extensions
func BenchmarkParse_WithExtensions(b *testing.B) {
	input := "2006-01-02T15:04:05Z[Asia/Tokyo][!u-ca=japanese][u-nu=latn]"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := Parse(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark parsing IXDTF with nanoseconds
func BenchmarkParse_WithNanoseconds(b *testing.B) {
	input := "2006-01-02T15:04:05.123456789Z[Asia/Tokyo]"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := Parse(input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark formatting simple RFC 3339
func BenchmarkFormat_RFC3339(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Format(benchmarkTimeUTC, benchmarkExtEmpty)
	}
}

// Benchmark formatting with timezone
func BenchmarkFormat_WithTimezone(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Format(benchmarkTimeJST, benchmarkExtTimezone)
	}
}

// Benchmark formatting with full extensions
func BenchmarkFormat_WithExtensions(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Format(benchmarkTimeJST, benchmarkExtFull)
	}
}

// Benchmark formatting with nanoseconds
func BenchmarkFormatNano_RFC3339(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FormatNano(benchmarkTimeUTC, benchmarkExtEmpty)
	}
}

// Benchmark formatting nanoseconds with extensions
func BenchmarkFormatNano_WithExtensions(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = FormatNano(benchmarkTimeJST, benchmarkExtFull)
	}
}

// Benchmark round trip (parse then format)
func BenchmarkRoundTrip_RFC3339(b *testing.B) {
	input := "2006-01-02T15:04:05Z"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t, ext, err := Parse(input)
		if err != nil {
			b.Fatal(err)
		}
		_ = Format(t, ext)
	}
}

// Benchmark round trip with extensions
func BenchmarkRoundTrip_WithExtensions(b *testing.B) {
	input := "2006-01-02T15:04:05Z[Asia/Tokyo][!u-ca=japanese][u-nu=latn]"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t, ext, err := Parse(input)
		if err != nil {
			b.Fatal(err)
		}
		_ = Format(t, ext)
	}
}

// Benchmark suffix parsing specifically
func BenchmarkParseSuffix(b *testing.B) {
	input := "[Asia/Tokyo][!u-ca=japanese][u-nu=latn][x-custom=value]"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = parseSuffix(input)
	}
}

// Benchmark appendSuffix specifically
func BenchmarkAppendSuffix(b *testing.B) {
	base := []byte("2006-01-02T15:04:05Z")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b := make([]byte, len(base))
		copy(b, base)
		_ = appendSuffix(b, benchmarkExtFull)
	}
}

// Benchmark comparison with standard time.Parse
func BenchmarkStandardTime_Parse(b *testing.B) {
	input := "2006-01-02T15:04:05Z"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := time.Parse(time.RFC3339, input)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark comparison with standard time.Format
func BenchmarkStandardTime_Format(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = benchmarkTimeUTC.Format(time.RFC3339)
	}
}
