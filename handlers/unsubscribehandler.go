package handlers

import (
	"log"

	"github.com/tonjun/pubsub"
)

type UnsubscribeHandler struct {
}

func NewUnsubscribeHandler() *UnsubscribeHandler {
	return &UnsubscribeHandler{}
}

func (h *UnsubscribeHandler) ProcessMessage(s pubsub.Server, c pubsub.Conn, mesg *pubsub.Message) {
	log.Printf("ProcessMessage: op: %s", mesg.OP)
}
