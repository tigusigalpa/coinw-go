// Package ratelimit implements the documented CoinW throttling policy:
//
//	Spot global:    100 req/s per IP AND 300 req/s per UID
//	Futures Group A (perpumPublic market data): 30 req/s per IP
//	Futures Group B (everything else futures):  100 req/s per UID
//	Affiliate:      5 req/s per IP and UID
//	Per-interface:  each endpoint's documented limit
//
// All governing buckets must grant a token before the request proceeds.
// The default limiter is in-process — inject a distributed implementation
// of Limiter for multi-worker deployments.
package ratelimit

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/tigusigalpa/coinw-go/endpoints"
)

// Limiter blocks until a request may proceed or the context is cancelled.
type Limiter interface {
	Acquire(ctx context.Context, endpoint endpoints.Definition, uid string) error
}

// bucket is a sliding-window log guarded by the parent mutex: it grants at
// most limit requests in any window. A token bucket with capacity N would
// allow a burst of N plus N refilled tokens, about 2N, within one second.
type bucket struct {
	limit  int
	window time.Duration
	grants []time.Time // the last limit grant times, a ring once full
	oldest int
}

// newBucket allows rate requests per second, or one per 1/rate seconds when
// rate is fractional or below one.
func newBucket(rate float64) *bucket {
	limit := int(math.Max(1, math.Floor(rate)))
	window := time.Duration(math.Ceil(float64(limit) / rate * float64(time.Second)))
	return &bucket{limit: limit, window: window, grants: make([]time.Time, 0, limit)}
}

// waitDuration returns how long until one more grant fits in the window.
func (b *bucket) waitDuration(now time.Time) time.Duration {
	if len(b.grants) < b.limit {
		return 0
	}
	if w := b.grants[b.oldest].Add(b.window).Sub(now); w > 0 {
		return w
	}
	return 0
}

func (b *bucket) commit(now time.Time) {
	if len(b.grants) < b.limit {
		b.grants = append(b.grants, now)
		return
	}
	b.grants[b.oldest] = now
	b.oldest = (b.oldest + 1) % b.limit
}

// Clock abstracts time for deterministic tests.
type Clock interface {
	Now() time.Time
}

type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }

// LimiterOption customises the default limiter (tests inject a fake clock).
type LimiterOption func(*InMemoryLimiter)

// WithClock overrides the time source.
func WithClock(c Clock) LimiterOption {
	return func(l *InMemoryLimiter) { l.clock = c }
}

// InMemoryLimiter is the default implementation — per-process token buckets.
// Not shared across processes.
type InMemoryLimiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	clock   Clock
}

// NewInMemoryLimiter creates the default limiter.
func NewInMemoryLimiter(opts ...LimiterOption) *InMemoryLimiter {
	l := &InMemoryLimiter{buckets: map[string]*bucket{}, clock: systemClock{}}
	for _, o := range opts {
		o(l)
	}
	return l
}

// Acquire waits until every governing bucket can grant a token, then commits
// all of them atomically (two-phase: peek slowest wait → commit).
func (l *InMemoryLimiter) Acquire(ctx context.Context, ep endpoints.Definition, uid string) error {
	bs := l.bucketsFor(ep, uid)
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		l.mu.Lock()
		now := l.clock.Now()
		var wait time.Duration
		for _, b := range bs {
			if w := b.waitDuration(now); w > wait {
				wait = w
			}
		}
		if wait <= 0 {
			for _, b := range bs {
				b.commit(now)
			}
			l.mu.Unlock()
			return nil
		}
		l.mu.Unlock()

		timer := time.NewTimer(wait)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
}

func (l *InMemoryLimiter) bucketsFor(ep endpoints.Definition, uid string) []*bucket {
	l.mu.Lock()
	defer l.mu.Unlock()

	type spec struct {
		key string
		rps float64
	}
	var specs []spec

	switch ep.RateGroup {
	case endpoints.GroupSpot:
		specs = append(specs, spec{"spot:ip", 100})
		if uid != "" {
			specs = append(specs, spec{"spot:uid:" + uid, 300})
		}
	case endpoints.GroupFuturesA:
		specs = append(specs, spec{"futuresA:ip", 30})
	case endpoints.GroupFuturesB:
		// An unspecified UID still needs a shared budget for this client.
		specs = append(specs, spec{"futuresB:uid:" + uid, 100})
	case endpoints.GroupAffiliate:
		specs = append(specs, spec{"affiliate:ip", 5})
		if uid != "" {
			specs = append(specs, spec{"affiliate:uid:" + uid, 5})
		}
	}

	ipRate, uidRate := ep.PerIPRPS, ep.PerUIDRPS
	if ipRate == 0 && uidRate == 0 {
		ipRate, uidRate = ep.PerInterfaceRPS, ep.PerInterfaceRPS
	}
	if ipRate > 0 {
		specs = append(specs, spec{"iface:" + ep.ID + ":ip", ipRate})
	}
	if uidRate > 0 {
		specs = append(specs, spec{"iface:" + ep.ID + ":uid:" + uid, uidRate})
	}

	out := make([]*bucket, 0, len(specs))
	for _, s := range specs {
		b, ok := l.buckets[s.key]
		if !ok {
			b = newBucket(s.rps)
			l.buckets[s.key] = b
		}
		out = append(out, b)
	}
	return out
}
