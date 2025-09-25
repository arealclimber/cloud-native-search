package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func searchHandler(searchService SearchService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "missing query parameter: q", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		resp, err := searchService.Search(ctx, query)
		if err != nil {
			http.Error(w, fmt.Sprintf("search error: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, fmt.Sprintf("encode error: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func main() {
	cfg := LoadConfig()

	// 初始化搜尋服務
	searchService := &FakeSearchService{}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/search", searchHandler(searchService))

	// 新增 metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	// 中介層：metrics → logging → recovery
	handler := MetricsMiddleware(LoggingMiddleware(RecoveryMiddleware(mux)))

	// 啟動 pprof (只有 build tag=debug 才會啟動)
	StartPprof()

	log.Printf("Server listening on %s", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, handler); err != nil {
		log.Fatal(err)
	}
}
