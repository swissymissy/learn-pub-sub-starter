package main

import (
	"fmt"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"

)

func main() {
	fmt.Println("Starting Peril server...")
	
	RabbitConnectStr := "amqp://guest:guest@localhost:5672/"		// where to connect to rabbitmq server
	connection, err := amqp.Dial(connectStr)				// create new connetion with rabbitmq
	if err != nil {
		fmt.Println("Can't connect to rabbitmq server")
		return
	}
	defer connection.Close()		// close the connection before exit program
	fmt.Println("Successfully connect to RabbitMQ server!")

	// wait for signal ctrl+c
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel
	fmt.Println("Rabbit connection closed.")
}
