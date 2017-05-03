package websocketserver

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	id   int64
	ws   *websocket.Conn
	send chan []byte

	onCloseCb func(c *Connection)
	onMesgCb  func(c *Connection, b []byte)
}

func NewConnection(id int64, ws *websocket.Conn) *Connection {
	return &Connection{
		id:   id,
		ws:   ws,
		send: make(chan []byte),
	}
}

func (c *Connection) ID() int64 {
	return c.id
}

func (c *Connection) Send(data []byte) error {
	log.Printf("Send: \"%s\"", string(data))

	select {
	case c.send <- data:
	case <-time.After(1 * time.Second):
		log.Printf("Send: timeout!!")
		return fmt.Errorf("Timeout while sending to channel")
	}

	return nil
}

func (c *Connection) Close() {
}

func (c *Connection) OnClose(fn func(c *Connection)) {
	c.onCloseCb = fn
}

func (c *Connection) OnMessage(fn func(c *Connection, b []byte)) {
	c.onMesgCb = fn
}

func (c *Connection) write(mt int, payload []byte) error {
	return c.ws.WriteMessage(mt, payload)
}

func (c *Connection) writePump() {
	defer func() {
		log.Printf("writePump: done")
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				log.Printf("writePump: send channel !ok")
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				log.Printf("WriteMessage error: %s", err.Error())
				return
			}
		}
	}
}

func (c *Connection) readPump() {
	for {
		t, message, err := c.ws.ReadMessage()
		if err != nil {
			log.Printf("read error: %s", err.Error())
			break
		}
		log.Printf("readPump: messageType: %d message: %s", t, message)
		if c.onMesgCb != nil {
			c.onMesgCb(c, message)
		}
	}
	c.ws.Close()
	c.onCloseCb(c)
}
