package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
)

// consumes all war messages published by handlerMove
func handlerWar(gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(rw gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")

		warOutCome, _, _ := gs.HandleWar(rw)
		switch warOutCome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue			// requeue to other client can try to consume it
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon, gamelogic.WarOutcomeYouWon, gamelogic.WarOutcomeDraw:
			return pubsub.Ack
		default:
			fmt.Println("Error: unkown war outcome")
			return pubsub.NackDiscard
		}
	}
}