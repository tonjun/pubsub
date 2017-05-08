package handlers

import (
	"log"

	"github.com/tonjun/pubsub"
	"github.com/tonjun/pubsub/store"
)

type SubscribeHandler struct {
	cfg         *pubsub.Config
	subscribers *store.Subscribers
}

func NewSubscribeHandler(cfg *pubsub.Config, s *store.Subscribers) *SubscribeHandler {
	return &SubscribeHandler{
		cfg:         cfg,
		subscribers: s,
	}
}

func (h *SubscribeHandler) ProcessMessage(s pubsub.Server, c pubsub.Conn, mesg *pubsub.Message) {
	if mesg.OP == pubsub.OPSubscribe {
		log.Printf("ProcessMessage: op: %s", mesg.OP)

		for _, topic := range mesg.Topics {
			h.subscribers.Add(topic, c)
		}

		resp := &pubsub.Message{
			OP: pubsub.OPSubscribeResponse,
			ID: mesg.ID,
			Data: map[string]string{
				"type": "success",
			},
		}
		c.Send(pubsub.ToBytes(resp))
	}
}
