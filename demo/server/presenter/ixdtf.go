// Package presenter はドメイン型 (entity) を generated レスポンス型へ変換する。
// unix_nano の 10 進文字列化 (D-3) や null / optional の変換を担う。
package presenter

import (
	"strconv"
	"time"

	"github.com/8beeeaaat/ixdtf/demo/server/entity"
	"github.com/8beeeaaat/ixdtf/demo/server/generated"
)

// ToParseResponse は ParseOutcome を generated.ParseResponse へ変換する。
func ToParseResponse(o entity.ParseOutcome) generated.ParseResponse {
	if !o.OK {
		return generated.ParseResponse{
			Ok:    false,
			Error: toParseError(o.Err),
		}
	}
	result := toParseResult(o)
	return generated.ParseResponse{
		Ok:     true,
		Result: &result,
	}
}

// ToFormatResponse は FormatOutcome を generated.FormatResponse へ変換する。
func ToFormatResponse(o entity.FormatOutcome) generated.FormatResponse {
	if !o.OK {
		return generated.FormatResponse{
			Ok:    false,
			Ixdtf: nil,
			Error: toParseError(o.Err),
		}
	}
	s := o.IXDTF
	return generated.FormatResponse{
		Ok:    true,
		Ixdtf: &s,
	}
}

// ToRoundtripResponse は RoundtripOutcome を generated.RoundtripResponse へ変換する。
func ToRoundtripResponse(o entity.RoundtripOutcome) generated.RoundtripResponse {
	return generated.RoundtripResponse{
		Parse:     ToParseResponse(o.Parse),
		Formatted: o.Formatted,
		Lossless:  o.Lossless,
	}
}

// ToNowResponse は FormattedNow を generated.NowResponse へ変換する。
func ToNowResponse(n entity.FormattedNow) generated.NowResponse {
	return generated.NowResponse{
		Ixdtf:    n.IXDTF,
		UnixNano: strconv.FormatInt(n.UnixNano, 10),
	}
}

func toParseResult(o entity.ParseOutcome) generated.ParseResult {
	return generated.ParseResult{
		Rfc3339:          o.Time.Format(time.RFC3339Nano),
		UnixNano:         strconv.FormatInt(o.Time.UnixNano(), 10),
		OffsetSeconds:    o.OffsetSeconds,
		TimeZone:         o.Ext.TimeZone,
		TimeZoneCritical: o.Ext.TimeZoneCritical,
		Tags:             toTags(o.Ext.Tags),
	}
}

func toTags(tags []entity.ExtensionTag) []generated.ExtensionTag {
	out := make([]generated.ExtensionTag, 0, len(tags))
	for _, t := range tags {
		out = append(out, generated.ExtensionTag{
			Key:      t.Key,
			Value:    t.Value,
			Critical: t.Critical,
		})
	}
	return out
}

func toParseError(err error) *generated.ParseError {
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	return &generated.ParseError{Message: msg}
}
