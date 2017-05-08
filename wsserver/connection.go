package wsserver

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Connection implements Conn interface in pubsub package
type Connection struct {
	id       uint64
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

// NewConnection creates a new instance of Connection
func NewConnection(id uint64, ws *websocket.Conn) *Connection {
	return &Connection{
		id:   id,
		ws:   ws,
		send: make(chan sendRequest),
	}
}

// ID returns the connection ID
func (c *Connection) ID() uint64 {
	return c.id
}

// Send sends the data to the websocket
func (c *Connection) Send(data []byte) error {
	log.Printf("%d Send: \"%s\"", c.id, string(data))
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

// Close closes the connection
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

func (c *Connection) onClose(fn func(c *Connection)) {
	c.onCloseCb = fn
}

func (c *Connection) onMessage(fn func(c *Connection, b []byte)) {
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
