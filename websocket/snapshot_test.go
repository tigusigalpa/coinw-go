package websocket

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	gorilla "github.com/gorilla/websocket"
	coinwerr "github.com/tigusigalpa/coinw-go/errors"
)

func TestPrivateReconnectReloginsSnapshotsAndResubscribes(t *testing.T) {
	var connections, snapshots, subscriptions atomic.Int32
	url := wsServer(t, func(c *gorilla.Conn) {
		n := connections.Add(1)
		_, login, err := c.ReadMessage()
		if err != nil {
			return
		}
		if len(login) == 0 {
			t.Error("missing login")
			return
		}
		if err = c.WriteMessage(gorilla.TextMessage, []byte(`{"channel":"login","data":{"result":true}}`)); err != nil {
			return
		}
		_, _, err = c.ReadMessage()
		if err != nil {
			return
		}
		subscriptions.Add(1)
		if snapshots.Load() < n {
			t.Error("subscription restored before the required snapshot hook")
		}
		if n == 1 {
			return
		}
		_ = c.WriteMessage(gorilla.TextMessage, []byte(`{"biz":"futures","type":"assets","data":[]}`))
		for {
			if _, _, err = c.ReadMessage(); err != nil {
				return
			}
		}
	})
	s := NewRawJsonSession(url, true, "fixture-key", "fixture-secret", WithOnConnect(func(ctx context.Context) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		snapshots.Add(1)
		return nil
	}), WithReconnectPolicy(&ReconnectPolicy{BaseDelay: time.Millisecond, MaxDelay: time.Millisecond}))
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := s.Connect(ctx); err != nil {
		t.Fatal(err)
	}
	if err := SubscribeStream(ctx, s, FuturesAssets()); err != nil {
		t.Fatal(err)
	}
	if err := s.Run(ctx, func(*Message) bool { return false }, func(error) bool { return true }); err != nil {
		t.Fatal(err)
	}
	if connections.Load() != 2 || snapshots.Load() != 2 || subscriptions.Load() != 2 {
		t.Fatalf("incomplete private recovery: connections=%d snapshots=%d subscriptions=%d", connections.Load(), snapshots.Load(), subscriptions.Load())
	}
}

func TestDepthRequestsSnapshotBeforeApplyingAndOnGap(t *testing.T) {
	var tracker DepthSequenceTracker
	start, end := int64(10), int64(12)
	requests := 0
	request := func(context.Context) error { requests++; return nil }
	result, err := tracker.TrackAndResync(context.Background(), &start, &end, request)
	if err != nil || result != DepthNeedsSnapshot || tracker.ExpectedNext() != 0 || requests != 1 {
		t.Fatal("unanchored increment was applied")
	}
	tracker.Reset(9)
	if tracker.Track(&start, &end) != DepthApplied {
		t.Fatal("anchored increment rejected")
	}
	start, end = 20, 21
	result, err = tracker.TrackAndResync(context.Background(), &start, &end, request)
	if err != nil || result != DepthGap || requests != 2 || tracker.ExpectedNext() != 13 {
		t.Fatal("gap did not request snapshot without advancing")
	}
	start, end = 1, 8
	if tracker.Track(&start, &end) != DepthDuplicate {
		t.Fatal("out-of-order old update was applied")
	}
}

func TestMalformedAcknowledgementsAndProfileSeparation(t *testing.T) {
	for _, raw := range []string{`{"channel":"login","data":{}}`, `{"channel":"subscribe","data":{"result":"true"}}`, `{"channel":"login","data":null}`, `{"channel":"login","data":{"result":true,"private":"value"}}`} {
		if _, err := decodeMessage([]byte(raw)); err == nil {
			t.Fatal("malformed acknowledgement accepted")
		}
	}
	s := NewRawJsonSession("wss://fixture.test", false, "", "", WithProfile(SpotProfile))
	if err := SubscribeStream(context.Background(), s, FuturesAssets()); err == nil {
		t.Fatal("private Futures subscription accepted on public Spot session")
	}
}

func TestBackpressureKeepsHeartbeatAlive(t *testing.T) {
	ping := make(chan struct{}, 1)
	entered := make(chan struct{})
	release := make(chan struct{})
	url := wsServer(t, func(c *gorilla.Conn) {
		if err := c.WriteMessage(gorilla.TextMessage, []byte(`{"type":"ticker","data":{}}`)); err != nil {
			return
		}
		for {
			_, raw, err := c.ReadMessage()
			if err != nil {
				return
			}
			if string(raw) == `{"event":"ping"}` {
				select {
				case ping <- struct{}{}:
				default:
				}
			}
		}
	})
	s := NewRawJsonSession(url, false, "", "", WithTimeouts(20*time.Millisecond, time.Second, time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- s.Run(ctx, func(*Message) bool { close(entered); <-release; return false }, func(error) bool { return false })
	}()
	select {
	case <-entered:
	case <-ctx.Done():
		close(release)
		t.Fatal("event callback did not start")
	}
	select {
	case <-ping:
	case <-time.After(time.Second):
		t.Error("backpressure blocked the heartbeat")
	}
	close(release)
	if err := <-done; err != nil {
		t.Fatal(err)
	}
}

func TestServerRejectionUsesErrorCallback(t *testing.T) {
	url := wsServer(t, func(c *gorilla.Conn) {
		_ = c.WriteMessage(gorilla.TextMessage, []byte(`{"errorCode":29001,"errorMsg":"rate limited"}`))
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				return
			}
		}
	})
	s := NewRawJsonSession(url, false, "", "")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	events, failures := 0, 0
	err := s.Run(ctx, func(*Message) bool { events++; return false }, func(error) bool { failures++; return false })
	var limit *coinwerr.RateLimitError
	if !errors.As(err, &limit) || limit.Code != "29001" || events != 0 || failures != 1 {
		t.Fatalf("rejected operation bypassed error delivery: %v", err)
	}
}
