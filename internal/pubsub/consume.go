package pubsub


type SimpleQueueType int 

const (
	SimpleQueueDurable SimpleQueueType =  iota
	SimpleQueueTransient
)

type AckType int 

const (
	Ack AckType = iota 
	NackRequeue
	NackDiscard
)