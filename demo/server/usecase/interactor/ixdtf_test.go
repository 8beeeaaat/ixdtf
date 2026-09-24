package interactor

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/8beeeaaat/ixdtf/demo/server/entity"
)

type sampleFile struct {
	Samples []sample `json:"samples"`
}

type sample struct {
	ID             string `json:"id"`
	Input          string `json:"input"`
	Strict         bool   `json:"strict"`
	ExpectServer   string `json:"expect_server"`
	ExpectErrorSub string `json:"expect_error_substr"`
}

// loadSamples は共有フィクスチャ (testdata/ixdtf_samples.json) を読み込む。
// パスはこのテストファイルからの相対。
func loadSamples(t *testing.T) sampleFile {
	t.Helper()
	p := filepath.Join("..", "..", "..", "testdata", "ixdtf_samples.json")
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read samples: %v", err)
	}
	var sf sampleFile
	if err := json.Unmarshal(data, &sf); err != nil {
		t.Fatalf("unmarshal samples: %v", err)
	}
	if len(sf.Samples) == 0 {
		t.Fatal("no samples loaded")
	}
	return sf
}

func TestParse_Samples(t *testing.T) {
	sf := loadSamples(t)
	i := NewIxdtf()
	for _, s := range sf.Samples {
		t.Run(s.ID, func(t *testing.T) {
			out := i.Parse(s.Input, s.Strict, false)
			switch s.ExpectServer {
			case "ok":
				if !out.OK {
					t.Fatalf("expected ok, got error: %v", out.Err)
				}
			case "error":
				if out.OK {
					t.Fatal("expected error, got ok")
				}
				if out.Err == nil {
					t.Fatal("expected non-nil Err")
				}
				// エラー文言は ixdtf ライブラリ原文のまま (A-3)。
				if s.ExpectErrorSub != "" && !strings.Contains(out.Err.Error(), s.ExpectErrorSub) {
					t.Fatalf("error %q does not contain %q", out.Err.Error(), s.ExpectErrorSub)
				}
			default:
				t.Fatalf("unknown expect_server %q", s.ExpectServer)
			}
		})
	}
}

func TestParse_Fields(t *testing.T) {
	out := NewIxdtf().Parse("2026-07-07T23:30:00+09:00[Asia/Tokyo]", false, false)
	if !out.OK {
		t.Fatalf("unexpected error: %v", out.Err)
	}
	if out.OffsetSeconds != 9*3600 {
		t.Errorf("offset = %d, want 32400", out.OffsetSeconds)
	}
	if out.Ext.TimeZone == nil || *out.Ext.TimeZone != "Asia/Tokyo" {
		t.Errorf("time_zone = %v, want Asia/Tokyo", out.Ext.TimeZone)
	}
	if len(out.Ext.Tags) != 0 {
		t.Errorf("tags = %+v, want empty", out.Ext.Tags)
	}
}

func TestParse_CriticalTag(t *testing.T) {
	out := NewIxdtf().Parse("2026-07-07T23:30:00+09:00[Asia/Tokyo][!u-ca=japanese]", false, false)
	if !out.OK {
		t.Fatalf("unexpected error: %v", out.Err)
	}
	if len(out.Ext.Tags) != 1 {
		t.Fatalf("tags len = %d, want 1", len(out.Ext.Tags))
	}
	tag := out.Ext.Tags[0]
	if tag.Key != "u-ca" || tag.Value != "japanese" || !tag.Critical {
		t.Errorf("tag = %+v", tag)
	}
	if out.Ext.TimeZoneCritical {
		t.Error("time_zone_critical should be false")
	}
}

func TestParse_CriticalTimezone(t *testing.T) {
	out := NewIxdtf().Parse("2026-07-07T23:30:00+09:00[!Asia/Tokyo]", false, false)
	if !out.OK {
		t.Fatalf("unexpected error: %v", out.Err)
	}
	if !out.Ext.TimeZoneCritical {
		t.Error("want time_zone_critical = true")
	}
	if out.Ext.TimeZone == nil || *out.Ext.TimeZone != "Asia/Tokyo" {
		t.Errorf("time_zone = %v", out.Ext.TimeZone)
	}
}

func TestParse_ValidateOnly(t *testing.T) {
	i := NewIxdtf()
	if out := i.Parse("2026-07-07T14:30:00Z", false, true); !out.OK {
		t.Fatalf("validate-only valid input: %v", out.Err)
	}
	out := i.Parse("not-a-timestamp", false, true)
	if out.OK {
		t.Fatal("expected error")
	}
	if !strings.Contains(out.Err.Error(), "cannot parse") {
		t.Errorf("error %q", out.Err.Error())
	}
}

func TestFormat(t *testing.T) {
	i := NewIxdtf()
	tz := "Asia/Tokyo"
	// 2026-07-07T14:30:00Z == 23:30 JST
	tm := time.Date(2026, 7, 7, 14, 30, 0, 0, time.UTC)
	ext := entity.Extensions{
		TimeZone: &tz,
		Tags:     []entity.ExtensionTag{{Key: "u-ca", Value: "japanese"}},
	}
	out := i.Format(tm.UnixNano(), ext)
	if !out.OK {
		t.Fatalf("unexpected error: %v", out.Err)
	}
	want := "2026-07-07T23:30:00+09:00[Asia/Tokyo][u-ca=japanese]"
	if out.IXDTF != want {
		t.Errorf("IXDTF = %q, want %q", out.IXDTF, want)
	}
}

func TestFormat_InvalidTimeZone(t *testing.T) {
	badTZ := "Not/AZone"
	out := NewIxdtf().Format(0, entity.Extensions{TimeZone: &badTZ})
	if out.OK || out.Err == nil {
		t.Error("expected domain error for invalid time zone")
	}
}

func TestRoundtrip_Lossless(t *testing.T) {
	i := NewIxdtf()
	cases := []struct {
		name     string
		input    string
		lossless bool
	}{
		{"timezone", "2026-07-07T23:30:00+09:00[Asia/Tokyo]", true},
		{"calendar", "2026-07-07T23:30:00+09:00[Asia/Tokyo][u-ca=japanese]", true},
		{"critical-tz", "2026-07-07T23:30:00+09:00[!Asia/Tokyo]", true},
		{"utc", "2026-07-07T14:30:00Z", true},
		// 末尾ゼロの小数秒は FormatNano で正規化され入力と一致しない。
		{"trailing-zero-frac", "2026-07-07T23:30:00.000+09:00[Asia/Tokyo]", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := i.Roundtrip(tc.input, false)
			if !out.Parse.OK {
				t.Fatalf("parse failed: %v", out.Parse.Err)
			}
			if out.Formatted == nil || out.Lossless == nil {
				t.Fatal("formatted / lossless should not be nil on success")
			}
			if *out.Lossless != tc.lossless {
				t.Errorf("lossless = %v (formatted %q), want %v", *out.Lossless, *out.Formatted, tc.lossless)
			}
		})
	}
}

func TestRoundtrip_ParseFailure(t *testing.T) {
	out := NewIxdtf().Roundtrip("not-a-timestamp", false)
	if out.Parse.OK {
		t.Fatal("expected parse failure")
	}
	if out.Formatted != nil || out.Lossless != nil {
		t.Error("formatted / lossless should be nil on parse failure")
	}
}

func TestNow(t *testing.T) {
	i := NewIxdtf()

	utc, err := i.Now("", "")
	if err != nil {
		t.Fatalf("utc: %v", err)
	}
	if !strings.HasSuffix(utc.IXDTF, "Z") {
		t.Errorf("utc IXDTF %q should end with Z", utc.IXDTF)
	}
	if strings.Contains(utc.IXDTF, "[") {
		t.Errorf("utc IXDTF %q should carry no annotation", utc.IXDTF)
	}
	if utc.UnixNano <= 0 {
		t.Errorf("unix_nano = %d", utc.UnixNano)
	}

	tokyo, err := i.Now("Asia/Tokyo", "")
	if err != nil {
		t.Fatalf("tokyo: %v", err)
	}
	if !strings.Contains(tokyo.IXDTF, "[Asia/Tokyo]") {
		t.Errorf("tokyo IXDTF %q missing [Asia/Tokyo]", tokyo.IXDTF)
	}

	cal, err := i.Now("Asia/Tokyo", "japanese")
	if err != nil {
		t.Fatalf("calendar: %v", err)
	}
	if !strings.Contains(cal.IXDTF, "[u-ca=japanese]") {
		t.Errorf("calendar IXDTF %q missing [u-ca=japanese]", cal.IXDTF)
	}

	if _, err := i.Now("Not/AZone", ""); err == nil {
		t.Error("expected error for invalid time zone")
	}
}
