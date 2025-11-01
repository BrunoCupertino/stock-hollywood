package internal

import (
	"log/slog"

	"github.com/anthdm/hollywood/actor"
)

type BroadcasterActor struct {
	quotesActors map[string]*actor.PID
	subscribers  map[string]map[*actor.PID]struct{}
	wildcard     map[*actor.PID]struct{}
}

func (a *BroadcasterActor) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		slog.Info("broadcaster actor has been started")
	case *QuoteSubscriptionRequest:
		if msg.Ticker == "*" {
			// todo how to findout all the tickers and spawn their actors?
			a.wildcard[ctx.Sender()] = struct{}{}
			return
		}

		if _, ok := a.quotesActors[msg.Ticker]; !ok {
			pid := ctx.Engine().Registry.GetPID("quote", msg.Ticker)
			if pid == nil {
				pid = ctx.Engine().Spawn(NewQuoteActor(msg.Ticker), "quote", actor.WithID(msg.Ticker))
			}

			a.quotesActors[msg.Ticker] = pid

			ctx.Send(pid, &Subscription{})
		}

		if _, ok := a.subscribers[msg.Ticker]; !ok {
			a.subscribers[msg.Ticker] = make(map[*actor.PID]struct{})
		}
		a.subscribers[msg.Ticker][ctx.Sender()] = struct{}{}
	case *OnQuote:
		ticker := msg.Ticker

		for subscriber := range a.subscribers[ticker] {
			ctx.Send(subscriber, &Quote{
				Ticker: ticker,
				Px:     msg.Px,
				// Date:   msg.Date,
			})
		}

		for subscriber := range a.wildcard {
			ctx.Send(subscriber, &Quote{
				Ticker: ticker,
				Px:     msg.Px,
				// Date:   msg.Date,
			})
		}
	}
}

func NewBroadcasterActor() actor.Producer {
	return func() actor.Receiver {
		return &BroadcasterActor{
			quotesActors: make(map[string]*actor.PID),
			subscribers:  make(map[string]map[*actor.PID]struct{}),
			wildcard:     make(map[*actor.PID]struct{}),
		}
	}
}
