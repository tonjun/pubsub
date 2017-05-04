package websocketserver

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/tonjun/pubsub"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestConnectionCallbacks(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	var srv pubsub.Server
	addr := ":7070"
	srv = NewWebSocketServer(&Options{
		ListenAddr: addr,
		Pattern:    "/ws",
	})
	assert.NotNil(t, srv)

	added := make(chan pubsub.Conn)
	closed := make(chan pubsub.Conn)
	connections := make(map[int64]pubsub.Conn)
	check1 := make(chan bool)
	check2 := make(chan bool)
	check3 := make(chan bool)

	go func() {
		for {
			select {
			case c := <-added:
				connections[c.ID()] = c

			case c := <-closed:
				delete(connections, c.ID())

			case <-check1:
				assert.Equal(t, 1, len(connections))
				conn := connections[1]
				assert.Equal(t, int64(1), conn.ID())

			case <-check2:
				assert.Equal(t, 2, len(connections))
				conn := connections[1]
				assert.Equal(t, int64(1), conn.ID())

				conn = connections[2]
				assert.Equal(t, int64(2), conn.ID())

			case <-check3:
				assert.Equal(t, 1, len(connections))

			}
		}
	}()

	go srv.Run()

	// assign callbacks
	srv.OnConnectionAdded(func(conn pubsub.Conn) {
		added <- conn
	})
	srv.OnConnectionWillClose(func(conn pubsub.Conn) {
		closed <- conn
	})

	time.Sleep(10 * time.Millisecond)

	var conn1 *websocket.Conn
	var conn2 *websocket.Conn
	var resp *http.Response
	var err error
	conn1, resp, err = websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost%s/ws", addr), nil)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, conn1)

	time.Sleep(10 * time.Millisecond)
	check1 <- true

	conn2, resp, err = websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost%s/ws", addr), nil)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, conn2)

	time.Sleep(10 * time.Millisecond)
	check2 <- true

	conn1.Close()

	time.Sleep(10 * time.Millisecond)
	check3 <- true

	srv.Stop()

	time.Sleep(100 * time.Millisecond)
}

func TestMessageCallback(t *testing.T) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	var srv pubsub.Server

	addr := ":7070"

	srv = NewWebSocketServer(&Options{
		ListenAddr: addr,
		Pattern:    "/ws",
	})
	assert.NotNil(t, srv)

	type M struct {
		Message []byte
		Conn    pubsub.Conn
	}

	//messages := make([]M, 0)
	var messages []M
	check1 := make(chan bool)
	check2 := make(chan bool)
	check3 := make(chan bool)
	addMesg := make(chan M)

	go func() {
		for {
			select {
			case m := <-addMesg:
				messages = append(messages, m)

			case <-check1:
				assert.Equal(t, 1, len(messages))
				assert.Equal(t, []byte("one"), messages[0].Message)
				assert.Equal(t, int64(1), messages[0].Conn.ID())

			case <-check2:
				assert.Equal(t, 3, len(messages))
				assert.Equal(t, []byte("one"), messages[0].Message)
				assert.Equal(t, []byte("two"), messages[1].Message)
				assert.Equal(t, []byte("three"), messages[2].Message)
				for _, m := range messages {
					assert.Equal(t, int64(1), m.Conn.ID())
				}

			case <-check3:
				assert.Equal(t, 4, len(messages))
				assert.Equal(t, []byte("four"), messages[3].Message)
				assert.Equal(t, int64(2), messages[3].Conn.ID())

			}
		}
	}()

	go srv.Run()

	srv.OnMessage(func(data []byte, c pubsub.Conn) {
		log.Printf("onMessage: \"%s\"", string(data))
		m := M{data, c}
		addMesg <- m
	})

	time.Sleep(10 * time.Millisecond)

	var conn1 *websocket.Conn
	var conn2 *websocket.Conn
	var resp *http.Response
	var err error
	conn1, resp, err = websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost%s/ws", addr), nil)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, conn1)

	err = conn1.WriteMessage(websocket.TextMessage, []byte("one"))
	assert.Nil(t, err)

	time.Sleep(10 * time.Millisecond)
	check1 <- true

	conn1.WriteMessage(websocket.TextMessage, []byte("two"))
	conn1.WriteMessage(websocket.TextMessage, []byte("three"))

	time.Sleep(10 * time.Millisecond)
	check2 <- true

	conn2, resp, err = websocket.DefaultDialer.Dial(fmt.Sprintf("ws://localhost%s/ws", addr), nil)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, conn2)
	conn2.WriteMessage(websocket.TextMessage, []byte("four"))

	time.Sleep(10 * time.Millisecond)
	check3 <- true

	srv.Stop()

	time.Sleep(100 * time.Millisecond)
}
