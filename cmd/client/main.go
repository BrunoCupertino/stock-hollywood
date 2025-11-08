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

	each := func(ctx *actor.Context, msg *internal.Quote) {
		if time.Now().UTC().Sub(msg.Date.AsTime()) > time.Millisecond*100 {
			slog.Info("new quote received",
				"ticker", msg.Ticker,
				"px", msg.Px, "id", ctx.PID().ID,
				"duration", time.Since(msg.Date.AsTime()),
				"now", time.Now().UTC(),
				"date", msg.Date.AsTime())
		}

		atomic.AddInt64(&quoteReceived, 1)

		if atomic.LoadInt64(&quoteReceived)%100_000 == 0 {
			slog.Info("100k quotes received from server")
		}
	}

	subscriber := func(ctx *actor.Context) {
		switch msg := ctx.Message().(type) {
		case actor.Started:
			slog.Info("subscriber actor has been started", "id", ctx.PID().ID)

			ctx.Send(remoteBroadcaster, &internal.QuoteSubscriptionRequest{
				Ticker: "APPL",
			})

			ctx.Send(remoteBroadcaster, &internal.QuoteSubscriptionRequest{
				Ticker: "GOOGL",
			})

			for i := range internal.TickersNum {
				ctx.Send(remoteBroadcaster, &internal.QuoteSubscriptionRequest{
					Ticker: fmt.Sprintf("TICKER%d", i),
				})
			}

			ctx.Send(remoteBroadcaster, &internal.QuoteSubscriptionRequest{
				Ticker: "LAST",
			})
		case *internal.Quote:
			each(ctx, msg)

		case *internal.QuoteBatch:
			for _, q := range msg.Quotes {
				each(ctx, q)
			}
		case actor.Stopped:
			slog.Warn("subscriber actor has been stopped")
		}
	}

	e.SpawnFunc(subscriber, "subcriber", actor.WithID("singleton"))

	select {}
}
