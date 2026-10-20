// Package websocket implements the documented CoinW stream adapters:
//   - raw JSON protocol (Spot public Method 2, Spot private, Futures
//     public/private on /perpum)
//   - Socket.IO/Engine.IO protocol (Spot public Method 1)
//   - DepthSequenceTracker for the incremental order book
package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	coinwerr "github.com/tigusigalpa/coinw-go/errors"
)

// Dialer abstracts the WS dial for tests. *websocket.Dialer satisfies it.
type Dialer interface {
	DialContext(ctx context.Context, urlStr string, requestHeader http.Header) (Conn, *http.Response, error)
}

// Conn is the subset of *websocket.Conn the sessions use.
type Conn interface {
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, p []byte, err error)
	SetReadDeadline(t time.Time) error
	Close() error
}

type gorillaDialer struct{ d *websocket.Dialer }

func (g gorillaDialer) DialContext(ctx context.Context, urlStr string, h http.Header) (Conn, *http.Response, error) {
	return g.d.DialContext(ctx, urlStr, h)
}

// DefaultDialer returns a sane production dialer.
func DefaultDialer() Dialer {
	return gorillaDialer{d: &websocket.Dialer{
		HandshakeTimeout: 15 * time.Second,
	}}
}

// Message is a decoded raw-JSON frame.
type Message struct {
	Channel   string                 `json:"channel"`
	Type      string                 `json:"type"`
	Biz       string                 `json:"biz"`
	Event     string                 `json:"event"`
	Data      json.RawMessage        `json:"data"`
	ErrorCode int                    `json:"errorCode"`
	ErrorMsg  string                 `json:"errorMsg"`
	Raw       map[string]interface{} `json:"-"`
	RawBytes  []byte                 `json:"-"`
}

// IsPong / IsLoginAck / IsSubscribeAck / IsEvent / IsError classify frames
// per the documented shapes.
func (m *Message) IsPong() bool     { return m.Event == "pong" }
func (m *Message) IsLoginAck() bool { return m.Channel == "login" }
func (m *Message) IsSubscribeAck() bool {
	return m.Channel == "subscribe" || m.Channel == "unsubscribe"
}
func (m *Message) IsEvent() bool {
	return m.Channel == "" && m.Type != "" && m.Event == ""
}
func (m *Message) IsError() bool {
	if m.ErrorCode != 0 {
		return true
	}
	if r, ok := m.Raw["result"].(bool); ok && !r {
		return true
	}
	var data struct {
		Result *bool `json:"result"`
	}
	if len(m.Data) > 0 && json.Unmarshal(m.Data, &data) == nil && data.Result != nil && !*data.Result {
		return true
	}
	return false
}

// AckSucceeded reports result:true in the data payload (login/sub acks).
func (m *Message) AckSucceeded() bool {
	var data struct {
		Result *bool `json:"result"`
	}
	if len(m.Data) == 0 || json.Unmarshal(m.Data, &data) != nil {
		return false
	}
	return data.Result != nil && *data.Result
}

func decodeMessage(raw []byte) (*Message, error) {
	var m Message
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, &coinwerr.ProtocolError{Message: "malformed WS frame (not a JSON object)"}
	}
	var rawMap map[string]interface{}
	_ = json.Unmarshal(raw, &rawMap)
	m.Raw = rawMap
	m.RawBytes = raw
	return &m, nil
}

func frameJSON(event string, params map[string]interface{}) ([]byte, error) {
	return json.Marshal(map[string]interface{}{"event": event, "params": params})
}

// subscription is a desired sub kept for automatic resubscribe.
type subscription struct {
	key    string
	params map[string]interface{}
}

// ReconnectPolicy — bounded exponential backoff.
type ReconnectPolicy struct {
	BaseDelay   time.Duration // default 500ms
	MaxDelay    time.Duration // default 15s
	MaxAttempts int           // 0 = unlimited until ctx cancelled
	attempts    int
}

func (p *ReconnectPolicy) defaults() {
	if p.BaseDelay <= 0 {
		p.BaseDelay = 500 * time.Millisecond
	}
	if p.MaxDelay <= 0 {
		p.MaxDelay = 15 * time.Second
	}
}

func (p *ReconnectPolicy) wait(ctx context.Context) error {
	p.defaults()
	if p.MaxAttempts > 0 && p.attempts >= p.MaxAttempts {
		return &coinwerr.CancelledError{Message: "max reconnect attempts reached"}
	}
	p.attempts++
	d := p.BaseDelay << min(p.attempts-1, 10)
	if d > p.MaxDelay {
		d = p.MaxDelay
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

func (p *ReconnectPolicy) reset() { p.attempts = 0 }

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ---------------------------------------------------------------------------
// RawJsonSession — Spot Method 2 / Spot private / Futures public+private
// ---------------------------------------------------------------------------

// RawJsonSession implements the documented raw JSON protocol:
//
//	{"event":"login","params":{"api_key":…,"passphrase":…}}
//	{"event":"sub","params":{"biz":…,"type":…,"pairCode":…,"interval":…}}
//	{"event":"unsub","params":{…}}   {"event":"ping"} → {"event":"pong"}
//
// Private pushes do not include an initial snapshot: after connect and after
// every reconnect pull the matching REST snapshot before applying events.
type RawJsonSession struct {
	url         string
	private     bool
	apiKey      string
	secretKey   string
	dialer      Dialer
	reconnect   *ReconnectPolicy
	pingPeriod  time.Duration
	ackTimeout  time.Duration
	readTimeout time.Duration

	mu   sync.Mutex
	conn Conn
	subs map[string]subscription
	auth bool
}

// RawOption customises a session.
type RawOption func(*RawJsonSession)

// WithReconnectPolicy sets a custom reconnect policy.
func WithReconnectPolicy(p *ReconnectPolicy) RawOption {
	return func(s *RawJsonSession) { s.reconnect = p }
}

// WithDialer injects a custom dialer (tests/mocks).
func WithDialer(d Dialer) RawOption {
	return func(s *RawJsonSession) { s.dialer = d }
}

// WithTimeouts tunes ping/ack/read timings.
func WithTimeouts(ping, ack, read time.Duration) RawOption {
	return func(s *RawJsonSession) {
		if ping > 0 {
			s.pingPeriod = ping
		}
		if ack > 0 {
			s.ackTimeout = ack
		}
		if read > 0 {
			s.readTimeout = read
		}
	}
}

// NewRawJsonSession creates a session. For private sessions pass apiKey and
// secretKey (documented login uses passphrase = secret key).
func NewRawJsonSession(url string, private bool, apiKey, secretKey string, opts ...RawOption) *RawJsonSession {
	s := &RawJsonSession{
		url:         url,
		private:     private,
		apiKey:      apiKey,
		secretKey:   secretKey,
		dialer:      DefaultDialer(),
		reconnect:   &ReconnectPolicy{},
		pingPeriod:  10 * time.Second,
		ackTimeout:  10 * time.Second,
		readTimeout: 30 * time.Second,
		subs:        map[string]subscription{},
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// IsAuthenticated reports whether the documented login succeeded.
func (s *RawJsonSession) IsAuthenticated() bool { s.mu.Lock(); defer s.mu.Unlock(); return s.auth }

// Connect opens the connection; performs the documented login for private
// sessions (blocks until the login ack or ackTimeout).
func (s *RawJsonSession) Connect(ctx context.Context) error {
	conn, resp, err := s.dialer.DialContext(ctx, s.url, nil)
	if err != nil {
		status := ""
		if resp != nil {
			status = " (HTTP " + resp.Status + ")"
		}
		return &coinwerr.TransportError{Message: "websocket dial failed" + status, Err: err}
	}
	s.mu.Lock()
	s.conn = conn
	s.auth = false
	s.mu.Unlock()

	if s.private {
		if s.apiKey == "" || s.secretKey == "" {
			conn.Close()
			return &coinwerr.AuthenticationError{Message: "private WS session requires api key and secret key"}
		}
		frame, _ := frameJSON("login", map[string]interface{}{
			"api_key": s.apiKey, "passphrase": s.secretKey,
		})
		if err := conn.WriteMessage(websocket.TextMessage, frame); err != nil {
			conn.Close()
			return &coinwerr.TransportError{Message: "failed to send login frame", Err: err}
		}
		if err := s.awaitLoginAck(ctx); err != nil {
			conn.Close()
			return err
		}
		s.mu.Lock()
		s.auth = true
		s.mu.Unlock()
	}
	return nil
}

func (s *RawJsonSession) awaitLoginAck(ctx context.Context) error {
	deadline := time.Now().Add(s.ackTimeout)
	for time.Now().Before(deadline) {
		if err := ctx.Err(); err != nil {
			return err
		}
		_ = s.conn.SetReadDeadline(deadline)
		_, raw, err := s.conn.ReadMessage()
		if err != nil {
			return &coinwerr.TransportError{Message: "failed reading login ack", Err: err}
		}
		m, err := decodeMessage(raw)
		if err != nil {
			continue
		}
		if m.IsLoginAck() {
			if m.AckSucceeded() {
				return nil
			}
			return &coinwerr.AuthenticationError{Message: "WS login rejected by server"}
		}
		if m.IsError() && m.ErrorCode != 0 {
			return &coinwerr.AuthenticationError{
				Message: fmt.Sprintf("WS login error %d: %s", m.ErrorCode, m.ErrorMsg),
			}
		}
	}
	return &coinwerr.TransportError{Message: "timed out waiting for WS login ack"}
}

// Subscribe registers and sends a sub frame. Params per docs:
// biz + type (+ pairCode / interval when the channel requires them).
func (s *RawJsonSession) Subscribe(params map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn == nil {
		return &coinwerr.TransportError{Message: "websocket is not connected — call Connect first"}
	}
	frame, err := frameJSON("sub", params)
	if err != nil {
		return err
	}
	if err := s.conn.WriteMessage(websocket.TextMessage, frame); err != nil {
		return &coinwerr.TransportError{Message: "subscribe send failed", Err: err}
	}
	s.subs[subKey(params)] = subscription{key: subKey(params), params: params}
	return nil
}

// Unsubscribe removes the registration and sends an unsub frame.
func (s *RawJsonSession) Unsubscribe(params map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn == nil {
		return &coinwerr.TransportError{Message: "websocket is not connected"}
	}
	frame, err := frameJSON("unsub", params)
	if err != nil {
		return err
	}
	if err := s.conn.WriteMessage(websocket.TextMessage, frame); err != nil {
		return &coinwerr.TransportError{Message: "unsubscribe send failed", Err: err}
	}
	delete(s.subs, subKey(params))
	return nil
}

// Ping sends the documented {"event":"ping"} keepalive.
func (s *RawJsonSession) Ping() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn == nil {
		return &coinwerr.TransportError{Message: "websocket is not connected"}
	}
	return s.conn.WriteMessage(websocket.TextMessage, []byte(`{"event":"ping"}`))
}

// ReadMessage reads and decodes the next frame. nil,nil on read timeout.
func (s *RawJsonSession) ReadMessage() (*Message, error) {
	s.mu.Lock()
	conn := s.conn
	s.mu.Unlock()
	if conn == nil {
		return nil, &coinwerr.TransportError{Message: "websocket is not connected"}
	}
	_ = conn.SetReadDeadline(time.Now().Add(s.readTimeout))
	mt, raw, err := conn.ReadMessage()
	if err != nil {
		if isTimeout(err) {
			return nil, nil
		}
		return nil, &coinwerr.TransportError{Message: "websocket read failed", Err: err}
	}
	if mt != websocket.TextMessage && mt != websocket.BinaryMessage {
		return nil, nil
	}
	return decodeMessage(raw)
}

// Run is the blocking event loop: ping keepalive, bounded-backoff
// reconnect, re-login and automatic resubscribe of the registered set.
// onEvent receives every decoded frame except pong. Run returns when ctx is
// cancelled or a fatal error escapes onError (return false to stop).
func (s *RawJsonSession) Run(ctx context.Context, onEvent func(*Message) bool, onError func(error) bool) error {
	lastPing := time.Time{}
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		s.mu.Lock()
		connected := s.conn != nil
		s.mu.Unlock()

		if !connected {
			if err := s.Connect(ctx); err != nil {
				if !onError(err) {
					return err
				}
				if err := s.reconnect.wait(ctx); err != nil {
					return err
				}
				continue
			}
			s.reconnect.reset()
			lastPing = time.Time{}
			// resubscribe desired set
			s.mu.Lock()
			subs := make([]subscription, 0, len(s.subs))
			for _, sub := range s.subs {
				subs = append(subs, sub)
			}
			conn := s.conn
			s.mu.Unlock()
			for _, sub := range subs {
				frame, _ := frameJSON("sub", sub.params)
				if err := conn.WriteMessage(websocket.TextMessage, frame); err != nil {
					break
				}
			}
		}

		if time.Since(lastPing) >= s.pingPeriod {
			if err := s.Ping(); err != nil {
				if !onError(err) {
					return err
				}
				s.closeQuietly()
				continue
			}
			lastPing = time.Now()
		}

		msg, err := s.ReadMessage()
		if err != nil {
			if !onError(err) {
				return err
			}
			s.closeQuietly()
			if err := s.reconnect.wait(ctx); err != nil {
				return err
			}
			continue
		}
		if msg == nil || msg.IsPong() {
			continue
		}
		if !onEvent(msg) {
			return nil
		}
	}
}

// Close closes the underlying connection.
func (s *RawJsonSession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn != nil {
		err := s.conn.Close()
		s.conn = nil
		s.auth = false
		return err
	}
	return nil
}

func (s *RawJsonSession) closeQuietly() { _ = s.Close() }

func subKey(params map[string]interface{}) string {
	raw, _ := json.Marshal(params)
	return string(raw)
}

func isTimeout(err error) bool {
	if ne, ok := err.(interface{ Timeout() bool }); ok {
		return ne.Timeout()
	}
	return false
}

// ---------------------------------------------------------------------------
// SocketIOSession — Spot public Method 1
// ---------------------------------------------------------------------------

// SocketIOMessage is a decoded Socket.IO event: event name + data payload.
type SocketIOMessage struct {
	Event string
	Data  json.RawMessage
}

// PublicTokenSource fetches the documented public WS token.
type PublicTokenSource func(ctx context.Context) (token, endpoint string, err error)

// SocketIOSession implements Spot public Method 1:
//
//  1. GET https://www.coinw.com/pusher/public-token → data.token (+endpoint)
//  2. wss://{host}?token={token}
//  3. 42["subscribe",{"args":"spot/match:BTC-USDT"}] etc.
//  4. keepalive 42["event",{"event":"ping"}] every ~10s
//  5. Engine.IO ping "2" → pong "3" ("2probe" → "3probe")
//
// Public-only per the documentation — no login.
type SocketIOSession struct {
	host        string
	tokenSource PublicTokenSource
	dialer      Dialer
	reconnect   *ReconnectPolicy
	pingPeriod  time.Duration
	readTimeout time.Duration

	mu   sync.Mutex
	conn Conn
	args map[string]struct{}
}

// SocketIOOption customises a session.
type SocketIOOption func(*SocketIOSession)

func WithSocketIODialer(d Dialer) SocketIOOption { return func(s *SocketIOSession) { s.dialer = d } }
func WithSocketIOReconnect(p *ReconnectPolicy) SocketIOOption {
	return func(s *SocketIOSession) { s.reconnect = p }
}
func WithSocketIOTimeouts(ping, read time.Duration) SocketIOOption {
	return func(s *SocketIOSession) {
		if ping > 0 {
			s.pingPeriod = ping
		}
		if read > 0 {
			s.readTimeout = read
		}
	}
}

// NewSocketIOSession creates a Method-1 session. host is the default
// wss URL (e.g. wss://ws.futurescw.info); the token endpoint may override it.
func NewSocketIOSession(host string, tokenSource PublicTokenSource, opts ...SocketIOOption) *SocketIOSession {
	s := &SocketIOSession{
		host:        host,
		tokenSource: tokenSource,
		dialer:      DefaultDialer(),
		reconnect:   &ReconnectPolicy{},
		pingPeriod:  10 * time.Second,
		readTimeout: 30 * time.Second,
		args:        map[string]struct{}{},
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// Connect fetches the public token and dials.
func (s *SocketIOSession) Connect(ctx context.Context) error {
	if s.tokenSource == nil {
		return &coinwerr.ValidationError{Message: "a PublicTokenSource is required"}
	}
	token, endpoint, err := s.tokenSource(ctx)
	if err != nil {
		return err
	}
	if token == "" {
		return &coinwerr.ProtocolError{Message: "public WS token response missing data.token"}
	}
	host := s.host
	if endpoint != "" {
		if hasWSScheme(endpoint) {
			host = endpoint
		} else {
			host = "wss://" + endpoint
		}
	}
	u, err := url.Parse(host)
	if err != nil {
		return &coinwerr.ValidationError{Message: "invalid Socket.IO ws url: " + err.Error()}
	}
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()

	conn, resp, err := s.dialer.DialContext(ctx, u.String(), nil)
	if err != nil {
		status := ""
		if resp != nil {
			status = " (HTTP " + resp.Status + ")"
		}
		return &coinwerr.TransportError{Message: "websocket dial failed" + status, Err: err}
	}
	s.mu.Lock()
	s.conn = conn
	s.mu.Unlock()
	return nil
}

func hasWSScheme(host string) bool {
	return len(host) > 6 && (host[:6] == "wss://" || host[:5] == "ws://")
}

// Subscribe sends 42["subscribe",{"args":"…"}] and registers for resubscribe.
func (s *SocketIOSession) Subscribe(args string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn == nil {
		return &coinwerr.TransportError{Message: "websocket is not connected"}
	}
	frame, err := sioEvent("subscribe", map[string]interface{}{"args": args})
	if err != nil {
		return err
	}
	if err := s.conn.WriteMessage(websocket.TextMessage, frame); err != nil {
		return &coinwerr.TransportError{Message: "subscribe send failed", Err: err}
	}
	s.args[args] = struct{}{}
	return nil
}

// Unsubscribe sends 42["unsubscribe",{"args":"…"}].
func (s *SocketIOSession) Unsubscribe(args string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn == nil {
		return &coinwerr.TransportError{Message: "websocket is not connected"}
	}
	frame, err := sioEvent("unsubscribe", map[string]interface{}{"args": args})
	if err != nil {
		return err
	}
	if err := s.conn.WriteMessage(websocket.TextMessage, frame); err != nil {
		return &coinwerr.TransportError{Message: "unsubscribe send failed", Err: err}
	}
	delete(s.args, args)
	return nil
}

// Ping emits the documented keepalive event.
func (s *SocketIOSession) Ping() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn == nil {
		return &coinwerr.TransportError{Message: "websocket is not connected"}
	}
	frame, err := sioEvent("event", map[string]interface{}{"event": "ping"})
	if err != nil {
		return err
	}
	return s.conn.WriteMessage(websocket.TextMessage, frame)
}

// ReadMessage returns the next Socket.IO event, or nil,nil for engine-level
// frames/timeouts. Engine.IO pings are answered automatically.
func (s *SocketIOSession) ReadMessage() (*SocketIOMessage, error) {
	s.mu.Lock()
	conn := s.conn
	s.mu.Unlock()
	if conn == nil {
		return nil, &coinwerr.TransportError{Message: "websocket is not connected"}
	}
	_ = conn.SetReadDeadline(time.Now().Add(s.readTimeout))
	mt, raw, err := conn.ReadMessage()
	if err != nil {
		if isTimeout(err) {
			return nil, nil
		}
		return nil, &coinwerr.TransportError{Message: "websocket read failed", Err: err}
	}
	if mt != websocket.TextMessage || len(raw) == 0 {
		return nil, nil
	}
	return s.handlePacket(conn, raw)
}

func (s *SocketIOSession) handlePacket(conn Conn, raw []byte) (*SocketIOMessage, error) {
	switch raw[0] {
	case '0': // Engine.IO open → request Socket.IO connect
		if err := conn.WriteMessage(websocket.TextMessage, []byte("40")); err != nil {
			return nil, &coinwerr.TransportError{Message: "failed to send Socket.IO connect", Err: err}
		}
		return nil, nil
	case '1':
		return nil, &coinwerr.TransportError{Message: "Engine.IO close packet received"}
	case '2': // Engine.IO ping → pong (preserve "probe" suffix)
		pong := []byte{'3'}
		if len(raw) > 1 {
			pong = append(pong, raw[1:]...)
		}
		if err := conn.WriteMessage(websocket.TextMessage, pong); err != nil {
			return nil, &coinwerr.TransportError{Message: "failed to send Engine.IO pong", Err: err}
		}
		return nil, nil
	case '3', '5', '6': // pong / upgrade / noop
		return nil, nil
	case '4':
		return s.handleSocketIO(conn, raw[1:])
	default:
		return nil, &coinwerr.ProtocolError{Message: fmt.Sprintf("unknown Engine.IO packet type: %c", raw[0])}
	}
}

func (s *SocketIOSession) handleSocketIO(conn Conn, payload []byte) (*SocketIOMessage, error) {
	if len(payload) == 0 {
		return nil, nil
	}
	switch payload[0] {
	case '0': // CONNECT
		return nil, nil
	case '1': // DISCONNECT — caller should reconnect
		return nil, &coinwerr.TransportError{Message: "Socket.IO disconnect received"}
	case '2': // EVENT
		var arr []json.RawMessage
		if err := json.Unmarshal(payload[1:], &arr); err != nil || len(arr) == 0 {
			return nil, &coinwerr.ProtocolError{Message: "malformed Socket.IO event packet"}
		}
		var event string
		if err := json.Unmarshal(arr[0], &event); err != nil {
			return nil, &coinwerr.ProtocolError{Message: "malformed Socket.IO event name"}
		}
		msg := &SocketIOMessage{Event: event}
		if len(arr) > 1 {
			msg.Data = arr[1]
		}
		return msg, nil
	case '3': // ACK — unused by the documented API
		return nil, nil
	case '4': // ERROR
		return nil, &coinwerr.ProtocolError{Message: "Socket.IO error packet: " + string(payload[1:])}
	default:
		return nil, &coinwerr.ProtocolError{Message: fmt.Sprintf("unsupported Socket.IO packet type: %c", payload[0])}
	}
}

// Run is the blocking loop with reconnect + resubscribe.
func (s *SocketIOSession) Run(ctx context.Context, onEvent func(*SocketIOMessage) bool, onError func(error) bool) error {
	lastPing := time.Time{}
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		s.mu.Lock()
		connected := s.conn != nil
		s.mu.Unlock()

		if !connected {
			if err := s.Connect(ctx); err != nil {
				if !onError(err) {
					return err
				}
				if err := s.reconnect.wait(ctx); err != nil {
					return err
				}
				continue
			}
			s.reconnect.reset()
			lastPing = time.Time{}
			s.mu.Lock()
			args := make([]string, 0, len(s.args))
			for a := range s.args {
				args = append(args, a)
			}
			conn := s.conn
			s.mu.Unlock()
			for _, a := range args {
				frame, _ := sioEvent("subscribe", map[string]interface{}{"args": a})
				if err := conn.WriteMessage(websocket.TextMessage, frame); err != nil {
					break
				}
			}
		}

		if time.Since(lastPing) >= s.pingPeriod {
			if err := s.Ping(); err != nil {
				if !onError(err) {
					return err
				}
				s.closeQuietly()
				continue
			}
			lastPing = time.Now()
		}

		msg, err := s.ReadMessage()
		if err != nil {
			if !onError(err) {
				return err
			}
			s.closeQuietly()
			if err := s.reconnect.wait(ctx); err != nil {
				return err
			}
			continue
		}
		if msg == nil {
			continue
		}
		if !onEvent(msg) {
			return nil
		}
	}
}

// Close closes the connection.
func (s *SocketIOSession) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn != nil {
		err := s.conn.Close()
		s.conn = nil
		return err
	}
	return nil
}

func (s *SocketIOSession) closeQuietly() { _ = s.Close() }

func sioEvent(event string, payload map[string]interface{}) ([]byte, error) {
	arr, err := json.Marshal([]interface{}{event, payload})
	if err != nil {
		return nil, &coinwerr.ProtocolError{Message: "failed to encode Socket.IO event"}
	}
	return append([]byte("42"), arr...), nil
}

// ---------------------------------------------------------------------------
// DepthSequenceTracker — Spot incremental order book (Method 2, type=depth)
// ---------------------------------------------------------------------------

// DepthTrackResult classifies a tracked depth event.
type DepthTrackResult int

const (
	DepthApplied    DepthTrackResult = iota // event applied in sequence
	DepthDuplicate                          // out-of-order old event — discard
	DepthGap                                // startSeq > expected → resync from snapshot
	DepthNoSequence                         // event lacks startSeq/endSeq
)

// DepthSequenceTracker detects sequence gaps in the spot incremental order
// book (startSeq/endSeq per event). On DepthGap the caller must re-fetch a
// snapshot (depth_snapshot channel) and call Reset with its last seq.
type DepthSequenceTracker struct {
	expectedNext int64
	hasExpected  bool
}

// Track feeds one event.
func (t *DepthSequenceTracker) Track(startSeq, endSeq *int64) DepthTrackResult {
	if startSeq == nil || endSeq == nil {
		return DepthNoSequence
	}
	if !t.hasExpected {
		t.expectedNext = *endSeq + 1
		t.hasExpected = true
		return DepthApplied
	}
	if *endSeq < t.expectedNext {
		return DepthDuplicate
	}
	if *startSeq > t.expectedNext {
		// do not advance across the hole — resync and Reset first
		return DepthGap
	}
	t.expectedNext = *endSeq + 1
	return DepthApplied
}

// Reset anchors the sequence after a fresh snapshot.
func (t *DepthSequenceTracker) Reset(lastSnapshotSeq int64) {
	t.expectedNext = lastSnapshotSeq + 1
	t.hasExpected = true
}

// ExpectedNext returns the next expected startSeq (0 before anchoring).
func (t *DepthSequenceTracker) ExpectedNext() int64 {
	if !t.hasExpected {
		return 0
	}
	return t.expectedNext
}
