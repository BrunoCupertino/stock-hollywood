package main

import (
	"log/slog"

	"github.com/BrunoCupertino/stock-hollywood/internal"
	"github.com/anthdm/hollywood/actor"
)

func main() {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		panic(err)
	}

	_ = e.Spawn(internal.NewFakeConnectorActor(), "connector", actor.WithID("fake"))

	broadcaster := e.Spawn(internal.NewBroadcasterActor(), "broadcaster", actor.WithID("singleton"))

	e.SpawnFunc(func(ctx *actor.Context) {
		switch msg := ctx.Message().(type) {
		case actor.Started:
			slog.Info("subscriber actor has been started")

			ctx.Send(broadcaster, &internal.QuoteSubscriptionRequest{
				Ticker: "GOOGL",
			})

			ctx.Send(broadcaster, &internal.QuoteSubscriptionRequest{
				Ticker: "APPL",
			})
		case *internal.Quote:
			slog.Info("new quote received", "ticker", msg.Ticker, "px", msg.Px)

		case actor.Stopped:
			slog.Warn("subscriber actor has been stopped")
		}
	}, "subcriber", actor.WithID("singleton"))

	select {}
}
