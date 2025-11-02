package main

import (
	"flag"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/BrunoCupertino/stock-hollywood/internal"
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
)

func main() {
	port := flag.String("port", ":3000", "port")
	flag.Parse()

	r := remote.New(*port, remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		panic(err)
	}

	var quoteReceived int64

	remoteBroadcaster := actor.NewPID(":4000", "broadcaster/singleton")

	subscriber := func(ctx *actor.Context) {
		switch msg := ctx.Message().(type) {
		case actor.Started:
			slog.Info("subscriber actor has been started", "id", ctx.PID().ID)

			ctx.Send(remoteBroadcaster, &internal.QuoteSubscriptionRequest{
				Ticker: "GOOGL",
			})

			ctx.Send(remoteBroadcaster, &internal.QuoteSubscriptionRequest{
				Ticker: "APPL",
			})

			// if ctx.PID().ID != "subcriber/1" {
			for i := range 35_000 {
				ctx.Send(remoteBroadcaster, &internal.QuoteSubscriptionRequest{
					Ticker: fmt.Sprintf("TICKER%d", i),
				})
			}
			// }
		case *internal.Quote:
			if time.Now().UTC().Sub(msg.Date.AsTime()) > time.Millisecond*5 {
				slog.Info("new quote received",
					"ticker", msg.Ticker,
					"px", msg.Px, "id", ctx.PID().ID,
					"duration", time.Since(msg.Date.AsTime()),
					"now", time.Now().UTC(),
					"date", msg.Date.AsTime())
			}

			atomic.AddInt64(&quoteReceived, 1)

			if atomic.LoadInt64(&quoteReceived)%1_000 == 0 {
				slog.Info("1k quotes received from server")
			}

			_ = msg
		case actor.Stopped:
			slog.Warn("subscriber actor has been stopped")
		}
	}

	e.SpawnFunc(subscriber, "subcriber", actor.WithID("1"))

	select {}
}
