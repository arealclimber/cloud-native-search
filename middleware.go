package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type LogEntry struct {
	Method   string        `json:"method"`
	Path     string        `json:"path"`
	Duration time.Duration `json:"duration"`
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		// 使用結構化日誌，避免字串格式化
		entry := LogEntry{
			Method:   r.Method,
			Path:     r.URL.Path,
			Duration: duration,
		}

		// 只記錄慢請求或取樣記錄
		if duration > 100*time.Millisecond {
			if data, err := json.Marshal(entry); err == nil {
				log.Println(string(data))
			}
		}
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
