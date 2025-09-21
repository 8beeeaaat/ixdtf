package ixdtf_test

import (
	"sort"
	"testing"
	"time"

	"github.com/8beeeaaat/ixdtf"
)

// Benchmark targets focus on core hot paths:
// 1. Format / FormatNano with different extension complexities
// 2. Parse / Validate for typical, extended, and error cases
// Goal: Provide ns/op & allocs metrics for public API.

func benchmarkTime() time.Time { return time.Date(2025, 1, 2, 3, 4, 5, 123456789, time.UTC) }

func BenchmarkFormat(b *testing.B) {
	base := benchmarkTime()
	benchTokyo := time.FixedZone("Asia/Tokyo", 9*3600)
	benchParis := time.FixedZone("Europe/Paris", 1*3600)
	cases := []struct {
		name string
		t    time.Time
		ext  *ixdtf.IXDTFExtensions
	}{
		{"noext", base, ixdtf.NewIXDTFExtensions(nil)},
		{"tz", base.In(benchTokyo), ixdtf.NewIXDTFExtensions(nil)},
		{"tz_specified", base, ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{Location: benchParis})},
		{
			"tags",
			base,
			ixdtf.NewIXDTFExtensions(
				&ixdtf.NewIXDTFExtensionsArgs{Tags: map[string]string{"u-ca": "gregory", "a": "1", "b": "2"}},
			),
		},
	}

	sort.Slice(cases, func(i, j int) bool { return cases[i].name < cases[j].name })

	b.ReportAllocs()
	for _, c := range cases {
		b.Run(c.name, func(sb *testing.B) {
			sb.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _ = ixdtf.Format(c.t, c.ext)
				}
			})
		})
	}
}

func BenchmarkFormatNano(b *testing.B) {
	base := benchmarkTime()
	benchTokyo := time.FixedZone("Asia/Tokyo", 9*3600)
	cases := []struct {
		name string
		t    time.Time
		ext  *ixdtf.IXDTFExtensions
	}{
		{"noext", base, ixdtf.NewIXDTFExtensions(nil)},
		{
			"tz_tags",
			base.In(benchTokyo),
			ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{Tags: map[string]string{"u-ca": "gregory"}}),
		},
	}

	sort.Slice(cases, func(i, j int) bool { return cases[i].name < cases[j].name })

	b.ReportAllocs()
	for _, c := range cases {
		b.Run(c.name, func(sb *testing.B) {
			sb.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _ = ixdtf.FormatNano(c.t, c.ext)
				}
			})
		})
	}
}

func BenchmarkParse(b *testing.B) {
	cases := []struct {
		name, input string
		strict      bool
	}{
		{"rfc3339", "2025-01-02T03:04:05Z", false},
		{"rfc3339_nano", "2025-01-02T03:04:05.123456789Z", false},
		{"extended_tz", "2025-06-07T08:09:10+09:00[Asia/Tokyo]", false},
		{"extended_tz_tags", "2025-06-07T08:09:10+01:00[Europe/Paris][u-ca=gregory]", false},
		{"mismatch_non_strict", "2025-06-01T12:00:00+09:00[America/New_York]", false},
		{"mismatch_strict", "2025-06-01T12:00:00+09:00[America/New_York]", true},
		{"invalid_suffix", "2025-01-01T00:00:00Z[unclosed", false},
	}

	sort.Slice(cases, func(i, j int) bool { return cases[i].name < cases[j].name })

	b.ReportAllocs()
	for _, c := range cases {
		b.Run(c.name, func(sb *testing.B) {
			sb.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _, _ = ixdtf.Parse(c.input, c.strict)
				}
			})
		})
	}
}

func BenchmarkValidate(b *testing.B) {
	cases := []struct {
		name, input string
		strict      bool
	}{
		{"rfc3339", "2025-01-02T03:04:05Z", false},
		{"extended", "2025-06-07T08:09:10+09:00[Asia/Tokyo]", false},
		{"extended_tags", "2025-06-07T08:09:10Z[u-ca=gregory]", false},
		{"mismatch_strict", "2025-06-01T12:00:00+09:00[America/New_York]", true},
		{"invalid_ext", "2025-01-01T00:00:00Z[INVALID=val]", false},
	}

	sort.Slice(cases, func(i, j int) bool { return cases[i].name < cases[j].name })

	b.ReportAllocs()
	for _, c := range cases {
		b.Run(c.name, func(sb *testing.B) {
			sb.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_ = ixdtf.Validate(c.input, c.strict)
				}
			})
		})
	}
}
