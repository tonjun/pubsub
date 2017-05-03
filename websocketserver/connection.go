package websocketserver

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	id       int64
	ws       *websocket.Conn
	send     chan sendRequest
	closed   bool
	closedMu sync.RWMutex

	onCloseCb func(c *Connection)
	onMesgCb  func(c *Connection, b []byte)
}

type sendRequest struct {
	data  []byte
	errCh chan error
}

func NewConnection(id int64, ws *websocket.Conn) *Connection {
	return &Connection{
		id:   id,
		ws:   ws,
		send: make(chan sendRequest),
	}
}

func (c *Connection) ID() int64 {
	return c.id
}

func (c *Connection) Send(data []byte) error {
	log.Printf("Send: \"%s\"", string(data))
	if c.isClosed() {
		return fmt.Errorf("connection closed")
	}
	req := &sendRequest{
		data:  data,
		errCh: make(chan error),
	}
	select {
	case c.send <- *req:
	case <-time.After(1 * time.Second):
		log.Printf("Send: timeout!!")
		return fmt.Errorf("Timeout while sending to channel")
	}
	err := <-req.errCh
	return err
}

func (c *Connection) Close() {
	log.Printf("Close")
	defer log.Printf("Close done")
	if c.isClosed() {
		log.Printf("Close: already closed")
		return
	}
	c.closedMu.Lock()
	c.closed = true
	c.closedMu.Unlock()
	c.ws.Close()
	c.onCloseCb(c)
	close(c.send)
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
		c.Close()
		//close(c.send)
	}()
	for {
		select {
		case req, ok := <-c.send:
			if !ok {
				log.Printf("writePump: send channel !ok")
				c.write(websocket.CloseMessage, []byte{})
				if req.errCh != nil {
					req.errCh <- fmt.Errorf("send channel !ok")
				} else {
					log.Printf("req.errCh is nil")
				}
				return
			}
			if err := c.write(websocket.TextMessage, req.data); err != nil {
				log.Printf("writePump: WriteMessage error: %s", err.Error())
				req.errCh <- err
				return
			}
			req.errCh <- nil
		}
	}
}

func (c *Connection) readPump() {
	defer func() {
		//c.ws.Close()
		//c.onCloseCb(c)
		c.Close()
		log.Printf("readPump: done")
	}()
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
}

func (c *Connection) isClosed() bool {
	c.closedMu.RLock()
	defer c.closedMu.RUnlock()
	return c.closed
}
