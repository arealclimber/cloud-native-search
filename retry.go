package main

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"
)

type Backoff struct {
	Base   time.Duration // 初始等待，例如 50ms
	Max    time.Duration // 上限，例如 2s
	Factor float64       // 指數成長倍率，例如 2.0
	Jitter float64       // 抖動比例 0~1，例如 0.2 代表 ±20%
	// 測試用：可注入 rand 以固定結果
	r *rand.Rand
}

func (b Backoff) withRand() Backoff {
	if b.r == nil {
		b2 := b
		b2.r = rand.New(rand.NewSource(time.Now().UnixNano()))
		return b2
	}
	return b
}

// 第 n 次（從 0 起算）的等待時間
func (b Backoff) Duration(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	// 指數成長
	mult := math.Pow(b.Factor, float64(attempt))
	d := time.Duration(float64(b.Base) * mult)
	if d > b.Max {
		d = b.Max
	}
	// 抖動：±Jitter
	if b.Jitter > 0 {
		b2 := b.withRand()
		j := (b2.r.Float64()*2 - 1) * b.Jitter // [-Jitter, +Jitter]
		d = time.Duration(float64(d) * (1 + j))
	}
	if d < 0 {
		d = 0
	}
	return d
}

// Retry: 以 backoff + jitter 重試 fn；遇到非重試錯誤、或 ctx 結束則返回
func Retry(ctx context.Context, maxAttempts int, b Backoff, fn func(context.Context) error, shouldRetry func(error) bool) error {
	if maxAttempts <= 0 {
		return errors.New("maxAttempts must be > 0")
	}
	var lastErr error
	for attempt := 0; attempt < maxAttempts; attempt++ {
		// 先檢查是否已取消
		select {
		case <-ctx.Done():
			return Wrap("retry canceled", ctx.Err())
		default:
		}

		err := fn(ctx)
		if err == nil {
			return nil
		}
		lastErr = err

		if !shouldRetry(err) || attempt == maxAttempts-1 {
			return lastErr
		}

		wait := b.Duration(attempt)
		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			timer.Stop()
			return Wrap("retry canceled", ctx.Err())
		case <-timer.C:
		}
	}
	return lastErr
}
