package websocketserver

import (
	"log"

	"github.com/gorilla/websocket"
)

type Connection struct {
	id int64
	ws *websocket.Conn

	onCloseCb func(c *Connection)
	onMesgCb  func(c *Connection, b []byte)
}

func NewConnection(id int64, ws *websocket.Conn) *Connection {
	return &Connection{
		id: id,
		ws: ws,
	}
}

func (c *Connection) ID() int64 {
	return c.id
}

func (c *Connection) Send(data []byte) error {
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

func (c *Connection) writePump() {
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
