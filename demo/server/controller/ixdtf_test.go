package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/8beeeaaat/ixdtf/demo/server/generated"
	"github.com/8beeeaaat/ixdtf/demo/server/usecase/interactor"
)

func newController() *Ixdtf {
	return NewIxdtf(interactor.NewIxdtf())
}

func TestParse_OK(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/ixdtf/parse",
		strings.NewReader(`{"input":"2026-07-07T23:30:00+09:00[Asia/Tokyo]","strict":false}`))
	newController().Parse(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d, body = %s", rec.Code, rec.Body.String())
	}
	var resp generated.ParseResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if !resp.Ok || resp.Result == nil {
		t.Fatalf("resp = %+v", resp)
	}
	if resp.Result.TimeZone == nil || *resp.Result.TimeZone != "Asia/Tokyo" {
		t.Errorf("time_zone = %v", resp.Result.TimeZone)
	}
}

func TestParse_DomainErrorIs200(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/ixdtf/parse",
		strings.NewReader(`{"input":"not-a-timestamp","strict":false}`))
	newController().Parse(rec, req)

	// 解析失敗は 200 のドメイン結果 (A-1)。
	if rec.Code != http.StatusOK {
		t.Fatalf("domain error should be 200, got %d", rec.Code)
	}
	var resp generated.ParseResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Ok {
		t.Error("ok should be false")
	}
	if resp.Error == nil || !strings.Contains(resp.Error.Message, "cannot parse") {
		t.Errorf("error = %+v", resp.Error)
	}
}

func TestParse_BadJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/ixdtf/parse", strings.NewReader(`{bad`))
	newController().Parse(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d", rec.Code)
	}
}

func TestParse_OversizedBodyIs400(t *testing.T) {
	rec := httptest.NewRecorder()
	huge := `{"input":"` + strings.Repeat("a", maxBodyBytes) + `","strict":false}`
	req := httptest.NewRequest(http.MethodPost, "/api/ixdtf/parse", strings.NewReader(huge))
	newController().Parse(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("oversized body should be 400, got %d", rec.Code)
	}
}

func TestParse_MissingRequiredField(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/ixdtf/parse", strings.NewReader(`{"input":"x"}`))
	newController().Parse(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d", rec.Code)
	}
	var e generated.BadRequestError
	_ = json.Unmarshal(rec.Body.Bytes(), &e)
	if e.Message == "" {
		t.Error("expected non-empty message")
	}
}

func TestFormat_OK(t *testing.T) {
	rec := httptest.NewRecorder()
	body := `{"unix_nano":"0","extensions":{"time_zone":"Asia/Tokyo","time_zone_critical":false,"tags":[{"key":"u-ca","value":"japanese","critical":false}]}}`
	req := httptest.NewRequest(http.MethodPost, "/api/ixdtf/format", strings.NewReader(body))
	newController().Format(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d, body = %s", rec.Code, rec.Body.String())
	}
	var resp generated.FormatResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if !resp.Ok || resp.Ixdtf == nil {
		t.Fatalf("resp = %+v", resp)
	}
	if !strings.Contains(*resp.Ixdtf, "[Asia/Tokyo]") || !strings.Contains(*resp.Ixdtf, "[u-ca=japanese]") {
		t.Errorf("ixdtf = %q", *resp.Ixdtf)
	}
}

func TestFormat_BadUnixNano(t *testing.T) {
	rec := httptest.NewRecorder()
	body := `{"unix_nano":"not-int","extensions":{"time_zone":null,"time_zone_critical":false,"tags":[]}}`
	req := httptest.NewRequest(http.MethodPost, "/api/ixdtf/format", strings.NewReader(body))
	newController().Format(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d", rec.Code)
	}
}

func TestFormat_MissingExtensions(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/ixdtf/format", strings.NewReader(`{"unix_nano":"0"}`))
	newController().Format(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d", rec.Code)
	}
}

func TestRoundtrip_OK(t *testing.T) {
	rec := httptest.NewRecorder()
	body := `{"input":"2026-07-07T23:30:00+09:00[Asia/Tokyo]","strict":false}`
	req := httptest.NewRequest(http.MethodPost, "/api/ixdtf/roundtrip", strings.NewReader(body))
	newController().Roundtrip(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d", rec.Code)
	}
	var resp generated.RoundtripResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if !resp.Parse.Ok {
		t.Error("parse should be ok")
	}
	if resp.Formatted == nil {
		t.Error("formatted should not be nil")
	}
	if resp.Lossless == nil || !*resp.Lossless {
		t.Errorf("lossless = %v", resp.Lossless)
	}
}

func TestNow_OK(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/now?time_zone=Asia/Tokyo", nil)
	newController().Now(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d", rec.Code)
	}
	var resp generated.NowResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Ixdtf == "" || resp.UnixNano == "" {
		t.Errorf("resp = %+v", resp)
	}
	if !strings.Contains(resp.Ixdtf, "[Asia/Tokyo]") {
		t.Errorf("ixdtf = %q", resp.Ixdtf)
	}
}

func TestNow_BadTimeZone(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/now?time_zone=Not/AZone", nil)
	newController().Now(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("code = %d", rec.Code)
	}
}
