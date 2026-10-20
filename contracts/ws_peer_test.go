package contracts_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	ws "github.com/tigusigalpa/coinw-go/websocket"
)

type peer struct {
	mu     sync.Mutex
	frames [][]byte
	writes [][]byte
}

func (p *peer) WriteMessage(_ int, b []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.writes = append(p.writes, append([]byte(nil), b...))
	return nil
}
func (p *peer) ReadMessage() (int, []byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if len(p.frames) == 0 {
		return 0, nil, io.EOF
	}
	b := p.frames[0]
	p.frames = p.frames[1:]
	return 1, b, nil
}
func (*peer) SetReadDeadline(time.Time) error { return nil }
func (*peer) Close() error                    { return nil }

type peerDialer struct{ p *peer }

func (d peerDialer) DialContext(context.Context, string, http.Header) (ws.Conn, *http.Response, error) {
	return d.p, nil, nil
}

func checkStream[T any](t *testing.T, stream ws.Stream[T], fixture, ackFile string, socketIO, private bool) {
	t.Helper()
	event, err := os.ReadFile(fixture)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p := &peer{}
	if socketIO {
		p.frames = [][]byte{[]byte(`0{"sid":"fixture","pingInterval":25000,"pingTimeout":20000}`), []byte(`40`), append([]byte(`42["subscribe",`), append(event, ']')...)}
		s := ws.NewSocketIOSession("wss://fixture.test", func(context.Context) (string, string, error) { return "fixture-only", "", nil }, ws.WithSocketIODialer(peerDialer{p}))
		defer func() { _ = s.Close() }()
		if err = s.Connect(ctx); err != nil {
			t.Fatal(err)
		}
		if err = ws.SubscribeSocketIO(ctx, s, stream); err != nil {
			t.Fatal(err)
		}
		m, err := s.ReadMessage()
		if err != nil {
			t.Fatal(err)
		}
		if _, err = ws.DecodeSocketIO(m, stream); err != nil {
			t.Fatal(err)
		}
		if err = ws.UnsubscribeSocketIO(ctx, s, stream); err != nil {
			t.Fatal(err)
		}
	} else {
		ack, err := os.ReadFile(ackFile)
		if err != nil {
			t.Fatal(err)
		}
		if private {
			p.frames = append(p.frames, []byte(`{"channel":"login","data":{"result":true}}`))
		}
		p.frames = append(p.frames, ack, event)
		s := ws.NewRawJsonSession("wss://fixture.test", private, "fixture-key", "fixture-secret", ws.WithDialer(peerDialer{p}), ws.WithOnConnect(func(context.Context) error { return nil }))
		defer func() { _ = s.Close() }()
		if err = s.Connect(ctx); err != nil {
			t.Fatal(err)
		}
		if err = ws.SubscribeStream(ctx, s, stream); err != nil {
			t.Fatal(err)
		}
		m, err := s.ReadMessage()
		if err != nil {
			t.Fatal(err)
		}
		a, err := ws.DecodeAcknowledgement(m.RawBytes)
		if err != nil || !a.Data.Result {
			t.Fatalf("acknowledgement: %v", err)
		}
		m, err = s.ReadMessage()
		if err != nil {
			t.Fatal(err)
		}
		data, err := ws.DecodeStream(m, stream)
		if err != nil {
			t.Fatal(err)
		}
		if _, err = json.Marshal(data); err != nil {
			t.Fatal(err)
		}
		if err = ws.UnsubscribeStream(ctx, s, stream); err != nil {
			t.Fatal(err)
		}
		var sub, unsub struct {
			Event  string                 `json:"event"`
			Params map[string]interface{} `json:"params"`
		}
		offset := 0
		if private {
			offset = 1
		}
		if err = json.Unmarshal(p.writes[offset], &sub); err != nil {
			t.Fatal(err)
		}
		if err = json.Unmarshal(p.writes[offset+1], &unsub); err != nil {
			t.Fatal(err)
		}
		if sub.Event != "sub" || unsub.Event != "unsub" {
			t.Fatal("subscription lifecycle changed")
		}
		want, _ := json.Marshal(stream.Params())
		actual, _ := json.Marshal(sub.Params)
		if string(want) != string(actual) {
			t.Fatal("wire subscription differs from typed contract")
		}
	}
}
