package main

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestBackoff_NoJitter(t *testing.T) {
	b := Backoff{Base: 10 * time.Millisecond, Max: 80 * time.Millisecond, Factor: 2, Jitter: 0}
	got := []time.Duration{
		b.Duration(0),
		b.Duration(1),
		b.Duration(2),
		b.Duration(3),
		b.Duration(4), // 封頂
	}
	want := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		40 * time.Millisecond,
		80 * time.Millisecond,
		80 * time.Millisecond,
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("attempt %d: got %v, want %v", i, got[i], want[i])
		}
	}
}

func TestIsRetryable_HTTP(t *testing.T) {
	// 5xx
	err := &HTTPStatusError{StatusCode: 502}
	if !IsRetryable(err) {
		t.Fatal("502 should be retryable")
	}
	// 4xx
	err = &HTTPStatusError{StatusCode: 400}
	if IsRetryable(err) {
		t.Fatal("400 should not be retryable")
	}
	// 429
	err = &HTTPStatusError{StatusCode: http.StatusTooManyRequests}
	if !IsRetryable(err) {
		t.Fatal("429 should be retryable")
	}
}

type tmpNetErr struct{}

func (tmpNetErr) Error() string   { return "net temp" }
func (tmpNetErr) Timeout() bool   { return false }
func (tmpNetErr) Temporary() bool { return true }

type timeoutNetErr struct{}

func (timeoutNetErr) Error() string   { return "net timeout" }
func (timeoutNetErr) Timeout() bool   { return true }
func (timeoutNetErr) Temporary() bool { return false }

func TestIsRetryable_netError(t *testing.T) {
	var e net.Error = tmpNetErr{}
	if !IsRetryable(e) {
		t.Fatal("temporary net.Error should be retryable")
	}
	e = timeoutNetErr{}
	if !IsRetryable(e) {
		t.Fatal("timeout net.Error should be retryable")
	}
}

func TestRetry_SucceedsAfterRetries(t *testing.T) {
	// 固定亂數以穩定 jitter（或直接 Jitter=0）
	b := Backoff{Base: 1 * time.Millisecond, Max: 3 * time.Millisecond, Factor: 2, Jitter: 0, r: rand.New(rand.NewSource(1))}
	ctx := context.Background()
	attempts := 0
	fn := func(context.Context) error {
		attempts++
		if attempts < 3 {
			return ErrTemporary // 前兩次暫時性錯誤
		}
		return nil
	}
	err := Retry(ctx, 5, b, fn, IsRetryable)
	if err != nil {
		t.Fatalf("expected success, got error %v", err)
	}
	if attempts != 3 {
		t.Fatalf("attempts got %d, want 3", attempts)
	}
}

func TestRetry_StopsOnNonRetryable(t *testing.T) {
	b := Backoff{Base: 1 * time.Millisecond, Max: 2 * time.Millisecond, Factor: 2, Jitter: 0}
	ctx := context.Background()
	attempts := 0
	fn := func(context.Context) error {
		attempts++
		return ErrBadRequest // 不可重試
	}
	err := Retry(ctx, 5, b, fn, IsRetryable)
	if !errors.Is(err, ErrBadRequest) {
		t.Fatalf("want ErrBadRequest, got %v", err)
	}
	if attempts != 1 {
		t.Fatalf("attempts got %d, want 1", attempts)
	}
}

func TestRetry_CancelContext(t *testing.T) {
	b := Backoff{Base: 50 * time.Millisecond, Max: 50 * time.Millisecond, Factor: 2, Jitter: 0}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	fn := func(context.Context) error { return ErrTemporary }
	err := Retry(ctx, 10, b, fn, IsRetryable)
	if err == nil || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("want context.DeadlineExceeded, got %v", err)
	}
}
