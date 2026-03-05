package pubsub

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

// consume message function JSON
func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
) error {
	// call helper func
	return helperSubscribe[T](
		conn,
		exchange,
		queueName,
		key,
		queueType,
		handler,
		func(data []byte) (T, error) {		// provide decoding logic
			var result T
			err := json.Unmarshal(data, &result)
			return result, err
		},
	)
}