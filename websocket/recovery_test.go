package websocket

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	gorilla "github.com/gorilla/websocket"
	coinwerr "github.com/tigusigalpa/coinw-go/errors"
)

// scriptedConn is a goroutine-safe Conn: reads block until a scripted frame
// is available or the connection is closed.
type scriptedConn struct {
	mu        sync.Mutex
	sent      [][]byte
	onWrite   func([]byte)
	incoming  chan []byte
	closed    chan struct{}
	closeOnce sync.Once
}

func newScriptedConn(frames ...string) *scriptedConn {
	c := &scriptedConn{incoming: make(chan []byte, len(frames)+1), closed: make(chan struct{})}
	for _, f := range frames {
		c.incoming <- []byte(f)
	}
	return c
}

func (c *scriptedConn) WriteMessage(_ int, data []byte) error {
	c.mu.Lock()
	c.sent = append(c.sent, data)
	hook := c.onWrite
	c.mu.Unlock()
	if hook != nil {
		hook(data)
	}
	return nil
}

func (c *scriptedConn) ReadMessage() (int, []byte, error) {
	select {
	case m := <-c.incoming:
		return gorilla.TextMessage, m, nil
	case <-c.closed:
		return 0, nil, errors.New("use of closed connection")
	}
}

func (c *scriptedConn) SetReadDeadline(time.Time) error { return nil }

func (c *scriptedConn) Close() error {
	c.closeOnce.Do(func() { close(c.closed) })
	return nil
}

func (c *scriptedConn) sentFrames() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]string, len(c.sent))
	for i, f := range c.sent {
		out[i] = string(f)
	}
	return out
}

// scriptedDialer hands out a new connection per dial.
type scriptedDialer struct {
	mu    sync.Mutex
	next  func(n int) *scriptedConn
	dials int
}

func (d *scriptedDialer) DialContext(context.Context, string, http.Header) (Conn, *http.Response, error) {
	d.mu.Lock()
	d.dials++
	n := d.dials
	d.mu.Unlock()
	return d.next(n), nil, nil
}

func (d *scriptedDialer) count() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.dials
}

func waitFor(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatal("condition not reached")
		}
		time.Sleep(time.Millisecond)
	}
}

var (
	btcTicker = map[string]interface{}{"biz": "futures", "type": "ticker", "pairCode": "BTC"}
	ethTicker = map[string]interface{}{"biz": "futures", "type": "ticker", "pairCode": "ETH"}
)

func TestRejectedSubscriptionKeepsConnectionAndOtherStreams(t *testing.T) {
	d := &scriptedDialer{next: func(int) *scriptedConn {
		return newScriptedConn(
			`{"biz":"futures","channel":"subscribe","data":{"result":false},"type":"ticker","pairCode":"ETH"}`,
			`{"biz":"futures","type":"ticker","pairCode":"BTC","data":{}}`,
		)
	}}
	s := NewRawJsonSession("wss://fixture.test", false, "", "", WithDialer(d),
		WithReconnectPolicy(&ReconnectPolicy{BaseDelay: time.Millisecond}))
	if err := s.Subscribe(btcTicker); err != nil {
		t.Fatal(err)
	}
	if err := s.Subscribe(ethTicker); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var rejections []error
	events := 0
	err := s.Run(ctx, func(*Message) bool { events++; return false }, func(err error) bool {
		rejections = append(rejections, err)
		return true
	})
	var apiErr *coinwerr.APIError
	if err != nil || events != 1 || len(rejections) != 1 || !errors.As(rejections[0], &apiErr) {
		t.Fatalf("rejection was not isolated: err=%v events=%d rejections=%v", err, events, rejections)
	}
	if d.count() != 1 {
		t.Fatalf("a refused subscription forced %d connections", d.count())
	}
	s.mu.Lock()
	_, keptETH := s.subs[subKey(ethTicker)]
	_, keptBTC := s.subs[subKey(btcTicker)]
	s.mu.Unlock()
	if keptETH || !keptBTC {
		t.Fatalf("desired set after rejection: eth=%v btc=%v", keptETH, keptBTC)
	}
}

func TestRejectedConnectionsBackOffUntilMaxAttempts(t *testing.T) {
	d := &scriptedDialer{next: func(int) *scriptedConn {
		return newScriptedConn(`{"errorCode":29001,"errorMsg":"rate limited"}`)
	}}
	s := NewRawJsonSession("wss://fixture.test", false, "", "", WithDialer(d),
		WithReconnectPolicy(&ReconnectPolicy{BaseDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond, MaxAttempts: 3}))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := s.Run(ctx, func(*Message) bool { return true }, func(error) bool { return true })
	var exhausted *coinwerr.CancelledError
	if !errors.As(err, &exhausted) || d.count() != 4 {
		t.Fatalf("rejected connections did not exhaust the policy: err=%v dials=%d", err, d.count())
	}
}

func TestCloseStopsActiveRun(t *testing.T) {
	d := &scriptedDialer{next: func(int) *scriptedConn { return newScriptedConn() }}
	s := NewRawJsonSession("wss://fixture.test", false, "", "", WithDialer(d),
		WithReconnectPolicy(&ReconnectPolicy{BaseDelay: time.Millisecond}))
	if err := s.Connect(context.Background()); err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() {
		done <- s.Run(context.Background(), func(*Message) bool { return true }, func(error) bool { return true })
	}()
	waitFor(t, func() bool { s.mu.Lock(); defer s.mu.Unlock(); return s.stopRun != nil })
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run after Close returned %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Close did not stop Run")
	}
	if d.count() != 1 {
		t.Fatalf("Run reconnected after Close: %d dials", d.count())
	}
}

func TestSocketIOSubscribesOfflineAndCloseStopsRun(t *testing.T) {
	conn := newScriptedConn(`0{"sid":"engine"}`, `40{"sid":"namespace"}`)
	d := &scriptedDialer{next: func(int) *scriptedConn { return conn }}
	s := NewSocketIOSession("wss://fixture.test", func(context.Context) (string, string, error) { return "TOK", "", nil },
		WithSocketIODialer(d))
	if err := s.Subscribe("spot/match:BTC-USDT"); err != nil {
		t.Fatalf("subscribe while disconnected: %v", err)
	}
	done := make(chan error, 1)
	go func() {
		done <- s.Run(context.Background(), func(*SocketIOMessage) bool { return true }, func(error) bool { return true })
	}()
	waitFor(t, func() bool { return len(conn.sentFrames()) >= 2 })
	if sent := conn.sentFrames(); sent[0] != "40" || !strings.Contains(sent[1], "spot/match:BTC-USDT") {
		t.Fatalf("offline subscription not restored after handshake: %q", sent)
	}
	waitFor(t, func() bool { s.mu.Lock(); defer s.mu.Unlock(); return s.stopRun != nil })
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run after Close returned %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Close did not stop Socket.IO Run")
	}
	if d.count() != 1 {
		t.Fatalf("Socket.IO Run reconnected after Close: %d dials", d.count())
	}
}

func TestSubscriptionChangesWhileDisconnectedApplyOnConnect(t *testing.T) {
	conn := newScriptedConn()
	s := NewRawJsonSession("wss://fixture.test", false, "", "",
		WithDialer(&scriptedDialer{next: func(int) *scriptedConn { return conn }}))
	if err := s.Subscribe(btcTicker); err != nil {
		t.Fatalf("subscribe while disconnected: %v", err)
	}
	if err := s.Subscribe(ethTicker); err != nil {
		t.Fatal(err)
	}
	if err := s.Unsubscribe(ethTicker); err != nil {
		t.Fatalf("unsubscribe while disconnected: %v", err)
	}
	if err := s.Connect(context.Background()); err != nil {
		t.Fatal(err)
	}
	sent := conn.sentFrames()
	if len(sent) != 1 || !strings.Contains(sent[0], `"pairCode":"BTC"`) {
		t.Fatalf("connect did not apply offline changes: %q", sent)
	}
}

func TestConnectionDroppedDuringResubscribeDoesNotPanic(t *testing.T) {
	var s *RawJsonSession
	d := &scriptedDialer{next: func(n int) *scriptedConn {
		c := newScriptedConn(`{"biz":"futures","type":"ticker","pairCode":"BTC","data":{}}`)
		if n == 1 {
			// A cancelled SubscribeStream drops the connection mid-restore.
			c.onWrite = func([]byte) { s.dropQuietly() }
		}
		return c
	}}
	s = NewRawJsonSession("wss://fixture.test", false, "", "", WithDialer(d),
		WithReconnectPolicy(&ReconnectPolicy{BaseDelay: time.Millisecond}))
	if err := s.Subscribe(btcTicker); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	events := 0
	err := s.Run(ctx, func(*Message) bool { events++; return false }, func(error) bool { return true })
	if err != nil || events != 1 || d.count() != 2 {
		t.Fatalf("recovery after a dropped restore: err=%v events=%d dials=%d", err, events, d.count())
	}
}
