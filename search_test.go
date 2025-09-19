package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchHandler(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		wantStatus int
		wantHits   int
	}{
		{
			name:       "valid query",
			query:      "golang",
			wantStatus: http.StatusOK,
			wantHits:   2,
		},
		{
			name:       "missing query",
			query:      "",
			wantStatus: http.StatusBadRequest,
			wantHits:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/search?q="+tt.query, nil)
			w := httptest.NewRecorder()

			searchHandler(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status got %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var resp SearchResponse
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("failed to parse JSON: %v", err)
				}
				if len(resp.Hits) != tt.wantHits {
					t.Errorf("hits got %d, want %d", len(resp.Hits), tt.wantHits)
				}
			}
		})
	}
}
