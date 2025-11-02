package internal

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type FakeConnectorActor struct {
	ticksPublished int64
	tickers        []string
	internal       time.Duration
	repeater       actor.SendRepeater
	me             *actor.PID
}

type RefreshQuotes struct{}

func (a *FakeConnectorActor) Receive(ctx *actor.Context) {
	switch ctx.Message().(type) {
	case actor.Started:
		slog.Info("fake connector actor has been starded")

		a.me = ctx.PID()

		a.repeater = ctx.Engine().SendRepeat(a.me, &RefreshQuotes{}, a.internal)

	case *RefreshQuotes:
		// now := time.Now()

		for _, t := range a.tickers {
			pid := ctx.Engine().Registry.GetPID("quote", t)
			if pid == nil {
				pid = ctx.Engine().Spawn(NewQuoteActor(t), "quote", actor.WithID(t))
			}

			a.ticksPublished++

			ctx.Send(pid, &OnQuote{
				Ticker: t,
				Px:     1.1 + float64(time.Now().UnixMilli()%10),
				Date:   time.Now().UTC(),
			})
		}

		if a.ticksPublished%1_000_000 == 0 {
			slog.Info("million ticks published")
		}

		// slog.Info("broadcast publisher time", "duration", time.Since(now))
	}

}

func NewFakeConnectorActor() actor.Producer {
	return func() actor.Receiver {
		tickers := make([]string, 0, 35_000)

		tickers = append(tickers, "APPL")
		tickers = append(tickers, "GOOGL")

		for i := range 35_000 - 2 {
			tickers = append(tickers, fmt.Sprintf("TICKER%d", i))
		}

		return &FakeConnectorActor{
			tickers: tickers,
			// tickers:  []string{"GOOGL", "APPL", "A", "B", "C", "D", "E", "F", "G", "H"},
			internal: time.Millisecond * 500,
		}
	}
}
