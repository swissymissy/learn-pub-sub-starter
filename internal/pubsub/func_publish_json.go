package pubsub

import (
	"encoding/json"
	"fmt"
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)


func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	// convert val to json byte
	jsonByte, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("Can't convert val to json bytes: %w", err)
	}

	// publish the message to the exchange with routing key
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{ContentType: "application/json", Body: jsonByte})
	if err != nil {
		return fmt.Errorf("Can't publish message to the exchange: %w", err)
	}

	return nil
}