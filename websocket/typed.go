package websocket

import (
	"context"
	"encoding/json"
	"strings"

	coinwerr "github.com/tigusigalpa/coinw-go/errors"
	"github.com/tigusigalpa/coinw-go/model"
)

type Profile string

const (
	SpotProfile    Profile = "exchange"
	FuturesProfile Profile = "futures"
)

func WithProfile(profile Profile) RawOption { return func(s *RawJsonSession) { s.profile = profile } }

// Stream binds an immutable documented subscription to its response contract.
type Stream[T any] struct {
	id, kind, pair, interval, args string
	profile                        Profile
	private                        bool
	needsPair, needsInterval       bool
}

func (s Stream[T]) ID() string { return s.id }
func (s Stream[T]) Params() map[string]interface{} {
	p := map[string]interface{}{"biz": string(s.profile), "type": s.kind}
	if s.pair != "" {
		p["pairCode"] = s.pair
	}
	if s.interval != "" {
		p["interval"] = s.interval
	}
	return p
}
func (s Stream[T]) validate(socketIO bool) error {
	if s.id == "" || (s.args != "") != socketIO {
		return &coinwerr.ValidationError{Message: "stream uses a different protocol adapter"}
	}
	if s.needsPair && strings.TrimSpace(s.pair) == "" {
		return &coinwerr.ValidationError{Message: "stream instrument must not be empty"}
	}
	if s.needsInterval && s.interval == "" {
		return &coinwerr.ValidationError{Message: "stream interval must not be empty"}
	}
	if s.interval != "" {
		allowed := "|1m|3m|5m|15m|30m|1h|2h|4h|6h|12h|1d|1w|1M|"
		if s.profile == FuturesProfile {
			allowed = "|1|3|5|15|30|1H|4H|1D|1W|1M|"
		}
		if !strings.Contains(allowed, "|"+s.interval+"|") {
			return &coinwerr.ValidationError{Message: "undocumented stream interval"}
		}
	}
	return nil
}

func SubscribeStream[T any](ctx context.Context, session *RawJsonSession, stream Stream[T]) error {
	return changeStream(ctx, session, stream, false)
}
func UnsubscribeStream[T any](ctx context.Context, session *RawJsonSession, stream Stream[T]) error {
	return changeStream(ctx, session, stream, true)
}
func changeStream[T any](ctx context.Context, s *RawJsonSession, stream Stream[T], unsubscribe bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := stream.validate(false); err != nil {
		return err
	}
	if s == nil {
		return &coinwerr.ValidationError{Message: "nil WS session"}
	}
	if s.private != stream.private || (s.profile != "" && s.profile != stream.profile) {
		return &coinwerr.AuthenticationError{Message: "stream requires its own public/private business profile"}
	}
	if stream.private && s.onConnect == nil && !unsubscribe {
		return &coinwerr.ValidationError{Message: "private typed subscriptions require WithOnConnect to fetch the corresponding REST snapshot"}
	}
	// Cancelling interrupts a blocked send; Run reconnects and restores the set.
	stop := context.AfterFunc(ctx, s.dropQuietly)
	defer stop()
	var err error
	if unsubscribe {
		err = s.Unsubscribe(stream.Params())
	} else {
		err = s.Subscribe(stream.Params())
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return err
}
func SubscribeSocketIO[T any](ctx context.Context, s *SocketIOSession, stream Stream[T]) error {
	return changeSocketIO(ctx, s, stream, false)
}
func UnsubscribeSocketIO[T any](ctx context.Context, s *SocketIOSession, stream Stream[T]) error {
	return changeSocketIO(ctx, s, stream, true)
}
func changeSocketIO[T any](ctx context.Context, s *SocketIOSession, stream Stream[T], unsubscribe bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := stream.validate(true); err != nil {
		return err
	}
	if s == nil {
		return &coinwerr.ValidationError{Message: "nil Socket.IO session"}
	}
	stop := context.AfterFunc(ctx, s.dropQuietly)
	defer stop()
	var err error
	if unsubscribe {
		err = s.Unsubscribe(stream.args)
	} else {
		err = s.Subscribe(stream.args)
	}
	if ctx.Err() != nil {
		return ctx.Err()
	}
	return err
}

type Acknowledgement struct {
	Biz      *string `json:"biz,omitempty"`
	Type     *string `json:"type,omitempty"`
	PairCode *string `json:"pairCode,omitempty"`
	Interval *string `json:"interval,omitempty"`
	Channel  string  `json:"channel" required:"true"`
	Data     struct {
		Result bool `json:"result" required:"true"`
	} `json:"data" required:"true"`
}

func DecodeAcknowledgement(raw []byte) (*Acknowledgement, error) {
	var ack Acknowledgement
	if err := model.Decode(raw, &ack); err != nil {
		return nil, &coinwerr.ProtocolError{Message: "malformed WS acknowledgement"}
	}
	if ack.Channel != "login" && ack.Channel != "subscribe" && ack.Channel != "unsubscribe" {
		return nil, &coinwerr.ProtocolError{Message: "unknown acknowledgement channel"}
	}
	return &ack, nil
}

type streamEnvelope[T any] struct {
	Biz      string                 `json:"biz" required:"true"`
	Type     string                 `json:"type" required:"true"`
	PairCode *string                `json:"pairCode,omitempty"`
	Interval *string                `json:"interval,omitempty"`
	Data     T                      `json:"data" required:"true"`
	TS       *model.StringOrInteger `json:"ts,omitempty"`
	N        *string                `json:"n,omitempty"`
}

// DecodeStream validates metadata and payload. It never orders events by
// timestamp, invents an initial snapshot, or converts quantity units.
func DecodeStream[T any](message *Message, stream Stream[T]) (T, error) {
	var zero T
	if message == nil {
		return zero, &coinwerr.ProtocolError{Message: "nil stream message"}
	}
	if err := stream.validate(false); err != nil {
		return zero, err
	}
	var e streamEnvelope[T]
	if err := model.Decode(message.RawBytes, &e); err != nil {
		return zero, &coinwerr.ProtocolError{Message: "stream contract " + stream.id + ": " + err.Error()}
	}
	if e.Biz != string(stream.profile) || e.Type != stream.kind {
		return zero, &coinwerr.ProtocolError{Message: "stream metadata does not match subscription"}
	}
	if stream.pair != "" && (e.PairCode == nil || !strings.EqualFold(*e.PairCode, stream.pair)) {
		return zero, &coinwerr.ProtocolError{Message: "stream pair does not match subscription"}
	}
	if stream.interval != "" && (e.Interval == nil || *e.Interval != stream.interval) {
		return zero, &coinwerr.ProtocolError{Message: "stream interval does not match subscription"}
	}
	return e.Data, nil
}
func DecodeSocketIO[T any](message *SocketIOMessage, stream Stream[T]) (T, error) {
	var zero T
	if message == nil || message.Event != "subscribe" {
		return zero, &coinwerr.ProtocolError{Message: "unexpected Socket.IO event"}
	}
	if err := stream.validate(true); err != nil {
		return zero, err
	}
	var e struct {
		Channel string `json:"channel" required:"true"`
		Subject string `json:"subject" required:"true"`
		Data    T      `json:"data" required:"true"`
	}
	if err := model.Decode(message.Data, &e); err != nil {
		return zero, &coinwerr.ProtocolError{Message: "Socket.IO stream contract " + stream.id + ": " + err.Error()}
	}
	if e.Channel != stream.args || e.Subject != strings.Split(stream.args, ":")[0] {
		return zero, &coinwerr.ProtocolError{Message: "Socket.IO channel does not match subscription"}
	}
	return e.Data, nil
}

// DecodeRawFrame exposes the checked raw envelope for explicit dispatch.
func DecodeRawFrame(raw json.RawMessage) (*Message, error) { return decodeMessage(raw) }
