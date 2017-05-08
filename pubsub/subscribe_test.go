package main_test

import (
	"log"
	"time"

	"github.com/tonjun/pubsub"
	"github.com/tonjun/wsclient"

	. "github.com/tonjun/pubsub/pubsub"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Subscribe", func() {

	var (
		server *PubSubServer
		client *wsclient.WSClient
		buffer *gbytes.Buffer
	)

	BeforeEach(func() {
		log.Printf("BeforeEach")
		buffer = gbytes.NewBuffer()
		connected := make(chan bool)

		// run server on port 7070
		cfg := &pubsub.Config{
			Addr: ":7070",
			Path: "/ws",
		}
		server = NewPubSubServer(cfg)
		go server.Main()

		// connect a client and write all incoming message to gbytes.Buffer
		client = wsclient.NewWSClient("ws://localhost:7070/ws")
		client.OnMessage(func(data []byte) {
			buffer.Write(data)
		})
		client.OnOpen(func() {
			connected <- true
		})
		client.OnError(func(err error) {
			log.Printf("reconnecting..")
			time.Sleep(10 * time.Millisecond)
			client.Connect()
		})
		client.Connect()

		<-connected
	})

	AfterEach(func() {
		log.Printf("AfterEach")
		server.Close()
	})

	It("Subscribe should get a successful response", func() {
		client.SendJSON(wsclient.M{
			"op":     "subscribe",
			"id":     "req1",
			"topics": []string{"t1"},
		})
		Eventually(buffer).Should(gbytes.Say(`{"op":"subscribe-response","id":"req1"}`))

		client.SendJSON(wsclient.M{
			"op":     "subscribe",
			"id":     "req2",
			"topics": []string{"t1", "t2", "t3"},
		})
		Eventually(buffer).Should(gbytes.Say(`{"op":"subscribe-response","id":"req2"}`))

	})

})
