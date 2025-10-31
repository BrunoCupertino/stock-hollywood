package main

import (
	"github.com/BrunoCupertino/stock-hollywood/internal"
	"github.com/anthdm/hollywood/actor"
)

func main() {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		panic(err)
	}

	_ = e.Spawn(internal.NewFakeConnectorActor(), "connector", actor.WithID("fake"))

	select {}
}
