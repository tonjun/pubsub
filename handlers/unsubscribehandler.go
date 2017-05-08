package handlers

import (
	//"log"

	"github.com/tonjun/pubsub"
	"github.com/tonjun/pubsub/store"
)

type UnsubscribeHandler struct {
	cfg         *pubsub.Config
	subscribers *store.Subscribers
}

func NewUnsubscribeHandler(c *pubsub.Config, s *store.Subscribers) *UnsubscribeHandler {
	return &UnsubscribeHandler{
		cfg:         c,
		subscribers: s,
	}
}

func (h *UnsubscribeHandler) ProcessMessage(s pubsub.Server, c pubsub.Conn, mesg *pubsub.Message) {
	//log.Printf("ProcessMessage: op: %s", mesg.OP)
}
