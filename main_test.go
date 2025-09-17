package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "healthz should return ok",
			url:        "/healthz",
			wantStatus: http.StatusOK,
			wantBody:   "ok\n",
		},
		{
			name:       "unknown path should return 404",
			url:        "/notfound",
			wantStatus: http.StatusNotFound,
			wantBody:   "404 page not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()

			// 使用我們的 handler 註冊路由
			mux := http.NewServeMux()
			mux.HandleFunc("/healthz", healthHandler)

			mux.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", w.Code, tt.wantStatus)
			}
			if w.Body.String() != tt.wantBody {
				t.Errorf("got body %q, want %q", w.Body.String(), tt.wantBody)
			}
		})
	}
}
