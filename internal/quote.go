package internal

import (
	"log/slog"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type QuoteActor struct {
	ticker      string
	quote       float64
	quotedAt    time.Time
	subscribers map[*actor.PID]struct{}
}

type Subscription struct {
	Subscriber *actor.PID
}

type OnQuote struct {
	Quote float64
	When  time.Time
}

func (a *QuoteActor) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		slog.Info("quote actor has been started", "ticker", a.ticker)
		// load last saved state
	case Subscription:
		a.subscribers[msg.Subscriber] = struct{}{}
	case *OnQuote:
		if msg.When.Before(a.quotedAt) {
			return
		}

		a.quote = msg.Quote
		a.quotedAt = msg.When

		for s := range a.subscribers {
			ctx.Send(s, msg)
		}

		slog.Info("quote broadcasted", "ticker", a.ticker, "quote", a.quote)
	}
}

func NewQuoteActor(ticker string) actor.Producer {
	return func() actor.Receiver {
		return &QuoteActor{
			ticker:      ticker,
			subscribers: make(map[*actor.PID]struct{}),
		}
	}
}
