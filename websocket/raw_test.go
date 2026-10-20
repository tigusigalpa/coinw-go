package websocket

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	coinwerr "github.com/tigusigalpa/coinw-go/errors"
)

// fakeConn is a scripted Conn: records sent frames, replays incoming frames.
type fakeConn struct {
	sent     [][]byte
	incoming [][]byte
	closed   bool
}

func (f *fakeConn) WriteMessage(_ int, data []byte) error {
	f.sent = append(f.sent, data)
	return nil
}

func (f *fakeConn) ReadMessage() (int, []byte, error) {
	if len(f.incoming) == 0 {
		return 0, nil, &timeoutErr{}
	}
	m := f.incoming[0]
	f.incoming = f.incoming[1:]
	return 1, m, nil
}

func (f *fakeConn) SetReadDeadline(time.Time) error { return nil }
func (f *fakeConn) Close() error                    { f.closed = true; return nil }

type timeoutErr struct{}

func (timeoutErr) Error() string   { return "i/o timeout" }
func (timeoutErr) Timeout() bool   { return true }
func (timeoutErr) Temporary() bool { return true }

// fakeDialer returns the same fakeConn.
type fakeDialer struct {
	conn *fakeConn
	urls []string
}

func (d *fakeDialer) DialContext(_ context.Context, urlStr string, _ http.Header) (Conn, *http.Response, error) {
	d.urls = append(d.urls, urlStr)
	return d.conn, nil, nil
}

func TestRawJsonPublicSession(t *testing.T) {
	conn := &fakeConn{incoming: [][]byte{
		[]byte(`{"biz":"exchange","channel":"subscribe","data":{"result":true},"type":"depth"}`),
		[]byte(`{"event":"pong"}`),
		[]byte(`{"biz":"exchange","type":"depth","data":{"startSeq":1,"endSeq":2}}`),
	}}
	d := &fakeDialer{conn: conn}
	s := NewRawJsonSession("wss://ws.futurescw.com", false, "", "", WithDialer(d))

	if err := s.Connect(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := s.Subscribe(map[string]interface{}{"biz": "exchange", "type": "depth", "pairCode": "78"}); err != nil {
		t.Fatal(err)
	}
	last := string(conn.sent[len(conn.sent)-1])
	if !strings.Contains(last, `"event":"sub"`) || !strings.Contains(last, `"pairCode":"78"`) {
		t.Fatalf("bad sub frame: %s", last)
	}
	if err := s.Ping(); err != nil {
		t.Fatal(err)
	}
	if string(conn.sent[len(conn.sent)-1]) != `{"event":"ping"}` {
		t.Fatal("bad ping frame")
	}
	m, _ := s.ReadMessage()
	if !m.IsSubscribeAck() || !m.AckSucceeded() {
		t.Fatal("expected subscribe ack")
	}
	m, _ = s.ReadMessage()
	if !m.IsPong() {
		t.Fatal("expected pong")
	}
	m, _ = s.ReadMessage()
	if !m.IsEvent() || m.Type != "depth" {
		t.Fatal("expected depth event")
	}
}

func TestRawJsonPrivateLogin(t *testing.T) {
	conn := &fakeConn{incoming: [][]byte{
		[]byte(`{"channel":"login","data":{"result":true}}`),
	}}
	s := NewRawJsonSession("wss://ws.futurescw.com/perpum", true, "K", "S",
		WithDialer(&fakeDialer{conn: conn}))
	if err := s.Connect(context.Background()); err != nil {
		t.Fatal(err)
	}
	frame := string(conn.sent[0])
	if !strings.Contains(frame, `"event":"login"`) ||
		!strings.Contains(frame, `"passphrase":"S"`) ||
		!strings.Contains(frame, `"api_key":"K"`) {
		t.Fatalf("bad login frame: %s", frame)
	}
	if !s.IsAuthenticated() {
		t.Fatal("expected authenticated session")
	}
}

func TestRawJsonPrivateLoginRejected(t *testing.T) {
	conn := &fakeConn{incoming: [][]byte{
		[]byte(`{"channel":"login","data":{"result":false}}`),
	}}
	s := NewRawJsonSession("wss://x", true, "K", "S", WithDialer(&fakeDialer{conn: conn}))
	err := s.Connect(context.Background())
	var ae *coinwerr.AuthenticationError
	if !errors.As(err, &ae) {
		t.Fatalf("expected AuthenticationError, got %v", err)
	}
}

func TestSocketIOSession(t *testing.T) {
	conn := &fakeConn{incoming: [][]byte{
		[]byte(`0{"sid":"abc"}`),
		[]byte(`40{"sid":"namespace"}`),
		[]byte(`2`),
		[]byte(`42["subscribe",{"args":"spot/match:BTC-USDT","success":true}]`),
	}}
	d := &fakeDialer{conn: conn}
	s := NewSocketIOSession("wss://ws.futurescw.info",
		func(context.Context) (string, string, error) { return "TOK", "", nil },
		WithSocketIODialer(d))

	if err := s.Connect(context.Background()); err != nil {
		t.Fatal(err)
	}
	u, _ := url.Parse(d.urls[0])
	if u.Path != "/socket.io/" || u.Query().Get("token") != "TOK" || u.Query().Get("EIO") != "4" || u.Query().Get("transport") != "websocket" {
		t.Fatalf("bad url: %s", d.urls[0])
	}
	if err := s.Subscribe("spot/match:BTC-USDT"); err != nil {
		t.Fatal(err)
	}
	last := string(conn.sent[len(conn.sent)-1])
	if last != `42["subscribe",{"args":"spot/match:BTC-USDT"}]` {
		t.Fatalf("bad subscribe frame: %s", last)
	}
	// Connect completed '0' → '40' before sending subscriptions.
	found40 := false
	for _, f := range conn.sent {
		if string(f) == "40" {
			found40 = true
		}
	}
	if !found40 {
		t.Fatal("expected client to send Socket.IO connect '40'")
	}
	// '2' engine ping → "3" pong
	if _, err := s.ReadMessage(); err != nil {
		t.Fatal(err)
	}
	if string(conn.sent[len(conn.sent)-1]) != "3" {
		t.Fatal("expected engine pong '3'")
	}
	// event frame
	m, err := s.ReadMessage()
	if err != nil || m == nil {
		t.Fatal(err)
	}
	if m.Event != "subscribe" || !strings.Contains(string(m.Data), "spot/match:BTC-USDT") {
		t.Fatalf("bad event: %+v", m)
	}
}

func TestDepthSequenceTracker(t *testing.T) {
	tr := &DepthSequenceTracker{}
	i64 := func(v int64) *int64 { return &v }

	if r := tr.Track(i64(100), i64(110)); r != DepthNeedsSnapshot {
		t.Fatalf("unanchored event must request a snapshot: %v", r)
	}
	tr.Reset(110)
	if r := tr.Track(i64(111), i64(120)); r != DepthApplied {
		t.Fatalf("sequential apply: %v", r)
	}
	if r := tr.Track(i64(95), i64(105)); r != DepthDuplicate {
		t.Fatalf("old event → duplicate: %v", r)
	}
	if r := tr.Track(i64(200), i64(210)); r != DepthGap {
		t.Fatalf("hole → gap: %v", r)
	}
	if tr.ExpectedNext() != 121 {
		t.Fatalf("gap must not advance expectation, got %d", tr.ExpectedNext())
	}
	tr.Reset(250)
	if r := tr.Track(i64(251), i64(260)); r != DepthApplied {
		t.Fatalf("post-resync apply: %v", r)
	}
	if r := tr.Track(nil, i64(1)); r != DepthNoSequence {
		t.Fatalf("missing seq: %v", r)
	}
}
