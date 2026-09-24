package presenter

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/8beeeaaat/ixdtf/demo/server/entity"
)

func TestToParseResponse_OK(t *testing.T) {
	tz := "Asia/Tokyo"
	tm := time.Date(2026, 7, 7, 23, 30, 0, 0, time.FixedZone("JST", 9*3600))
	o := entity.ParseOutcome{
		OK:            true,
		Time:          tm,
		OffsetSeconds: 9 * 3600,
		Ext: entity.Extensions{
			TimeZone:         &tz,
			TimeZoneCritical: true,
			Tags:             []entity.ExtensionTag{{Key: "u-ca", Value: "japanese", Critical: false}},
		},
	}
	resp := ToParseResponse(o)
	if !resp.Ok {
		t.Fatal("Ok should be true")
	}
	if resp.Result == nil {
		t.Fatal("Result should not be nil")
	}
	if resp.Error != nil {
		t.Error("Error should be nil on success")
	}
	r := resp.Result
	if r.Rfc3339 != "2026-07-07T23:30:00+09:00" {
		t.Errorf("rfc3339 = %q", r.Rfc3339)
	}
	if r.UnixNano != strconv.FormatInt(tm.UnixNano(), 10) {
		t.Errorf("unix_nano = %q", r.UnixNano)
	}
	if r.OffsetSeconds != 32400 {
		t.Errorf("offset = %d", r.OffsetSeconds)
	}
	if r.TimeZone == nil || *r.TimeZone != "Asia/Tokyo" {
		t.Errorf("time_zone = %v", r.TimeZone)
	}
	if !r.TimeZoneCritical {
		t.Error("time_zone_critical should be true")
	}
	if len(r.Tags) != 1 || r.Tags[0].Key != "u-ca" || r.Tags[0].Value != "japanese" || r.Tags[0].Critical {
		t.Errorf("tags = %+v", r.Tags)
	}
}

func TestToParseResponse_Error(t *testing.T) {
	resp := ToParseResponse(entity.ParseOutcome{OK: false, Err: errors.New("boom")})
	if resp.Ok {
		t.Error("Ok should be false")
	}
	if resp.Result != nil {
		t.Error("Result should be nil on error")
	}
	if resp.Error == nil || resp.Error.Message != "boom" {
		t.Errorf("error = %+v", resp.Error)
	}
}

func TestToFormatResponse(t *testing.T) {
	ok := ToFormatResponse(entity.FormatOutcome{OK: true, IXDTF: "2026-07-07T23:30:00+09:00[Asia/Tokyo]"})
	if !ok.Ok || ok.Ixdtf == nil || *ok.Ixdtf != "2026-07-07T23:30:00+09:00[Asia/Tokyo]" {
		t.Errorf("ok resp = %+v", ok)
	}
	if ok.Error != nil {
		t.Error("Error should be nil on success")
	}

	bad := ToFormatResponse(entity.FormatOutcome{OK: false, Err: errors.New("nope")})
	if bad.Ok {
		t.Error("Ok should be false")
	}
	if bad.Ixdtf != nil {
		t.Error("Ixdtf should be nil on error")
	}
	if bad.Error == nil || bad.Error.Message != "nope" {
		t.Errorf("error = %+v", bad.Error)
	}
}

func TestToRoundtripResponse(t *testing.T) {
	formatted := "2026-07-07T23:30:00+09:00[Asia/Tokyo]"
	lossless := true
	o := entity.RoundtripOutcome{
		Parse:     entity.ParseOutcome{OK: true, Time: time.Unix(0, 0).UTC()},
		Formatted: &formatted,
		Lossless:  &lossless,
	}
	resp := ToRoundtripResponse(o)
	if !resp.Parse.Ok {
		t.Error("parse should be ok")
	}
	if resp.Formatted == nil || *resp.Formatted != formatted {
		t.Errorf("formatted = %v", resp.Formatted)
	}
	if resp.Lossless == nil || !*resp.Lossless {
		t.Errorf("lossless = %v", resp.Lossless)
	}
}

func TestToNowResponse(t *testing.T) {
	// int64 が JS の MAX_SAFE_INTEGER を超えても 10 進文字列で運ぶ (D-3)。
	resp := ToNowResponse(entity.FormattedNow{IXDTF: "x", UnixNano: 1783434600123456789})
	if resp.Ixdtf != "x" {
		t.Errorf("ixdtf = %q", resp.Ixdtf)
	}
	if resp.UnixNano != "1783434600123456789" {
		t.Errorf("unix_nano = %q", resp.UnixNano)
	}
}
