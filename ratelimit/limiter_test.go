package ratelimit

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/tigusigalpa/coinw-go/endpoints"
)

type fakeClock struct{ now time.Time }

func (f *fakeClock) Now() time.Time          { return f.now }
func (f *fakeClock) Advance(d time.Duration) { f.now = f.now.Add(d) }

func ep(group string, rps float64) endpoints.Definition {
	return endpoints.Definition{
		ID: "test.ep", Method: "GET", Path: "/v1/x",
		Auth: endpoints.AuthPublic, RateGroup: group,
		PerInterfaceRPS: rps, ReadOnly: true,
	}
}

func TestFractionalRateAllowsOneRequestAndRefills(t *testing.T) {
	c := &fakeClock{now: time.Unix(0, 0)}
	l := NewInMemoryLimiter(WithClock(c))
	bs := l.bucketsFor(ep(endpoints.GroupSpot, 0.5), "")
	b := bs[len(bs)-1]
	if w := b.waitDuration(c.Now()); w != 0 {
		t.Fatalf("first request must be available: %v", w)
	}
	b.commit(c.Now())
	if w := b.waitDuration(c.Now()); w != 2*time.Second {
		t.Fatalf("expected one token per 2s, got %v", w)
	}
	c.Advance(2 * time.Second)
	if w := b.waitDuration(c.Now()); w != 0 {
		t.Fatalf("fractional bucket did not refill: %v", w)
	}
}

func TestPerInterfaceIPBudgetIsSharedAcrossUIDs(t *testing.T) {
	c := &fakeClock{now: time.Unix(0, 0)}
	l := NewInMemoryLimiter(WithClock(c))
	e := ep(endpoints.GroupFuturesB, 1)
	if err := l.Acquire(context.Background(), e, "one"); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if err := l.Acquire(ctx, e, "two"); err == nil {
		t.Fatal("changing UID must not bypass the per-IP interface limit")
	}
}

func TestFuturesGlobalBudgetWithoutUID(t *testing.T) {
	l := NewInMemoryLimiter()
	if len(l.bucketsFor(ep(endpoints.GroupFuturesB, 0), "")) == 0 {
		t.Fatal("missing UID must not disable the group B limit")
	}
}

// limiter that waits on the fake clock instead of real time — used by
// overriding Acquire's sleep via a short-circuit: we measure calls, not
// real elapsed time.

func TestBucketWaitAndRefill(t *testing.T) {
	c := &fakeClock{now: time.Unix(0, 0)}
	b := newBucket(2)
	if w := b.waitDuration(c.Now()); w != 0 {
		t.Fatalf("expected no wait, got %v", w)
	}
	b.commit(c.Now())
	b.commit(c.Now())
	if w := b.waitDuration(c.Now()); w <= 0 {
		t.Fatal("expected positive wait when bucket empty")
	}
	c.Advance(time.Second)
	if w := b.waitDuration(c.Now()); w != 0 {
		t.Fatalf("expected refilled bucket, got wait %v", w)
	}
}

func TestBucketNeverExceedsRateWithinOneSecond(t *testing.T) {
	for _, rate := range []float64{100, 30, 5, 1.5} {
		c := &fakeClock{now: time.Unix(0, 0)}
		b := newBucket(rate)
		// A greedy caller takes every grant as soon as it is offered.
		var grants []time.Time
		for c.Now().Before(time.Unix(3, 0)) {
			if w := b.waitDuration(c.Now()); w > 0 {
				c.Advance(w)
				continue
			}
			b.commit(c.Now())
			grants = append(grants, c.Now())
		}
		limit := int(math.Ceil(rate))
		for i := range grants {
			inWindow := 0
			for _, g := range grants[i:] {
				if g.Sub(grants[i]) < time.Second {
					inWindow++
				}
			}
			if inWindow > limit {
				t.Fatalf("rate %v: %d grants within one second from %v", rate, inWindow, grants[i].Sub(time.Unix(0, 0)))
			}
		}
		if len(grants) < int(3*rate*0.9) {
			t.Fatalf("rate %v: only %d grants in 3s", rate, len(grants))
		}
	}
}

func TestAcquireRespectsGlobalSpotBucket(t *testing.T) {
	c := &fakeClock{now: time.Unix(0, 0)}
	l := NewInMemoryLimiter(WithClock(c))
	ctx := context.Background()
	e := ep(endpoints.GroupSpot, 0)

	// 100 capacity → exhaust without blocking (each acquire commits 1 token;
	// refill driven by fake clock stays at t=0)
	for i := 0; i < 100; i++ {
		if err := l.Acquire(ctx, e, ""); err != nil {
			t.Fatalf("acquire %d: %v", i, err)
		}
	}
	// Next acquire must wait ~10ms for a refill — but the real timer waits on
	// wall clock; use a cancelled context to verify the wait path exists.
	ctxCancel, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	if err := l.Acquire(ctxCancel, e, ""); err == nil {
		t.Fatal("expected context deadline while waiting for token")
	}
}

func TestAcquireHonoursCancelledContext(t *testing.T) {
	l := NewInMemoryLimiter()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := l.Acquire(ctx, ep(endpoints.GroupFuturesA, 0), "")
	if err == nil {
		t.Fatal("expected ctx error on cancelled acquire")
	}
}
