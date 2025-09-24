package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// 總請求數，依路徑與狀態碼區分
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests",
		},
		[]string{"path", "code"},
	)

	// 請求延遲，依路徑區分
	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Request latency (seconds)",
			Buckets: prometheus.DefBuckets, // 預設 buckets (0.005s ~ 10s)
		},
		[]string{"path"},
	)
)

func init() {
	// 註冊到全域 registry
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}
