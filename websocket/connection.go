package websocket

import (
	"context"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// sessionConn enforces gorilla's one-reader/one-writer contract. Close is
// independent of those locks so cancellation also interrupts blocked I/O.
type sessionConn struct {
	Conn
	ctx       context.Context
	readMu    sync.Mutex
	writeMu   sync.Mutex
	closeOnce sync.Once
	done      chan struct{}
	closeErr  error
	hookOnce  sync.Once
	hookErr   error
}

func newSessionConn(ctx context.Context, conn Conn) *sessionConn {
	c := &sessionConn{Conn: conn, ctx: ctx, done: make(chan struct{})}
	if ctx.Done() != nil {
		go func() {
			select {
			case <-ctx.Done():
				_ = c.Close()
			case <-c.done:
			}
		}()
	}
	return c
}

func (c *sessionConn) active() bool {
	select {
	case <-c.done:
		return false
	default:
		return true
	}
}

// prepare runs the per-connection hook once, before the first read. Frames
// that arrive while it runs wait in the socket, so none are lost.
func (c *sessionConn) prepare(hook func(context.Context) error) error {
	if hook == nil {
		return nil
	}
	c.hookOnce.Do(func() { c.hookErr = hook(c.ctx) })
	return c.hookErr
}

func (c *sessionConn) WriteMessage(mt int, data []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if err := c.ctx.Err(); err != nil {
		return err
	}
	if w, ok := c.Conn.(interface{ SetWriteDeadline(time.Time) error }); ok {
		deadline := time.Now().Add(10 * time.Second)
		if d, ok := c.ctx.Deadline(); ok && d.Before(deadline) {
			deadline = d
		}
		if err := w.SetWriteDeadline(deadline); err != nil {
			return err
		}
	}
	return c.Conn.WriteMessage(mt, data)
}

func (c *sessionConn) read(timeout time.Duration) (int, []byte, error) {
	c.readMu.Lock()
	defer c.readMu.Unlock()
	if err := c.ctx.Err(); err != nil {
		return 0, nil, err
	}
	deadline := time.Now().Add(timeout)
	if d, ok := c.ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}
	if err := c.SetReadDeadline(deadline); err != nil {
		return 0, nil, err
	}
	return c.ReadMessage()
}

func (c *sessionConn) Close() error {
	c.closeOnce.Do(func() {
		close(c.done)
		c.closeErr = c.Conn.Close()
	})
	return c.closeErr
}

// runIO keeps heartbeats independent of incoming traffic. The consumer runs
// synchronously, providing backpressure without an unbounded event queue.
func runIO(ctx context.Context, period time.Duration, ping func() error, read func() (interface{}, error), onEvent func(interface{}) bool, conn *sessionConn) error {
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	pingErr := make(chan error, 1)
	done := make(chan struct{})
	go func() {
		defer close(done)
		ticker := time.NewTicker(period)
		defer ticker.Stop()
		for {
			select {
			case <-runCtx.Done():
				_ = conn.Close()
				return
			case <-ticker.C:
				if err := ping(); err != nil {
					pingErr <- err
					_ = conn.Close()
					return
				}
			}
		}
	}()
	defer func() { cancel(); <-done }()
	for {
		m, err := read()
		if ctx.Err() != nil {
			return ctx.Err()
		}
		select {
		case err = <-pingErr:
		default:
		}
		if err != nil {
			return err
		}
		if m != nil && !onEvent(m) {
			return nil
		}
	}
}

// Ensure the production connection supports bounded writes.
var _ interface{ SetWriteDeadline(time.Time) error } = (*websocket.Conn)(nil)
