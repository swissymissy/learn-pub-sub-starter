package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// declare and bind a queue
func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName, 
	key string,
	queueType SimpleQueueType,
) (*amqp.Channel, amqp.Queue, error) {
	
	// create new channel 
	newChan, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("Can't create new channel: %s", err)
	}

	// declare new queue
	durable := false
	autoDelete := false
	exclusive := false
	if queueType == SimpleQueueDurable {
		durable =  true
	} else {
		autoDelete = true
		exclusive = true
	}

	// create new queue: durable or transient
	newQueue, err := newChan.QueueDeclare(queueName, durable, autoDelete, exclusive, false , nil)
	if err != nil {
		return nil, amqp.Queue{} , fmt.Errorf("Can't declare new queue:%s", err) 
	}
	
	// bind the queue to the exchange
	err = newChan.QueueBind(queueName, key, exchange, false, nil )
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("Can't bind queue: %s", err)
	}

	return newChan, newQueue, nil
}