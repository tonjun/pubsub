package handlers

import (
	"log"

	"github.com/tonjun/pubsub"
	"github.com/tonjun/pubsub/store"
)

type PublishHandler struct {
	cfg         *pubsub.Config
	subscribers *store.Subscribers
}

func NewPublishHandler(c *pubsub.Config, s *store.Subscribers) *PublishHandler {
	return &PublishHandler{
		cfg:         c,
		subscribers: s,
	}
}

func (h *PublishHandler) ProcessMessage(s pubsub.Server, c pubsub.Conn, mesg *pubsub.Message) {

	if mesg.OP == pubsub.OPPublish {
		log.Printf("ProcessMessage: op: %s", mesg.OP)
		resp := &pubsub.Message{
			OP: pubsub.OPPublishResponse,
			ID: mesg.ID,
			Data: map[string]string{
				"type": "success",
			},
		}
		c.Send(pubsub.ToBytes(resp))

		for _, topic := range mesg.Topics {
			subs := h.subscribers.GetSubscribers(topic)
			for _, s := range subs {
				s.Send(pubsub.ToBytes(mesg))
			}
		}

	}

}
