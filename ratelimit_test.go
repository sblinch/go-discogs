package discogs

import (
	"context"
	"io"
	"testing"
	"time"
)

func TestRateLimit_Update(t *testing.T) {
	rl := &RateLimit{}
	tests := []struct {
		name                   string
		total, used, remaining int
	}{
		{"initial", 10, 6, 4},
		{"change", 10, 7, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl.Update(tt.total, tt.used, tt.remaining)
			total, used, remaining, updated := rl.Get()
			if total != tt.total || used != tt.used || remaining != tt.remaining {
				t.Errorf("set values %d, %d, %d, got values %d, %d, %d", tt.total, tt.used, tt.remaining, total, used, remaining)
			}
			if updated.Before(time.Now().Add(-time.Second)) || updated.After(time.Now()) {
				t.Errorf("unexpected update time %s", updated.String())
			}
		})
	}
}

func TestRateLimit_Call(t *testing.T) {
	rl := &RateLimit{}
	ctx := context.Background()

	tests := []struct {
		name                   string
		total, used, remaining int
		fresh                  bool
		attempts               []error
		expectErr              error
		expectDelay            time.Duration
	}{
		{"fresh data and remaining", 10, 6, 4, true, []error{nil}, nil, 0},
		{"fresh data and zero remaining", 10, 10, 0, true, []error{nil}, nil, minimumRateLimitDelay},
		{"stale data and zero remaining", 10, 10, 0, false, []error{nil}, nil, 0},

		{"fresh data and remaining with error", 10, 6, 4, true, []error{io.ErrUnexpectedEOF}, io.ErrUnexpectedEOF, 0},
		{"fresh data and zero remaining with error", 10, 10, 0, true, []error{io.ErrUnexpectedEOF}, io.ErrUnexpectedEOF, minimumRateLimitDelay},

		{"fresh data and zero remaining and rate limited", 10, 10, 0, true, []error{ErrTooManyRequests, nil}, nil, minimumRateLimitDelay * 3},
		{"fresh data and zero remaining and rate limited twice", 10, 10, 0, true, []error{ErrTooManyRequests, ErrTooManyRequests, nil}, nil, minimumRateLimitDelay * 7},

		{"fresh data and remaining and rate limited", 10, 6, 4, true, []error{ErrTooManyRequests, nil}, nil, minimumRateLimitDelay},
		{"fresh data and remaining and rate limited twice", 10, 6, 4, true, []error{ErrTooManyRequests, ErrTooManyRequests, nil}, nil, minimumRateLimitDelay * 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rl.Update(tt.total, tt.used, tt.remaining)
			if !tt.fresh {
				rl.updated = time.Now().Add(-time.Minute)
			}
			attempts := tt.attempts[:]
			slept := time.Duration(0)

			request := func() error {
				err := attempts[0]
				attempts = attempts[1:]
				return err
			}

			sleep := func(ctx context.Context, duration time.Duration) error {
				slept += duration
				return nil
			}

			err := rl.call(ctx, request, sleep)

			if err != tt.expectErr {
				t.Errorf("Expected error %v, got error %v", tt.expectErr, err)
			}
			if slept != tt.expectDelay {
				t.Errorf("Expected delay %v, got delay %v", tt.expectDelay.String(), slept.String())
			}

		})
	}
}
