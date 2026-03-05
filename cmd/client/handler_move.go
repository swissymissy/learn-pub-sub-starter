package main 

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

// this func will be called whenever a player move their army
func handlerMove(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.ArmyMove) pubsub.AckType{
	return func(move gamelogic.ArmyMove) pubsub.AckType{
		defer fmt.Print("> ")
		username := gs.GetUsername()
		routingKey := routing.WarRecognitionsPrefix + "." + username

		moveOutCome := gs.HandleMove(move)
		switch moveOutCome {
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack 
		case gamelogic.MoveOutcomeMakeWar:
			// publish a message
			if err := pubsub.PublishJSON(
				ch,
				routing.ExchangePerilTopic,
				routingKey,
				gamelogic.RecognitionOfWar{
					Attacker: move.Player,
					Defender: gs.GetPlayerSnap(),
				},
			); err != nil {
				fmt.Printf("Error publishing message: %s\n", err)
				return pubsub.NackRequeue	// put message in queue again, and retry
			}
			return pubsub.Ack
		default:
			return pubsub.NackDiscard
		}
	}
} 