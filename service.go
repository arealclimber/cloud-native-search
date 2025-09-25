// service.go
package main

import "context"

// SearchResult 是單筆搜尋結果
type SearchResult struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// SearchResponse 是搜尋回應
type SearchResponse struct {
	Query string         `json:"query"`
	Hits  []SearchResult `json:"hits"`
}

// SearchService 定義搜尋服務介面
type SearchService interface {
	Search(ctx context.Context, query string) (SearchResponse, error)
}

// FakeSearchService 是假的實作，先回固定資料
type FakeSearchService struct{}

func (s *FakeSearchService) Search(ctx context.Context, query string) (SearchResponse, error) {
	return SearchResponse{
		Query: query,
		Hits: []SearchResult{
			{ID: 1, Title: "Learning Go"},
			{ID: 2, Title: "Go Concurrency Patterns"},
		},
	}, nil
}
