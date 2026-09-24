// Package interactor は IXDTF ユースケースのビジネスロジック実装。
// ixdtf ライブラリの Parse / Validate / FormatNano を呼び出し、
// strict・validate_only・lossless 判定などのドメイン判断を担う。
package interactor

import (
	"time"

	"github.com/8beeeaaat/ixdtf"
	"github.com/8beeeaaat/ixdtf/demo/server/entity"
	"github.com/8beeeaaat/ixdtf/demo/server/usecase/inputport"
)

// Ixdtf は ixdtf ライブラリを用いた IxdtfInputPort 実装。
type Ixdtf struct{}

var _ inputport.IxdtfInputPort = (*Ixdtf)(nil)

// NewIxdtf は Ixdtf interactor を生成する。
func NewIxdtf() *Ixdtf {
	return &Ixdtf{}
}

// Parse は IXDTF 文字列を解析する。
// validateOnly の場合は先に Validate を実行し、そのエラーを優先する (A-3)。
// 解析エラーは ixdtf ライブラリの原文のまま Outcome.Err に保持する。
func (i *Ixdtf) Parse(input string, strict bool, validateOnly bool) entity.ParseOutcome {
	if validateOnly {
		if err := ixdtf.Validate(input, strict); err != nil {
			return entity.ParseOutcome{Err: err}
		}
	}
	t, ix, err := ixdtf.Parse(input, strict)
	if err != nil {
		return entity.ParseOutcome{Err: err}
	}
	_, offset := t.Zone()
	return entity.ParseOutcome{
		OK:            true,
		Time:          t,
		Ext:           entity.ExtensionsFromIXDTF(ix),
		OffsetSeconds: offset,
	}
}

// Format は unix_nano と拡張から IXDTF 文字列を生成する。
// TimeZone 指定時は time.LoadLocation で時刻の Location に適用する。
// LoadLocation 失敗はドメインエラー (Outcome.Err) として扱う。
func (i *Ixdtf) Format(unixNano int64, ext entity.Extensions) entity.FormatOutcome {
	loc, err := resolveLocation(ext.TimeZone)
	if err != nil {
		return entity.FormatOutcome{Err: err}
	}
	t := time.Unix(0, unixNano).UTC()
	if loc != nil {
		t = t.In(loc)
	}
	s, err := ixdtf.FormatNano(t, ext.IntoIXDTF(loc))
	if err != nil {
		return entity.FormatOutcome{Err: err}
	}
	return entity.FormatOutcome{OK: true, IXDTF: s}
}

// Roundtrip は Parse 成功時のみ FormatNano で往復させ、lossless を判定する。
func (i *Ixdtf) Roundtrip(input string, strict bool) entity.RoundtripOutcome {
	parse := i.Parse(input, strict, false)
	out := entity.RoundtripOutcome{Parse: parse}
	if !parse.OK {
		return out
	}
	// 注釈を再現するためゾーン名から Location を再解決する。
	// 解析が成功しているため通常は解決できる (失敗時は nil のまま整形)。
	loc, _ := resolveLocation(parse.Ext.TimeZone)
	formatted, err := ixdtf.FormatNano(parse.Time, parse.Ext.IntoIXDTF(loc))
	if err != nil {
		return out
	}
	lossless := formatted == input
	out.Formatted = &formatted
	out.Lossless = &lossless
	return out
}

// Now はサーバー現在時刻を IXDTF 表記で返す。
// timeZone が空なら UTC、非空なら time.LoadLocation で解決する。
// LoadLocation 失敗は error 戻り値 (controller が 400 を返す)。
// calendar が非空なら u-ca タグを付与する。
func (i *Ixdtf) Now(timeZone string, calendar string) (entity.FormattedNow, error) {
	loc := time.UTC
	if timeZone != "" {
		l, err := time.LoadLocation(timeZone)
		if err != nil {
			return entity.FormattedNow{}, err
		}
		loc = l
	}
	t := time.Now().In(loc)

	args := &ixdtf.NewIXDTFExtensionsArgs{
		Tags:     map[string]string{},
		Critical: map[string]bool{},
	}
	if loc != time.UTC {
		args.Location = loc
	}
	if calendar != "" {
		args.Tags[ixdtf.ExtensionUnicodeCalendar] = calendar
	}
	s, err := ixdtf.FormatNano(t, ixdtf.NewIXDTFExtensions(args))
	if err != nil {
		return entity.FormattedNow{}, err
	}
	return entity.FormattedNow{IXDTF: s, UnixNano: t.UnixNano()}, nil
}

// resolveLocation は TimeZone 名から *time.Location を解決する。
// nil または空文字なら (nil, nil) を返す (注釈なし)。
func resolveLocation(timeZone *string) (*time.Location, error) {
	if timeZone == nil || *timeZone == "" {
		return nil, nil
	}
	return time.LoadLocation(*timeZone)
}
