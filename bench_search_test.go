package main

import (
	"bytes"
	"encoding/json"
	"testing"
)

type respEnvelope struct {
	Query string         `json:"query"`
	Hits  []SearchResult `json:"hits"`
}

var resp = respEnvelope{
	Query: "golang",
	Hits:  smallHits,
}

// 偽資料來源
var smallHits = []SearchResult{
	{ID: 1, Title: "Learning Go"},
	{ID: 2, Title: "Go Concurrency Patterns"},
	{ID: 3, Title: "High Performance Go"},
}

// 1A) 動態 append（可能反覆擴容）
func buildHitsAppend(src []SearchResult) []SearchResult {
	return append([]SearchResult{}, src...)
}

// 1B) 預先配置容量（避免擴容）
func buildHitsPrealloc(src []SearchResult) []SearchResult {
	out := make([]SearchResult, 0, len(src))
	return append(out, src...)
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

func BenchmarkJSONMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(resp)
	}
}

func BenchmarkJSONEncoder_ReusedBuffer(b *testing.B) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	for i := 0; i < b.N; i++ {
		buf.Reset()
		_ = enc.Encode(resp) // 注意 Encode 會在結尾加 '\n'
	}
}
