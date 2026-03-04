package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
)

// func to be passed in SubscribeJSON. 
// Will be called each time a new msg is consumed
func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType{
	return func( playingState routing.PlayingState) pubsub.AckType {
		defer fmt.Print("> ")
		gs.HandlePause(playingState)
		return pubsub.Ack
	}
}


