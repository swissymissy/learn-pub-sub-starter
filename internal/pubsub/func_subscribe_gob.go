package pubsub

import (
	"encoding/gob"
	"bytes"

	amqp "github.com/rabbitmq/amqp091-go"
)

// consumer game log message in gob
func SubscribeGob[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
) error {
	// call helper function
	return helperSubscribe[T](
		conn,
		exchange,
		queueName,
		key,
		queueType,
		handler,
		func(data []byte) (T, error) {
			buf := bytes.NewBuffer(data) 	// wrap byte slice in a buffer
			dec := gob.NewDecoder(buf)		// attach decoder to buffer
			var result T 					// result generic type
			err := dec.Decode(&result)		// write the result to "result"
			return result, err
		},
	)
}