package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

// func to be passed in SubscribeJSON. 
// Will be called each time a new msg is consumed
func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState){
	return func( playingState routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(playingState)
	}
}


