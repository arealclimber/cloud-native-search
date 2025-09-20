package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// 假資料
type SearchResult struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type SearchResponse struct {
	Query string         `json:"query"`
	Hits  []SearchResult `json:"hits"`
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "missing query parameter: q", http.StatusBadRequest)
		return
	}

	// 建立帶有 timeout 的 context
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// 模擬下游呼叫（例如 Elasticsearch）
	resultCh := make(chan SearchResponse, 1)
	go func() {
		// 假裝需要 3 秒（比 timeout 長）
		time.Sleep(3 * time.Second)
		resultCh <- SearchResponse{
			Query: query,
			Hits: []SearchResult{
				{ID: 1, Title: "Learning Go"},
				{ID: 2, Title: "Go Concurrency Patterns"},
			},
		}
	}()

	select {
	case <-ctx.Done():
		// timeout 或被取消
		http.Error(w, "search timeout", http.StatusGatewayTimeout)
		return
	case resp := <-resultCh:
		// 正常拿到結果
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, fmt.Sprintf("encode error: %v", err), http.StatusInternalServerError)
		}
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func main() {
	cfg := LoadConfig()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/search", searchHandler)

	handler := LoggingMiddleware(RecoveryMiddleware(mux))

	log.Printf("Server listening on %s", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, handler); err != nil {
		log.Fatal(err)
	}
}
