package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
)

// writing logs to disk
func handlerLogs() func(routing.GameLog) pubsub.AckType {
	return func(gamelog routing.GameLog) pubsub.AckType{
		defer fmt.Print("> ")

		err := gamelogic.WriteLog(gamelog)
		if err != nil {
			fmt.Printf("Can't write log to disk: %s\n", err)
			return pubsub.NackRequeue
		}
		return pubsub.Ack
	}
}