package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	gorilla "github.com/gorilla/websocket"
)

func wsServer(t *testing.T, handler func(*gorilla.Conn)) string {
	t.Helper()
	upgrader := gorilla.Upgrader{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		handler(conn)
	}))
	t.Cleanup(srv.Close)
	return "ws" + strings.TrimPrefix(srv.URL, "http")
}

func TestRunCancellationInterruptsBlockedRead(t *testing.T) {
	url := wsServer(t, func(c *gorilla.Conn) {
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				return
			}
		}
	})
	s := NewRawJsonSession(url, false, "", "")
	if err := s.Connect(context.Background()); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- s.Run(ctx, func(*Message) bool { return true }, func(error) bool { return true }) }()
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected cancellation, got %v", err)
		}
	case <-time.After(time.Second):
		s.Close()
		t.Fatal("Run did not interrupt a blocked read")
	}
}

func TestPrivateLoginCancellationAndFailureLeaveNoConnection(t *testing.T) {
	url := wsServer(t, func(c *gorilla.Conn) {
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				return
			}
		}
	})
	s := NewRawJsonSession(url, true, "K", "S")
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
	defer cancel()
	err := s.Connect(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected login cancellation, got %v", err)
	}
	if s.conn != nil || s.IsAuthenticated() {
		t.Fatal("failed login left a connected session")
	}
	conn := &fakeConn{incoming: [][]byte{[]byte(`{"channel":"login","data":{"result":false}}`)}}
	s = NewRawJsonSession("wss://example.test", true, "K", "S", WithDialer(&fakeDialer{conn: conn}))
	if err := s.Connect(context.Background()); err == nil {
		t.Fatal("expected login rejection")
	}
	if s.conn != nil || !conn.closed {
		t.Fatal("rejected connection was retained")
	}
}

func TestHeartbeatRunsWhileReadIsBlocked(t *testing.T) {
	ping := make(chan struct{}, 1)
	url := wsServer(t, func(c *gorilla.Conn) {
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
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- s.Run(ctx, func(*Message) bool { return true }, func(error) bool { return false }) }()
	select {
	case <-ping:
	case <-time.After(time.Second):
		t.Error("heartbeat was blocked by ReadMessage")
	}
	cancel()
	<-done
}

func TestReadTimeoutReconnects(t *testing.T) {
	var connections atomic.Int32
	url := wsServer(t, func(c *gorilla.Conn) {
		if connections.Add(1) >= 2 {
			_ = c.WriteMessage(gorilla.TextMessage, []byte(`{"type":"ticker","data":{}}`))
		}
		for {
			if _, _, err := c.ReadMessage(); err != nil {
				return
			}
		}
	})
	s := NewRawJsonSession(url, false, "", "", WithTimeouts(time.Second, time.Second, 40*time.Millisecond), WithReconnectPolicy(&ReconnectPolicy{BaseDelay: time.Millisecond, MaxDelay: time.Millisecond}))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := s.Run(ctx, func(*Message) bool { return false }, func(error) bool { return true })
	if err != nil || connections.Load() < 2 {
		t.Fatalf("timeout did not reconnect: connections=%d err=%v", connections.Load(), err)
	}
}

func TestConcurrentWritesAndImmutableSubscription(t *testing.T) {
	var connections atomic.Int32
	resubscribed := make(chan []byte, 1)
	url := wsServer(t, func(c *gorilla.Conn) {
		n := connections.Add(1)
		for {
			_, raw, err := c.ReadMessage()
			if err != nil {
				return
			}
			if n == 2 && strings.Contains(string(raw), `"event":"sub"`) {
				resubscribed <- raw
				return
			}
		}
	})
	s := NewRawJsonSession(url, false, "", "", WithTimeouts(time.Second, time.Second, time.Second), WithReconnectPolicy(&ReconnectPolicy{BaseDelay: time.Millisecond}))
	if err := s.Connect(context.Background()); err != nil {
		t.Fatal(err)
	}
	params := map[string]interface{}{"type": "depth", "nested": map[string]interface{}{"pairCode": "BTC"}}
	if err := s.Subscribe(params); err != nil {
		t.Fatal(err)
	}
	params["nested"].(map[string]interface{})["pairCode"] = "CHANGED"
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				if err := s.Ping(); err != nil {
					t.Error(err)
					return
				}
			}
		}()
	}
	wg.Wait()
	_ = s.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- s.Run(ctx, func(*Message) bool { return true }, func(error) bool { return true }) }()
	select {
	case raw := <-resubscribed:
		if strings.Contains(string(raw), "CHANGED") || !strings.Contains(string(raw), "BTC") {
			t.Errorf("subscription changed: %s", raw)
		}
	case <-ctx.Done():
		t.Error("missing resubscription")
	}
	cancel()
	<-done
}

func TestRawDecoderPreservesNumbersAndRejectsNull(t *testing.T) {
	if _, err := decodeMessage([]byte("null")); err == nil {
		t.Fatal("null accepted as message")
	}
	m, err := decodeMessage([]byte(`{"type":"depth","seq":9007199254740993}`))
	if err != nil {
		t.Fatal(err)
	}
	if n, ok := m.Raw["seq"].(json.Number); !ok || n.String() != "9007199254740993" {
		t.Fatalf("sequence precision lost: %#v", m.Raw["seq"])
	}
}

func TestSocketIOHandshakeCompletesBeforeSubscribe(t *testing.T) {
	conn := &fakeConn{incoming: [][]byte{
		[]byte(`0{"sid":"engine"}`), []byte(`2probe`), []byte(`40{"sid":"namespace"}`),
	}}
	s := NewSocketIOSession("wss://example.test", func(context.Context) (string, string, error) { return "TOK", "", nil }, WithSocketIODialer(&fakeDialer{conn: conn}))
	if err := s.Connect(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := s.Subscribe("spot/match:BTC-USDT"); err != nil {
		t.Fatal(err)
	}
	if len(conn.sent) != 3 || string(conn.sent[0]) != "40" || string(conn.sent[1]) != "3probe" || !strings.HasPrefix(string(conn.sent[2]), "42") {
		t.Fatalf("bad handshake order: %q", conn.sent)
	}
}
