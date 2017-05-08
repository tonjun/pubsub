package handlers

import (
	"log"

	"github.com/tonjun/pubsub"
)

type SubscribeHandler struct {
}

func NewSubscribeHandler() *SubscribeHandler {
	return &SubscribeHandler{}
}

func (h *SubscribeHandler) ProcessMessage(s pubsub.Server, c pubsub.Conn, mesg *pubsub.Message) {
	log.Printf("ProcessMessage: op: %s", mesg.OP)
	if mesg.OP == pubsub.OPSubscribe {
		resp := &pubsub.Message{
			OP: "subscribe-response",
			ID: mesg.ID,
		}
		c.Send(pubsub.ToBytes(resp))
	}
}
