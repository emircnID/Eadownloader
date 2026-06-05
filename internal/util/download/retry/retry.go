package retry

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

func IsStatus(statusCode int) bool {
	switch statusCode {
	case http.StatusRequestTimeout,
		http.StatusTooEarly,
		http.StatusTooManyRequests,
		http.StatusInternalServerError,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	default:
		return false
	}
}

func Sleep(ctx context.Context, attempt int, headers http.Header) error {
	delay := time.Duration(1<<min(attempt, 4)) * time.Second

	if headers != nil {
		if retryAfter := headers.Get("Retry-After"); retryAfter != "" {
			if seconds, err := strconv.Atoi(retryAfter); err == nil && seconds > 0 {
				delay = time.Duration(seconds) * time.Second
			} else if retryTime, err := http.ParseTime(retryAfter); err == nil {
				if until := time.Until(retryTime); until > 0 {
					delay = until
				}
			}
		}
	}

	if delay > 30*time.Second {
		delay = 30 * time.Second
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(delay):
		return nil
	}
}
