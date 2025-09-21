//go:build slow

package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// 目前的實作：handler 會在 2s 超時、下游模擬 3s -> 必定 504
func TestSearchHandler_ServerTimeout(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/search?q=golang", nil)
	rr := httptest.NewRecorder()

	// 直接呼叫目前的 handler
	searchHandler(rr, req)

	if rr.Code != http.StatusGatewayTimeout {
		t.Fatalf("status got %d, want %d", rr.Code, http.StatusGatewayTimeout)
	}
	want := "search timeout\n"
	if rr.Body.String() != want {
		t.Fatalf("body got %q, want %q", rr.Body.String(), want)
	}
}

// 客戶端主動取消（早於 2s timeout），也應拿到 504
func TestSearchHandler_ClientCancel(t *testing.T) {
	t.Parallel()

	// 建立一個會在 200ms 取消的 context，包在 request 裡
	parent := context.Background()
	ctx, cancel := context.WithTimeout(parent, 200*time.Millisecond)
	defer cancel()

	req := httptest.NewRequest(http.MethodGet, "/search?q=golang", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	searchHandler(rr, req)

	if rr.Code != http.StatusGatewayTimeout {
		t.Fatalf("status got %d, want %d", rr.Code, http.StatusGatewayTimeout)
	}
}

// 缺少 q 參數 -> 400
func TestSearchHandler_BadRequest(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	rr := httptest.NewRecorder()

	searchHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status got %d, want %d", rr.Code, http.StatusBadRequest)
	}
}
