package main_test

import (
	"log"
	"time"

	"github.com/onsi/gomega/gbytes"
	"github.com/tonjun/pubsub"
	"github.com/tonjun/wsclient"

	. "github.com/tonjun/pubsub/pubsub"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Publish", func() {

	var (
		server *PubSubServer

		client1 *wsclient.WSClient
		client2 *wsclient.WSClient

		buffer1 *gbytes.Buffer
		buffer2 *gbytes.Buffer
	)

	BeforeEach(func() {
		buffer1 = gbytes.NewBuffer()
		buffer2 = gbytes.NewBuffer()

		connected := make(chan bool)

		// run server on port 7070
		cfg := &pubsub.Config{
			Addr: ":7070",
			Path: "/ws",
		}
		server = NewPubSubServer(cfg)
		go server.Main()

		// connect a client1 and write all incoming message to gbytes.Buffer
		client1 = wsclient.NewWSClient("ws://localhost:7070/ws")
		client1.OnMessage(func(data []byte) {
			buffer1.Write(data)
		})
		client1.OnOpen(func() {
			connected <- true
		})
		client1.OnError(func(err error) {
			log.Printf("reconnecting..")
			time.Sleep(10 * time.Millisecond)
			client1.Connect()
		})
		client1.Connect()
		<-connected

		// connect a client2 and write all incoming message to gbytes.Buffer
		client2 = wsclient.NewWSClient("ws://localhost:7070/ws")
		client2.OnMessage(func(data []byte) {
			buffer2.Write(data)
		})
		client2.OnOpen(func() {
			connected <- true
		})
		client2.OnError(func(err error) {
			log.Printf("reconnecting..")
			time.Sleep(10 * time.Millisecond)
			client2.Connect()
		})
		client2.Connect()
		<-connected

	})

	AfterEach(func() {
		server.Close()
		//time.Sleep(1000 * time.Millisecond)
	})

	It("Publish should send the message to all the subscribers", func(done Done) {

		client1.SendJSON(wsclient.M{
			"op":     "subscribe",
			"id":     "req1",
			"topics": []string{"t1"},
		})
		client2.SendJSON(wsclient.M{
			"op":     "subscribe",
			"id":     "req1",
			"topics": []string{"t1"},
		})
		Eventually(buffer1).Should(gbytes.Say(
			`{"op":"subscribe_response","id":"req1","data":{"type":"success"}}`,
		))
		Eventually(buffer2).Should(gbytes.Say(
			`{"op":"subscribe_response","id":"req1","data":{"type":"success"}}`,
		))

		client1.SendJSON(wsclient.M{
			"op":     "publish",
			"id":     "pub1",
			"topics": []string{"t1"},
			"data": map[string]string{
				"body": "hi all",
			},
		})

		Eventually(buffer1).Should(gbytes.Say(
			`{"op":"publish_response","id":"pub1","data":{"type":"success"}}`,
		))
		Eventually(buffer1).Should(gbytes.Say(
			`{"op":"publish","id":"pub1","topics":\["t1"\],"data":\{"body":"hi all"\}}`,
		))
		Eventually(buffer2).Should(gbytes.Say(
			`{"op":"publish","id":"pub1","topics":\["t1"\],"data":\{"body":"hi all"\}}`,
		))

		close(done)
	})

})
