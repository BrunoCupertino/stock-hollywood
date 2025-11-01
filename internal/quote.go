package internal

import (
	"log/slog"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type QuoteActor struct {
	ticker      string
	px          float64
	updatedAt   time.Time
	subscribers map[*actor.PID]struct{}
}

type Subscription struct{}

type OnQuote struct {
	Ticker string
	Px     float64
	Date   time.Time
}

func (a *QuoteActor) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		slog.Info("quote actor has been started", "ticker", a.ticker)

		broadcaster := ctx.Engine().Registry.GetPID("broadcaster", "singleton")

		// todo load the last quote from the database
		a.px = 1
		a.updatedAt = time.Now().AddDate(0, 0, -1)

		ctx.Send(broadcaster, &OnQuote{
			Ticker: a.ticker,
			Px:     a.px,
			Date:   a.updatedAt,
		})

	case *Subscription:
		a.subscribers[ctx.Sender()] = struct{}{}
	case *OnQuote:
		if msg.Date.Before(a.updatedAt) {
			return
		}

		a.px = msg.Px
		a.updatedAt = msg.Date

		for s := range a.subscribers {
			ctx.Forward(s)
		}
	}
}

func NewQuoteActor(t string) actor.Producer {
	return func() actor.Receiver {
		return &QuoteActor{
			ticker:      t,
			subscribers: make(map[*actor.PID]struct{}),
		}
	}
}
