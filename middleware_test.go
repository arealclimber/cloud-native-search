package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// 一個會 panic 的 handler，模擬程式發生未捕捉錯誤
func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("boom")
}

func TestRecoveryMiddleware_Returns500OnPanic(t *testing.T) {
	// 用 RecoveryMiddleware 包住一個會 panic 的 handler
	h := RecoveryMiddleware(http.HandlerFunc(panicHandler))

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status code = %d, want %d", rr.Code, http.StatusInternalServerError)
	}

	// 我們在 RecoveryMiddleware 裡用 http.Error 寫入錯誤訊息，
	// 預設為 "internal server error\n"
	wantBody := "internal server error\n"
	if rr.Body.String() != wantBody {
		t.Fatalf("body = %q, want %q", rr.Body.String(), wantBody)
	}
}
