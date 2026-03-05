package pubsub

import (
	"fmt"
	"encoding/gob"
	"bytes"
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// publish the message in game_logs queue
func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	// convert val to gob bytes
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)			// write to buffer
	err := enc.Encode(val)				// encode val T into bytes
	if err != nil {
		return err
	}
	result := buf.Bytes()				// extract the bytes from buffer

	// publish message 
	if err = ch.PublishWithContext(
		context.Background(),
		exchange,
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/gob",
			Body: result,
		},
	); err != nil {
		return fmt.Errorf("Can't publish message: %w", err)
	}
	return nil
}