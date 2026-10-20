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

// bucket is a classic token bucket guarded by the parent mutex.
type bucket struct {
	capacity float64
	rate     float64 // tokens per second
	tokens   float64
	lastFill time.Time
}

func newBucket(capacity, rate float64, now time.Time) *bucket {
	return &bucket{capacity: capacity, rate: rate, tokens: capacity, lastFill: now}
}

// waitDuration returns how long until one token is available.
func (b *bucket) waitDuration(now time.Time) time.Duration {
	b.refill(now)
	if b.tokens >= 1.0 {
		return 0
	}
	missing := 1.0 - b.tokens
	return time.Duration(math.Ceil(missing / b.rate * float64(time.Second)))
}

func (b *bucket) refill(now time.Time) {
	elapsed := now.Sub(b.lastFill).Seconds()
	if elapsed <= 0 {
		return
	}
	b.tokens += elapsed * b.rate
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
	b.lastFill = now
}

func (b *bucket) commit(now time.Time) {
	b.refill(now)
	if b.tokens >= 1.0 {
		b.tokens -= 1.0
	}
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
	now := l.clock.Now()

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
			b = newBucket(math.Max(1, s.rps), s.rps, now)
			l.buckets[s.key] = b
		}
		out = append(out, b)
	}
	return out
}
