package coinw_test

import (
	"context"
	"os"
	"testing"
	"time"

	coinw "github.com/tigusigalpa/coinw-go"
	"github.com/tigusigalpa/coinw-go/model"
	ws "github.com/tigusigalpa/coinw-go/websocket"
)

// These tests make real network requests only after an explicit environment
// opt-in. No live response payload or credential is written to test output.
func integrationClient(t *testing.T) *coinw.Client {
	t.Helper()
	if os.Getenv("COINW_INTEGRATION") != "1" {
		t.Skip("set COINW_INTEGRATION=1 and COINW_TEST_BASE_URL to run live tests")
	}
	base := requiredEnv(t, "COINW_TEST_BASE_URL")
	c, err := coinw.NewClient(coinw.WithRestBaseURL(base), coinw.WithMaxReadRetries(0), coinw.WithRequestTimeout(10*time.Second), coinw.WithSpotCredentials(os.Getenv("COINW_TEST_SPOT_API_KEY"), os.Getenv("COINW_TEST_SPOT_SECRET_KEY")), coinw.WithFuturesCredentials(os.Getenv("COINW_TEST_FUTURES_API_KEY"), os.Getenv("COINW_TEST_FUTURES_SECRET_KEY")), coinw.WithWebSocketURLs(os.Getenv("COINW_TEST_SPOT_WS_URL"), os.Getenv("COINW_TEST_FUTURES_WS_URL"), os.Getenv("COINW_TEST_SOCKETIO_URL"), os.Getenv("COINW_TEST_TOKEN_URL")))
	if err != nil {
		t.Fatal(err)
	}
	return c
}
func requiredEnv(t *testing.T, key string) string {
	t.Helper()
	v := os.Getenv(key)
	if v == "" {
		t.Fatalf("explicit integration configuration requires %s", key)
	}
	return v
}
func liveContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func TestIntegrationPublicREST(t *testing.T) {
	c := integrationClient(t)
	if _, err := c.SpotMarket().GetInstrumentInformation(liveContext(t)); err != nil {
		t.Fatal(err)
	}
	if _, err := c.FuturesMarket().GetInstrumentInfo(liveContext(t), ""); err != nil {
		t.Fatal(err)
	}
}
func TestIntegrationPrivateReadOnly(t *testing.T) {
	c := integrationClient(t)
	t.Run("Spot", func(t *testing.T) {
		if os.Getenv("COINW_TEST_SPOT_API_KEY") == "" || os.Getenv("COINW_TEST_SPOT_SECRET_KEY") == "" {
			t.Skip("explicit Spot test credentials are required")
		}
		if _, err := c.SpotAccount().GetSpotAccountBalance(liveContext(t)); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("Futures", func(t *testing.T) {
		if os.Getenv("COINW_TEST_FUTURES_API_KEY") == "" || os.Getenv("COINW_TEST_FUTURES_SECRET_KEY") == "" {
			t.Skip("explicit Futures test credentials are required")
		}
		if _, err := c.FuturesAssets().GetFuturesAccountAssets(liveContext(t)); err != nil {
			t.Fatal(err)
		}
	})
}

func TestIntegrationPublicWebSockets(t *testing.T) {
	c := integrationClient(t)
	t.Run("SpotRaw", func(t *testing.T) {
		if os.Getenv("COINW_TEST_SPOT_WS_URL") == "" {
			t.Skip("explicit Spot WS URL required")
		}
		s := c.SpotPublicWS()
		defer func() { _ = s.Close() }()
		ctx := liveContext(t)
		if err := s.Connect(ctx); err != nil {
			t.Fatal(err)
		}
		stream := ws.SpotAllTickers()
		if err := ws.SubscribeStream(ctx, s, stream); err != nil {
			t.Fatal(err)
		}
		for {
			m, err := s.ReadMessage()
			if err != nil {
				t.Fatal(err)
			}
			if m != nil && m.IsEvent() {
				if _, err = ws.DecodeStream(m, stream); err != nil {
					t.Fatal(err)
				}
				return
			}
		}
	})
	t.Run("FuturesRaw", func(t *testing.T) {
		if os.Getenv("COINW_TEST_FUTURES_WS_URL") == "" {
			t.Skip("explicit Futures WS URL required")
		}
		s := c.FuturesPublicWS()
		defer func() { _ = s.Close() }()
		ctx := liveContext(t)
		if err := s.Connect(ctx); err != nil {
			t.Fatal(err)
		}
		instrument, _ := model.USDTInstrument("BTC")
		stream := ws.FuturesTicker(instrument)
		if err := ws.SubscribeStream(ctx, s, stream); err != nil {
			t.Fatal(err)
		}
		for {
			m, err := s.ReadMessage()
			if err != nil {
				t.Fatal(err)
			}
			if m != nil && m.IsEvent() {
				if _, err = ws.DecodeStream(m, stream); err != nil {
					t.Fatal(err)
				}
				return
			}
		}
	})
	t.Run("SpotSocketIO", func(t *testing.T) {
		if os.Getenv("COINW_TEST_SOCKETIO_URL") == "" || os.Getenv("COINW_TEST_TOKEN_URL") == "" {
			t.Skip("explicit Socket.IO and public-token URLs required")
		}
		s := c.SpotSocketIOWS()
		defer func() { _ = s.Close() }()
		ctx := liveContext(t)
		if err := s.Connect(ctx); err != nil {
			t.Fatal(err)
		}
		stream := ws.SpotTickerSocketIO("BTC-USDT")
		if err := ws.SubscribeSocketIO(ctx, s, stream); err != nil {
			t.Fatal(err)
		}
		for {
			m, err := s.ReadMessage()
			if err != nil {
				t.Fatal(err)
			}
			if m != nil {
				if _, err = ws.DecodeSocketIO(m, stream); err != nil {
					t.Fatal(err)
				}
				return
			}
		}
	})
}

func TestIntegrationMutationCancelSpecifiedOrder(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("live mutation tests are disabled in CI")
	}
	if os.Getenv("COINW_TEST_MUTATIONS") != "1" {
		t.Skip("separate COINW_TEST_MUTATIONS=1 opt-in required")
	}
	c := integrationClient(t)
	requiredEnv(t, "COINW_TEST_FUTURES_API_KEY")
	requiredEnv(t, "COINW_TEST_FUTURES_SECRET_KEY")
	id := requiredEnv(t, "COINW_TEST_CANCEL_ORDER_ID")
	// The human supplies an existing expendable order. The suite never places,
	// withdraws, transfers, or creates a position to prepare this test.
	if _, err := c.FuturesTrade().CancelAnOrder(liveContext(t), id); err != nil {
		t.Fatal(err)
	}
}
