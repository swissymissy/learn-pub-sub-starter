package main

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

// consumes all war messages published by handlerMove and publish game logs
func handlerWar(gs *gamelogic.GameState, ch *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(rw gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		username := gs.GetUsername()
	
		warOutCome, winner, loser := gs.HandleWar(rw)
		switch warOutCome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue			// requeue so another client can try to consume it
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			msg := fmt.Sprintf("%s won a war against %s", winner, loser)
			if err := publishGameLog(ch, username, msg); err != nil {
				fmt.Printf("Error publishing message: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeYouWon:
			msg := fmt.Sprintf("%s won a war against %s", winner, loser)
			if err := publishGameLog(ch, username, msg); err != nil {
				fmt.Printf("Error publishing message: %s\n", err)
				return pubsub.NackRequeue
			} 
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			msg := fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
			if err := publishGameLog(ch, username, msg); err != nil {
				fmt.Printf("Error publishing message: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		default:
			fmt.Println("Error: unknown war outcome")
			return pubsub.NackDiscard
		}
	}
}

// helper function
func publishGameLog(ch *amqp.Channel, username, msg string) error {
	// create GameLog struct
	gameLogStruct := routing.GameLog{
		CurrentTime: time.Now(),
		Message: msg,
		Username: username,
	}

	//create routing key
	routingKey := routing.GameLogSlug + "." + username

	//call pubsub.PublishGob
	if err := pubsub.PublishGob(
		ch,
		routing.ExchangePerilTopic,
		routingKey,
		gameLogStruct,
	); err != nil {
		return err
	}
	return nil
}