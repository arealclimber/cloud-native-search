package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	resp := SearchResponse{
		Query: query,
		Hits: []SearchResult{
			{ID: 1, Title: "Learning Go"},
			{ID: 2, Title: "Go Concurrency Patterns"},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("encode error: %v", err), http.StatusInternalServerError)
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
