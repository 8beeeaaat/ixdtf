// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ixdtf_test

import (
	"fmt"
	"time"

	"github.com/8beeeaaat/ixdtf"
)

func ExampleParse() {
	// Parse basic RFC 3339 without extensions
	t1, ext1, _ := ixdtf.Parse("2006-01-02T15:04:05Z")
	fmt.Printf("Basic: %v, timezone=%s, tags=%d\n", t1, ext1.TimeZone, len(ext1.Tags))

	// Parse with timezone suffix
	t2, ext2, _ := ixdtf.Parse("2006-01-02T15:04:05Z[UTC]")
	fmt.Printf("With timezone: %v, timezone=%s\n", t2, ext2.TimeZone)

	// Parse with extension tag
	t3, ext3, _ := ixdtf.Parse("2006-01-02T15:04:05Z[u-ca=japanese]")
	fmt.Printf("With extension: %v, u-ca=%s\n", t3, ext3.Tags["u-ca"])

	// Parse with critical extension
	t4, ext4, _ := ixdtf.Parse("2006-01-02T15:04:05Z[!u-ca=japanese]")
	fmt.Printf("Critical extension: %v, u-ca=%s, critical=%t\n", t4, ext4.Tags["u-ca"], ext4.Critical["u-ca"])

	// Output:
	// Basic: 2006-01-02 15:04:05 +0000 UTC, timezone=, tags=0
	// With timezone: 2006-01-02 15:04:05 +0000 UTC, timezone=UTC
	// With extension: 2006-01-02 15:04:05 +0000 UTC, u-ca=japanese
	// Critical extension: 2006-01-02 15:04:05 +0000 UTC, u-ca=japanese, critical=true
}

func ExampleFormat() {
	t := time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)

	// Format basic time
	ext1 := ixdtf.NewIXDTFExtensions()
	fmt.Println(ixdtf.Format(t, ext1))

	// Format with timezone
	ext2 := ixdtf.NewIXDTFExtensions()
	ext2.TimeZone = "UTC"
	fmt.Println(ixdtf.Format(t, ext2))

	// Format with extension
	ext3 := ixdtf.NewIXDTFExtensions()
	ext3.Tags["u-ca"] = "japanese"
	fmt.Println(ixdtf.Format(t, ext3))

	// Format with critical extension
	ext4 := ixdtf.NewIXDTFExtensions()
	ext4.Tags["u-ca"] = "japanese"
	ext4.Critical["u-ca"] = true
	fmt.Println(ixdtf.Format(t, ext4))

	// Output:
	// 2006-01-02T15:04:05Z
	// 2006-01-02T15:04:05Z[UTC]
	// 2006-01-02T15:04:05Z[u-ca=japanese]
	// 2006-01-02T15:04:05Z[!u-ca=japanese]
}
