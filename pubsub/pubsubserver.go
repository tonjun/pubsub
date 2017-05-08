package main

import (
	"encoding/json"
	"log"

	"github.com/tonjun/pubsub"
	"github.com/tonjun/pubsub/handlers"
	"github.com/tonjun/pubsub/store"
	"github.com/tonjun/pubsub/wsserver"
)

type PubSubServer struct {
	srv         *wsserver.WSServer
	cfg         *pubsub.Config
	handlers    []pubsub.Handler
	subscribers *store.Subscribers
}

func NewPubSubServer(cfg *pubsub.Config) *PubSubServer {
	opts := &wsserver.Options{
		ListenAddr: cfg.Addr,
		Path:       cfg.Path,
	}
	return &PubSubServer{
		srv:         wsserver.NewWSServer(opts),
		cfg:         cfg,
		handlers:    make([]pubsub.Handler, 0),
		subscribers: store.NewSubscribers(),
	}
}

func (ps *PubSubServer) Main() {

	var h pubsub.Handler

	ps.subscribers.Init()

	h = handlers.NewConnectHandler()
	ps.handlers = append(ps.handlers, h)

	h = handlers.NewSubscribeHandler(ps.cfg, ps.subscribers)
	ps.handlers = append(ps.handlers, h)

	h = handlers.NewPublishHandler(ps.cfg, ps.subscribers)
	ps.handlers = append(ps.handlers, h)

	h = handlers.NewUnsubscribeHandler(ps.cfg, ps.subscribers)
	ps.handlers = append(ps.handlers, h)

	ps.srv.OnMessage(ps.onMessage)
	ps.srv.Run()
}

func (ps *PubSubServer) Close() {
	ps.srv.Stop()
}

func (ps *PubSubServer) onMessage(data []byte, c pubsub.Conn) {

	//go func() {
	mesg := &pubsub.Message{}
	err := json.Unmarshal(data, mesg)
	if err != nil {
		log.Printf("WARNING: Invalid packet from client")
		log.Printf("WARNING: \"%s\"", string(data))
		c.Close()
		return
	}
	for _, h := range ps.handlers {
		h.ProcessMessage(ps.srv, c, mesg)
	}
	//}()
}
