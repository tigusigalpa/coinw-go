// CoinW Go SDK quickstart.
//
// Credentials come from environment variables — never hardcode keys.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	coinw "github.com/tigusigalpa/coinw-go"
)

func main() {
	opts := []coinw.ClientOption{}
	if k := os.Getenv("COINW_SPOT_API_KEY"); k != "" {
		opts = append(opts, coinw.WithSpotCredentials(k, os.Getenv("COINW_SPOT_SECRET_KEY")))
	}
	if k := os.Getenv("COINW_FUTURES_API_KEY"); k != "" {
		opts = append(opts, coinw.WithFuturesCredentials(k, os.Getenv("COINW_FUTURES_SECRET_KEY")))
	}
	if uid := os.Getenv("COINW_UID"); uid != "" {
		opts = append(opts, coinw.WithUID(uid))
	}

	c, err := coinw.NewClient(opts...)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// --- Public market data (no credentials) ---
	res, err := c.SpotMarket().GetRecentTrades(ctx, "BTC_USDT", 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res.Data))

	// --- Private reads ---
	// fmt.Println(c.SpotAccount().GetSpotAccountBalance(ctx))
	// fmt.Println(c.FuturesPosition().GetAllCurrentPositions(ctx))

	// --- Futures public WebSocket ---
	// ws := c.FuturesPublicWS()
	// _ = ws.Connect(context.Background())
	// _ = ws.Subscribe(map[string]interface{}{
	//     "biz": "perpetual", "type": "depth", "pairCode": "BTC_USDT",
	// })
	// for {
	//     msg, err := ws.ReadMessage()
	//     if err != nil { panic(err) }
	//     if msg != nil { fmt.Println(msg.Type, string(msg.Data)) }
	// }
}
