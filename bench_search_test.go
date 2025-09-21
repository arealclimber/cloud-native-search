package main

import (
	"testing"
)

// 偽資料來源
var smallHits = []SearchResult{
	{ID: 1, Title: "Learning Go"},
	{ID: 2, Title: "Go Concurrency Patterns"},
	{ID: 3, Title: "High Performance Go"},
}

// 1A) 動態 append（可能反覆擴容）
func buildHitsAppend(src []SearchResult) []SearchResult {
	out := []SearchResult{}
	for _, h := range src {
		out = append(out, h)
	}
	return out
}

// 1B) 預先配置容量（避免擴容）
func buildHitsPrealloc(src []SearchResult) []SearchResult {
	out := make([]SearchResult, 0, len(src))
	for _, h := range src {
		out = append(out, h)
	}
	return out
}

func BenchmarkBuildHitsAppend(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = buildHitsAppend(smallHits)
	}
}

func BenchmarkBuildHitsPrealloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = buildHitsPrealloc(smallHits)
	}
}
