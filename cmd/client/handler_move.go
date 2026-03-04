package main 

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
)

// this func will be called whenever a player move their army
func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.AckType{
	return func(move gamelogic.ArmyMove) pubsub.AckType{
		defer fmt.Print("> ")
		moveOutCome := gs.HandleMove(move)
		switch moveOutCome {
		case gamelogic.MoveOutComeSafe, gamelogic.MoveOutcomeMakeWar:
			return pubsub.Ack 
		default:
			return pubsub.NackDiscard
		}
	}
} 