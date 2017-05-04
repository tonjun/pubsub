package handlers

import (
	"log"

	"github.com/tonjun/pubsub"
)

type ConnectHandler struct {
}

func NewConnectHandler() *ConnectHandler {
	return &ConnectHandler{}
}

func (h *ConnectHandler) ProcessMessage(s pubsub.Server, c pubsub.Conn, mesg *pubsub.Message) {
	log.Printf("ProcessMessage: op: %s", mesg.OP)
}
