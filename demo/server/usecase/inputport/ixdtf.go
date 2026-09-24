// Package inputport は usecase 入口のインターフェース (controller が依存する契約) を定義する。
package inputport

import "github.com/8beeeaaat/ixdtf/demo/server/entity"

// IxdtfInputPort は IXDTF の解析・整形ユースケースの契約。
// interactor が実装し、controller が依存する。
type IxdtfInputPort interface {
	Parse(input string, strict bool, validateOnly bool) entity.ParseOutcome
	Format(unixNano int64, ext entity.Extensions) entity.FormatOutcome
	Roundtrip(input string, strict bool) entity.RoundtripOutcome
	Now(timeZone string, calendar string) (entity.FormattedNow, error)
}
