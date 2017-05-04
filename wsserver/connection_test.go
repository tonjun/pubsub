package wsserver

import (
	"fmt"
	"log"
	//"net/http"
	//"net/http/httputil"
	"testing"
	"time"

	"github.com/tonjun/pubsub"
	"github.com/tonjun/wstester"

	//"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestConnectionSend(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// manage connections
	added := make(chan pubsub.Conn)
	closed := make(chan pubsub.Conn)
	send := make(chan string)
	connections := make(map[uint64]pubsub.Conn)
	go func() {
		for {
			select {
			case c := <-added:
				connections[c.ID()] = c
			case c := <-closed:
				delete(connections, c.ID())

			case m := <-send:
				for _, c := range connections {
					// send 10 messages
					for i := 0; i < 10; i++ {
						go func(con pubsub.Conn) {
							err := con.Send([]byte(m))
							assert.Nil(t, err)
						}(c)
					}
				}
			}
		}
	}()

	// create server and add callbacks
	var srv pubsub.Server
	addr := ":7070"
	srv = NewWSServer(&Options{
		ListenAddr: addr,
		Path:       "/ws",
	})
	assert.NotNil(t, srv)
	srv.OnConnectionAdded(func(conn pubsub.Conn) {
		added <- conn
	})
	srv.OnConnectionWillClose(func(conn pubsub.Conn) {
		closed <- conn
	})
	go srv.Run()
	time.Sleep(10 * time.Millisecond)

	// connect to the server
	wt := wstester.NewWSTester(&wstester.Options{
		ServerAddr: fmt.Sprintf("ws://localhost%s/ws", addr),
		Count:      10,
	})
	wt.Start()

	send <- "hello world"

	time.Sleep(100 * time.Millisecond)

	conns, _ := wt.GetConnections()
	for _, c := range conns {
		assert.Equal(t, 10, len(c.Messages))
		for _, m := range c.Messages {
			//log.Printf("%d messageList: %s", i, string(c.Messages[0]))
			assert.Equal(t, "hello world", string(m))
		}
	}

	time.Sleep(100 * time.Millisecond)

	// stop the server
	srv.Stop()
	time.Sleep(100 * time.Millisecond)

	wt.Stop()
	time.Sleep(100 * time.Millisecond)

}

func TestConnectionSendOnClosed(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// manage connections
	added := make(chan pubsub.Conn)
	closed := make(chan pubsub.Conn)
	send := make(chan string)
	connections := make(map[uint64]pubsub.Conn)
	go func() {
		for {
			select {
			case c := <-added:
				connections[c.ID()] = c
			case c := <-closed:
				log.Printf("deleting connection id: %d", c.ID())
				delete(connections, c.ID())

			case m := <-send:
				for _, c := range connections {
					for i := 0; i < 10; i++ {
						go func(con pubsub.Conn) {
							err := con.Send([]byte(m))
							if err != nil {
								log.Printf("TEST Send errror: %s", err.Error())
							}
							//assert.Nil(t, err)
						}(c)
					}
				}
			}
		}
	}()

	// create server and add callbacks
	var srv pubsub.Server
	addr := ":7070"
	srv = NewWSServer(&Options{
		ListenAddr: addr,
		Path:       "/ws",
	})
	assert.NotNil(t, srv)
	srv.OnConnectionAdded(func(conn pubsub.Conn) {
		added <- conn
	})
	srv.OnConnectionWillClose(func(conn pubsub.Conn) {
		//closed <- conn
	})
	go srv.Run()
	time.Sleep(10 * time.Millisecond)

	// connect to the server
	wt := wstester.NewWSTester(&wstester.Options{
		ServerAddr: fmt.Sprintf("ws://localhost%s/ws", addr),
		Count:      1,
	})
	wt.Start()

	send <- "hello world Closed socket"

	time.Sleep(100 * time.Millisecond)

	conns, _ := wt.GetConnections()
	for _, c := range conns {
		assert.Equal(t, 10, len(c.Messages))
		for _, m := range c.Messages {
			//log.Printf("%d messageList: %s", i, string(c.Messages[0]))
			assert.Equal(t, "hello world Closed socket", string(m))
		}
	}
	conns[0].WS.Close()

	time.Sleep(100 * time.Millisecond)

	send <- "hello world Closed socket"

	// stop the server
	srv.Stop()
	time.Sleep(100 * time.Millisecond)

	wt.Stop()
	time.Sleep(100 * time.Millisecond)

}
