package discogs

import (
	"context"
	"errors"
	"sync"
	"time"
)

type RateLimit struct {
	mu        sync.Mutex
	total     int
	used      int
	remaining int
	updated   time.Time
}

// Update sets the rate limiting parameters received from the headers of a Discogs API call.
func (r *RateLimit) Update(total, used, remaining int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.total = total
	r.used = used
	r.remaining = remaining
	r.updated = time.Now()
}

// Get retrieves the most recent rate limiting parameters and the time at which they were set.
func (r *RateLimit) Get() (total, used, remaining int, updated time.Time) {
	r.mu.Lock()
	defer r.mu.Unlock()

	total = r.total
	used = r.used
	remaining = r.remaining
	updated = r.updated
	return
}

// Call invokes f() when the rate limiting metrics indicate that it's likely safe to do so and, if a rate limiting
// error is returned, repeats the call with exponential backoff until it returns any value other than ErrTooManyRequests.
func (r *RateLimit) Call(ctx context.Context, f func() error) error {

	t := time.NewTimer(time.Minute)
	t.Stop()
	defer t.Stop()

	sleep := func(ctx context.Context, d time.Duration) error {
		t.Reset(d)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			return nil
		}
	}

	return r.call(ctx, f, sleep)

}

const minimumRateLimitDelay = 2500 * time.Millisecond

// call is the inner implementation of Call which accepts a sleep function that can be mocked during testing.
func (r *RateLimit) call(ctx context.Context, f func() error, sleep func(context.Context, time.Duration) error) error {
	delay := minimumRateLimitDelay
	first := true

	for {
		_, _, remaining, when := r.Get()

		// pause if the rate limiting metrics are reasonably fresh and we have no remaining permitted requests, OR if
		// we just received ErrTooManyRequests regardless of how many requests Discogs claims we have remaining;
		// Discogs seems to report the pre-request X-Discogs-Ratelimit-Used value, so we're out of requests when remaining==1
		if !first || time.Now().Sub(when) < 10*time.Second && remaining <= 1 {
			if err := sleep(ctx, delay); err != nil {
				return err
			}
			delay *= 2
		}

		err := f()
		if !errors.Is(err, ErrTooManyRequests) {
			return err
		}
		first = false
	}
}
