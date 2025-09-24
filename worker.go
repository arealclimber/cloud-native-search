package main

import (
	"context"
	"sync"
)

// WorkItem 定義一個要處理的任務
type WorkItem struct {
	ID   int
	Task func(context.Context) error
}

// RunWorkerPool 執行一個有限併發的 worker pool
//   - ctx：用於取消整個 pool
//   - maxWorkers：同時最多跑多少 goroutine
//   - items：要處理的任務 slice
func RunWorkerPool(ctx context.Context, maxWorkers int, items []WorkItem) []error {
	var wg sync.WaitGroup
	errs := make([]error, len(items))

	// channel 傳輸任務
	workCh := make(chan WorkItem)

	// 啟動 worker
	for w := 0; w < maxWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range workCh {
				select {
				case <-ctx.Done():
					// 整個 pool 被取消
					return
				default:
					errs[item.ID] = item.Task(ctx)
				}
			}
		}()
	}

	// 丟任務
	go func() {
		defer close(workCh)
		for _, item := range items {
			select {
			case <-ctx.Done():
				return
			case workCh <- item:
			}
		}
	}()

	wg.Wait()
	return errs
}
