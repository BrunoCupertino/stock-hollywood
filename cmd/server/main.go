package main

import (
	"github.com/BrunoCupertino/stock-hollywood/internal"
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
)

func main() {
	r := remote.New(":4000", remote.NewConfig().WithBufferSize(1024*1024*32))

	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		panic(err)
	}

	_ = e.Spawn(internal.NewFakeConnectorActor(), "connector", actor.WithID("fake"))

	_ = e.Spawn(internal.NewBroadcasterActor(), "broadcaster", actor.WithID("singleton"))

	// subscriber := func(ctx *actor.Context) {
	// 	switch msg := ctx.Message().(type) {
	// 	case actor.Started:
	// 		slog.Info("subscriber actor has been started", "id", ctx.PID().ID)

	// 		ctx.Send(broadcaster, &internal.QuoteSubscriptionRequest{
	// 			Ticker: "GOOGL",
	// 		})

	// 		ctx.Send(broadcaster, &internal.QuoteSubscriptionRequest{
	// 			Ticker: "APPL",
	// 		})

	// 		if ctx.PID().ID != "subcriber/1" {
	// 			for i := range 100 {
	// 				ctx.Send(broadcaster, &internal.QuoteSubscriptionRequest{
	// 					Ticker: fmt.Sprintf("TICKER%d", i),
	// 				})
	// 			}
	// 		}
	// 	case *internal.Quote:
	// 		if time.Now().UTC().Sub(msg.Date.AsTime()) > time.Millisecond*5 {
	// 			slog.Info("new quote received",
	// 				"ticker", msg.Ticker,
	// 				"px", msg.Px, "id", ctx.PID().ID,
	// 				"duration", time.Since(msg.Date.AsTime()),
	// 				"now", time.Now().UTC(),
	// 				"date", msg.Date.AsTime())
	// 		}
	// 		_ = msg
	// 	case actor.Stopped:
	// 		slog.Warn("subscriber actor has been stopped")
	// 	}
	// }

	// e.SpawnFunc(subscriber, "subcriber", actor.WithID("1"))

	// time.Sleep(time.Second * 3)

	// for i := range 10 {
	// 	e.SpawnFunc(subscriber, "subcriber", actor.WithID(fmt.Sprintf("%d", i+2)))
	// }

	select {}
}
