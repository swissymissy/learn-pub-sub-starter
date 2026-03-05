package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// subscribe helper
func helperSubscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
	decoder func([]byte) (T , error),
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
		defer channel.Close()
		for eachMsg := range deliveryChan{
			// decode each msg from bytes to generic type T
			var data T
			data, err = decoder(eachMsg.Body)
			if err != nil {
				fmt.Printf("Error decoding message: %s\n", err)
				continue 
			}
			ackType := handler(data)
			switch ackType {
			case Ack:
				fmt.Println("Ack type called")
				eachMsg.Ack(false)
			case NackRequeue:
				fmt.Println("NackRequeue type called")
				eachMsg.Nack(false, true)
			case NackDiscard:
				fmt.Println("Nack Discard type called")
				eachMsg.Nack(false, false)
			}
		}
	}()
	return nil 
}