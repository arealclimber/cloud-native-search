package main

import (
	"fmt"
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ok")
}

func main() {
	cfg := LoadConfig()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)

	handler := LoggingMiddleware(RecoveryMiddleware(mux))

	log.Printf("Server listening on %s", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, handler); err != nil {
		log.Fatal(err)
	}
}
