package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
)

var (
	// 範例：語意化錯誤
	ErrTimeout    = errors.New("timeout")
	ErrTemporary  = errors.New("temporary")
	ErrBadRequest = errors.New("bad request")
)

// Wrap: 以操作名稱包裝錯誤（可多層疊代）
func Wrap(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", op, err)
}

// IsRetryable: 判斷是否「值得重試」
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// 1) 內建語意
	if errors.Is(err, ErrTimeout) || errors.Is(err, ErrTemporary) {
		return true
	}
	if errors.Is(err, ErrBadRequest) {
		return false
	}

	// 2) net.Error: Timeout 或 Temporary
	var ne net.Error
	if errors.As(err, &ne) {
		return ne.Timeout() || ne.Temporary()
	}

	// 3) HTTP 狀態碼（若錯誤內含）
	var he *HTTPStatusError
	if errors.As(err, &he) {
		// 5xx 可重試，429（Too Many Requests）也可視情況重試
		if he.StatusCode == http.StatusTooManyRequests {
			return true
		}
		return he.StatusCode >= 500 && he.StatusCode <= 599
	}

	// 預設保守：不重試
	return false
}

// HTTPStatusError：把下游回來的 HTTP 狀態碼帶進錯誤
type HTTPStatusError struct {
	StatusCode int
	Body       string
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("http status %d: %s", e.StatusCode, e.Body)
}
