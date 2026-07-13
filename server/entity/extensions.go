// Package entity はドメインオブジェクトを定義する。
// ixdtf ライブラリ (RFC 9557) の型との双方向変換もここに置く。
package entity

import (
	"sort"
	"time"

	"github.com/8beeeaaat/ixdtf"
)

// ExtensionTag は IXDTF 拡張タグ 1 件 (キー / 値 / critical フラグ)。
type ExtensionTag struct {
	Key      string
	Value    string
	Critical bool
}

// Extensions は IXDTF サフィックス情報のドメイン表現。
// タグは順序付きスライスで保持する (D-4: 表示の安定性と "!" critical
// フラグをタグ単位で持たせるため)。
type Extensions struct {
	TimeZone         *string
	TimeZoneCritical bool
	Tags             []ExtensionTag
}

// ExtensionsFromIXDTF は ixdtf ライブラリの map ベース拡張を
// ドメインの順序付きスライスへ変換する。
//
// ライブラリが map で返すため元の出現順は失われている。JSON 全経路での
// 順序安定性 (D-4) を保つため、タグはキー昇順ソートで決定的に並べる。
func ExtensionsFromIXDTF(ix *ixdtf.IXDTFExtensions) Extensions {
	var ext Extensions
	if ix == nil {
		ext.Tags = []ExtensionTag{}
		return ext
	}
	if ix.Location != nil {
		name := ix.Location.String()
		ext.TimeZone = &name
	}
	ext.TimeZoneCritical = ix.CriticalLocation

	keys := make([]string, 0, len(ix.Tags))
	for k := range ix.Tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	ext.Tags = make([]ExtensionTag, 0, len(keys))
	for _, k := range keys {
		ext.Tags = append(ext.Tags, ExtensionTag{
			Key:      k,
			Value:    ix.Tags[k],
			Critical: ix.Critical[k],
		})
	}
	return ext
}

// IntoIXDTF はドメイン拡張を ixdtf ライブラリ型へ変換する。
// loc はタイムゾーン注釈に用いる解決済みロケーション
// (呼び出し側が time.LoadLocation 済み。未指定なら nil)。
func (e Extensions) IntoIXDTF(loc *time.Location) *ixdtf.IXDTFExtensions {
	tags := make(map[string]string, len(e.Tags))
	critical := make(map[string]bool, len(e.Tags))
	for _, t := range e.Tags {
		tags[t.Key] = t.Value
		if t.Critical {
			critical[t.Key] = true
		}
	}
	return ixdtf.NewIXDTFExtensions(&ixdtf.NewIXDTFExtensionsArgs{
		Location:         loc,
		CriticalLocation: e.TimeZoneCritical,
		Tags:             tags,
		Critical:         critical,
	})
}
