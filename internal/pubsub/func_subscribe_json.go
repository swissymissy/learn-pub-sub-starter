package pubsub

import (
	"fmt"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)


func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T),
) error {

	// make sure the given queue exists and bounded to exchange
	channel, queue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return fmt.Errorf("Queue does not exist or not bounded to exchange: %w", err)
	}

	// create a delivery message channel
	deliveryChan , err := channel.Consume(queue.Name, "", false, false, false, false, nil) 
	if err != nil {
		return fmt.Errorf("Error deliver queued messages: %w", err)
	}

	// goroutine to acknowledge msg and remove msg from queue in channel
	go func() {
		for eachMsg := range deliveryChan{
			// unmarshal each msg from json bytes to generic type T
			var data T
			if err = json.Unmarshal(eachMsg.Body, &data); err != nil {
				log.Printf("Error unmarshaling channel msg: %v", err)
				continue
			}
			handler(data)
			// acknowledge the message to remove it from the queue
			if err = eachMsg.Ack(false); err != nil {
				log.Printf("Error removing msg from queue: %v", err)
				continue
			}
		}
	}()

	return nil
}