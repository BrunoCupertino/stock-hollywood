package main

import (
	"github.com/BrunoCupertino/stock-hollywood/internal"
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
)

func main() {
	// r := remote.New(":4000", remote.NewConfig().WithBufferSize(1024*4))
	r := remote.New(":4000", remote.NewConfig())

	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(r))
	if err != nil {
		panic(err)
	}

	_ = e.Spawn(internal.NewFakeConnectorActor(), "connector", actor.WithID("fake"))

	// todo: broadcaster by consumer
	_ = e.Spawn(internal.NewBroadcasterActor(), "broadcaster", actor.WithID("singleton"))

	select {}
}
