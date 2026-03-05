package main

import (
	"fmt"
	"os/signal"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func main() {
	fmt.Println("Starting Peril client...")

	RabbitConnectStr := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(RabbitConnectStr)
	if err != nil {
		fmt.Printf("Can't connect to rabbitmq server: %s", err)
		return
	}
	defer connection.Close()
	fmt.Println("Successfully connect to rabbitmq server!")

	// create a channel
	pubChan, err := connection.Channel()
	if err != nil {
		fmt.Printf("Can't create new channel: %s", err)
		return 
	}

	// prompt user for a username
	usrname, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Printf("Can't get username: %s", err)
		return
	}
	
	// create new game state of user
	gameState := gamelogic.NewGameState(usrname)
	pauseQueueName := "pause." + usrname
	
	// start a routine to listen to pause msges in background
	if err = pubsub.SubscribeJSON(
		connection,
		routing.ExchangePerilDirect,
		pauseQueueName,
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(gameState),
	); err != nil {
		fmt.Printf("Error subscribing pause: %s", err)
		return 
	}

	// start routine to listen to army move msges in background
	moveQueueName := "army_moves." + usrname
	moveRoutingKey := routing.ArmyMovesPrefix + "." + "*"
	if err = pubsub.SubscribeJSON(
		connection,
		routing.ExchangePerilTopic,
		moveQueueName,
		moveRoutingKey,
		pubsub.SimpleQueueTransient,
		handlerMove(gameState, pubChan),
	); err != nil {
		fmt.Printf("Error subscribing army_moves: %s", err)
		return
	}

	// start routine to listen to war message
	warQueueName := routing.WarRecognitionsPrefix
	warRoutingKey := routing.WarRecognitionsPrefix + ".#"
	if err = pubsub.SubscribeJSON(
		connection,
		routing.ExchangePerilTopic,
		warQueueName,
		warRoutingKey,
		pubsub.SimpleQueueDurable,
		handlerWar(gameState),
	); err != nil {
		fmt.Printf("Error subscribing war message: %s\n", err)
		return
	}

	// REPL
	for {											
		// get user's input
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}

		firstWord := input[0]

		switch firstWord {
		case "spawn":
			err = gameState.CommandSpawn(input)
			if err != nil {
				fmt.Printf("%s\n", err)
				continue
			}
		case "move":
			mv, err := gameState.CommandMove(input)
			if err != nil {
				fmt.Printf("%s\n", err)
				continue
			}
			//publish the move
			if err = pubsub.PublishJSON(pubChan, routing.ExchangePerilTopic, moveQueueName, mv); err != nil {
				fmt.Printf("%s\n", err)
				continue
			}
			fmt.Println("Units have successfully moved to new location!")
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet")
		case "quit":
			gamelogic.PrintQuit()
			return 
		default:
			fmt.Println("Unknown command")
			continue
		}
	}

	// wait for signal ctrl+c
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel
	fmt.Println("Rabbit connection closed.")
}
