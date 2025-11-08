package internal

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/anthdm/hollywood/actor"
)

const TickersNum = 100_000

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
		now := time.Now().UTC()

		for _, t := range a.tickers {
			pid := ctx.Engine().Registry.GetPID("quote", t)
			if pid == nil {
				pid = ctx.Engine().Spawn(NewQuoteActor(t), "quote", actor.WithID(t))
			}

			a.ticksPublished++

			q := &OnQuote{
				Ticker: t,
				Px:     1.1 + float64(time.Now().UnixMilli()%10),
				Date:   time.Now().UTC(),
			}

			if t == "LAST" {
				q.Date = now
			}

			ctx.Send(pid, q)

			if a.ticksPublished%100_000 == 0 {
				slog.Info("100k ticks published")
			}
		}

		// a.repeater.Stop()

		// slog.Info("broadcast publisher time", "duration", time.Since(now))
	}

}

func NewFakeConnectorActor() actor.Producer {
	return func() actor.Receiver {
		tickers := make([]string, 0, TickersNum+3)

		tickers = append(tickers, "APPL")
		tickers = append(tickers, "GOOGL")

		for i := range TickersNum + 2 {
			tickers = append(tickers, fmt.Sprintf("TICKER%d", i))
		}

		tickers = append(tickers, "LAST")

		return &FakeConnectorActor{
			tickers:  tickers,
			internal: time.Second * 1,
		}
	}
}
