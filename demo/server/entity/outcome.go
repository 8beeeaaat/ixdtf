package entity

import "time"

// ParseOutcome は Parse / Validate のドメイン結果。
//
// 解析失敗は Go の error 戻り値ではなく Err フィールドに保持する。
// 不正入力の観察はこのアプリのドメイン上の正常系であり (D-1)、controller が
// 200 / 4xx を選り分ける分岐を持たずに済むようにするため。
type ParseOutcome struct {
	OK            bool
	Time          time.Time
	Ext           Extensions
	OffsetSeconds int
	Err           error
}

// FormatOutcome は Format のドメイン結果。
type FormatOutcome struct {
	OK    bool
	IXDTF string
	Err   error
}

// RoundtripOutcome は Parse → Format 往復のドメイン結果。
type RoundtripOutcome struct {
	Parse     ParseOutcome
	Formatted *string // Parse 成功時のみ設定
	Lossless  *bool   // Formatted == 入力文字列
}

// FormattedNow はサーバー現在時刻の IXDTF 表記。
type FormattedNow struct {
	IXDTF    string
	UnixNano int64
}
