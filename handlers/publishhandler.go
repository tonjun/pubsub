package handlers

import (
	"log"

	"github.com/tonjun/pubsub"
)

type PublishHandler struct {
}

func NewPublishHandler() *PublishHandler {
	return &PublishHandler{}
}

func (h *PublishHandler) ProcessMessage(s pubsub.Server, c pubsub.Conn, mesg *pubsub.Message) {
	log.Printf("ProcessMessage: op: %s", mesg.OP)
}
