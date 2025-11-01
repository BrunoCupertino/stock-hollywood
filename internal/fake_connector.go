package internal

import (
	"log/slog"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type FakeConnectorActor struct {
	tickers  []string
	internal time.Duration
	repeater actor.SendRepeater
	me       *actor.PID
}

type RefreshQuotes struct{}

func (a *FakeConnectorActor) Receive(ctx *actor.Context) {
	switch ctx.Message().(type) {
	case actor.Started:
		slog.Info("fake connector actor has been starded")

		a.me = ctx.PID()

		a.repeater = ctx.Engine().SendRepeat(a.me, &RefreshQuotes{}, a.internal)

	case *RefreshQuotes:
		for _, t := range a.tickers {
			pid := ctx.Engine().Registry.GetPID("quote", t)
			if pid == nil {
				pid = ctx.Engine().Spawn(NewQuoteActor(t), "quote", actor.WithID(t))
			}
			ctx.Send(pid, &OnQuote{
				Ticker: t,
				Px:     1.1 + float64(time.Now().UnixMilli()%10),
				Date:   time.Now(),
			})
		}
	}

}

func NewFakeConnectorActor() actor.Producer {
	return func() actor.Receiver {
		return &FakeConnectorActor{
			tickers:  []string{"GOOGL", "APPL", "UND", "XP", "NU", "BB"},
			internal: time.Millisecond * 5,
		}
	}
}
