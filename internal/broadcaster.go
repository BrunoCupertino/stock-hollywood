package internal

import (
	"log/slog"

	"github.com/anthdm/hollywood/actor"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const batchSize = 100

type BroadcasterActor struct {
	quotesActors map[string]*actor.PID
	subscribers  map[string]map[*actor.PID]struct{}
	batch        []*Quote
}

func (a *BroadcasterActor) Receive(ctx *actor.Context) {
	switch msg := ctx.Message().(type) {
	case actor.Started:
		slog.Info("broadcaster actor has been started")
	case *QuoteSubscriptionRequest:
		quotePID, ok := a.quotesActors[msg.Ticker]
		if !ok {
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

		ctx.Send(quotePID, &Snapshot{
			pid: ctx.Sender(),
		})

		slog.Info("subscription made for", "id", ctx.Sender().ID)
	case *OnQuote:
		ticker := msg.Ticker

		// now := time.Now()

		// snapshot
		if msg.PID != nil {
			ctx.Send(msg.PID, &Quote{
				Ticker: ticker,
				Px:     msg.Px,
				Date:   timestamppb.New(msg.Date),
			})
			return
		}

		a.batch = append(a.batch, &Quote{
			Ticker: ticker,
			Px:     msg.Px,
			Date:   timestamppb.New(msg.Date),
		})

		if len(a.batch) == batchSize {
			for subscriber := range a.subscribers[ticker] {
				ctx.Send(subscriber, &QuoteBatch{
					Quotes: a.batch,
				})
			}
			a.batch = a.batch[:0]
		}

		// for subscriber := range a.subscribers[ticker] {
		// 	ctx.Send(subscriber, &Quote{
		// 		Ticker: ticker,
		// 		Px:     msg.Px,
		// 		Date:   timestamppb.New(msg.Date),
		// 	})
		// }

		// slog.Info("broadcast broadcaster time",
		// 	"duration", time.Now().UTC().Sub(msg.Date),
		// 	"len", len(a.subscribers[ticker]))
	}
}

func NewBroadcasterActor() actor.Producer {
	return func() actor.Receiver {
		return &BroadcasterActor{
			quotesActors: make(map[string]*actor.PID),
			subscribers:  make(map[string]map[*actor.PID]struct{}),
			batch:        make([]*Quote, 0, batchSize),
		}
	}
}
