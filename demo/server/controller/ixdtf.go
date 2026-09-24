// Package controller は HTTP ハンドラ群。JSON デコード → inputport 呼び出し
// → presenter → JSON エンコードを担う。リクエスト自体の不備のみ 4xx を返し (A-2)、
// 解析失敗は presenter 経由で 200 のドメイン結果 (ok: false) として返す (A-1)。
package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/8beeeaaat/ixdtf/demo/server/entity"
	"github.com/8beeeaaat/ixdtf/demo/server/generated"
	"github.com/8beeeaaat/ixdtf/demo/server/presenter"
	"github.com/8beeeaaat/ixdtf/demo/server/usecase/inputport"
)

// Ixdtf は IXDTF API の HTTP ハンドラ群。
type Ixdtf struct {
	usecase inputport.IxdtfInputPort
}

// NewIxdtf は controller を生成する。
func NewIxdtf(usecase inputport.IxdtfInputPort) *Ixdtf {
	return &Ixdtf{usecase: usecase}
}

// Parse は POST /api/ixdtf/parse を処理する。
func (c *Ixdtf) Parse(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Input        *string `json:"input"`
		Strict       *bool   `json:"strict"`
		ValidateOnly *bool   `json:"validate_only"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeBadRequest(w, err.Error())
		return
	}
	if body.Input == nil || body.Strict == nil {
		writeBadRequest(w, "missing required field: input and strict are required")
		return
	}
	validateOnly := body.ValidateOnly != nil && *body.ValidateOnly
	outcome := c.usecase.Parse(*body.Input, *body.Strict, validateOnly)
	writeJSON(w, http.StatusOK, presenter.ToParseResponse(outcome))
}

// Format は POST /api/ixdtf/format を処理する。
func (c *Ixdtf) Format(w http.ResponseWriter, r *http.Request) {
	var body formatRequestBody
	if err := decodeJSON(w, r, &body); err != nil {
		writeBadRequest(w, err.Error())
		return
	}
	if body.UnixNano == nil || body.Extensions == nil {
		writeBadRequest(w, "missing required field: unix_nano and extensions are required")
		return
	}
	unixNano, err := strconv.ParseInt(*body.UnixNano, 10, 64)
	if err != nil {
		writeBadRequest(w, "invalid unix_nano: "+err.Error())
		return
	}
	ext, err := body.Extensions.toEntity()
	if err != nil {
		writeBadRequest(w, err.Error())
		return
	}
	outcome := c.usecase.Format(unixNano, ext)
	writeJSON(w, http.StatusOK, presenter.ToFormatResponse(outcome))
}

// Roundtrip は POST /api/ixdtf/roundtrip を処理する。
func (c *Ixdtf) Roundtrip(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Input  *string `json:"input"`
		Strict *bool   `json:"strict"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		writeBadRequest(w, err.Error())
		return
	}
	if body.Input == nil || body.Strict == nil {
		writeBadRequest(w, "missing required field: input and strict are required")
		return
	}
	outcome := c.usecase.Roundtrip(*body.Input, *body.Strict)
	writeJSON(w, http.StatusOK, presenter.ToRoundtripResponse(outcome))
}

// Now は GET /api/now を処理する。
// timeZone が不正な場合のみ 400 (LoadLocation 失敗はリクエスト不備扱い)。
func (c *Ixdtf) Now(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	now, err := c.usecase.Now(q.Get("time_zone"), q.Get("calendar"))
	if err != nil {
		writeBadRequest(w, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, presenter.ToNowResponse(now))
}

// formatRequestBody は FormatRequest のデコード用。必須フィールド欠落を
// 検出するためポインタで受ける。
type formatRequestBody struct {
	UnixNano   *string           `json:"unix_nano"`
	Extensions *formatExtensions `json:"extensions"`
}

type formatExtensions struct {
	TimeZone         *string     `json:"time_zone"`
	TimeZoneCritical *bool       `json:"time_zone_critical"`
	Tags             []formatTag `json:"tags"`
}

type formatTag struct {
	Key      *string `json:"key"`
	Value    *string `json:"value"`
	Critical *bool   `json:"critical"`
}

// toEntity は FormatRequest.extensions を entity.Extensions へ変換する。
// タグの必須フィールド (key / value / critical) 欠落はリクエスト不備。
func (e *formatExtensions) toEntity() (entity.Extensions, error) {
	ext := entity.Extensions{TimeZone: e.TimeZone}
	if e.TimeZoneCritical != nil {
		ext.TimeZoneCritical = *e.TimeZoneCritical
	}
	ext.Tags = make([]entity.ExtensionTag, 0, len(e.Tags))
	for _, t := range e.Tags {
		if t.Key == nil || t.Value == nil || t.Critical == nil {
			return entity.Extensions{}, errors.New("invalid extension tag: key, value, critical are required")
		}
		ext.Tags = append(ext.Tags, entity.ExtensionTag{
			Key:      *t.Key,
			Value:    *t.Value,
			Critical: *t.Critical,
		})
	}
	return ext, nil
}

// maxBodyBytes はリクエストボディの上限。IXDTF 文字列は高々数 KB のため
// 64 KiB で十分な余裕がある。超過時は decode エラーとなり 400 を返す。
const maxBodyBytes = 64 << 10

func decodeJSON(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	return json.NewDecoder(r.Body).Decode(v)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeBadRequest(w http.ResponseWriter, msg string) {
	writeJSON(w, http.StatusBadRequest, generated.BadRequestError{Message: msg})
}
