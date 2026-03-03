package main

import (
	"fmt"
	"os/signal"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"

)

func main() {
	fmt.Println("Starting Peril server...")
	
	RabbitConnectStr := "amqp://guest:guest@localhost:5672/"	// where to connect to rabbitmq server
	connection, err := amqp.Dial(RabbitConnectStr)				// create new connetion with rabbitmq
	if err != nil {
		fmt.Printf("Can't connect to rabbitmq server: %s", err)
		return
	}
	defer connection.Close()									// close the connection before exiting program
	fmt.Println("Successfully connect to RabbitMQ server!")

	// create new durable queue and bind to exchange
	bindKey := routing.GameLogSlug + "." + "*"
	newChan, _, err := pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		bindKey,
		pubsub.SimpleQueueDurable,
	)
	if err != nil {
		fmt.Printf("Can't create durable queue: %s\n", err)
		return
	}

	gamelogic.PrintServerHelp()
	// RELP 
	for {
		//get user's input
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		firstWord := input[0]
		switch firstWord {
		case "pause":
			fmt.Println("Sending a pause message...")
			// publish a message to the exchange
			err = pubsub.PublishJSON(newChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			if err != nil {
				fmt.Printf("Can't publish message to the exchange: %s", err)
				continue
			}
		case "resume":
			fmt.Println("Sending resume message...")
			err = pubsub.PublishJSON(newChan, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			if err != nil {
				fmt.Printf("Can't publish message to the exchange: %s", err)
				continue
			}
		case "quit":
			fmt.Println("Exiting server...")
			return
		default: 
			fmt.Println("Unknown command.")
		}
	}

	// wait for signal ctrl+c
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel
	fmt.Println("Rabbit connection closed.")
}
