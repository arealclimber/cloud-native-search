package main

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunWorkerPool_AllSuccess(t *testing.T) {
	ctx := context.Background()
	n := 5
	items := make([]WorkItem, n)
	for i := 0; i < n; i++ {
		id := i
		items[i] = WorkItem{
			ID: id,
			Task: func(ctx context.Context) error {
				return nil
			},
		}
	}
	errs := RunWorkerPool(ctx, 2, items)
	for _, err := range errs {
		if err != nil {
			t.Fatalf("expected nil error, got %v", err)
		}
	}
}

func TestRunWorkerPool_WithError(t *testing.T) {
	ctx := context.Background()
	items := []WorkItem{
		{ID: 0, Task: func(ctx context.Context) error { return errors.New("boom") }},
		{ID: 1, Task: func(ctx context.Context) error { return nil }},
	}
	errs := RunWorkerPool(ctx, 1, items)
	if errs[0] == nil || errs[1] != nil {
		t.Fatalf("unexpected errs: %#v", errs)
	}
}

func TestRunWorkerPool_CancelEarly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var counter int32
	items := []WorkItem{
		{ID: 0, Task: func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			time.Sleep(50 * time.Millisecond)
			return nil
		}},
		{ID: 1, Task: func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		}},
	}

	// 立刻取消
	cancel()

	errs := RunWorkerPool(ctx, 2, items)
	if counter > 1 {
		t.Fatalf("expected at most 1 task executed, got %d", counter)
	}
	if errs[0] != nil && !errors.Is(errs[0], context.Canceled) {
		// 部分可能是 context.Canceled
		t.Logf("errs: %#v", errs)
	}
}
