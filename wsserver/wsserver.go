/*
Package websocketserver implements Server and Conn pubsub interface using websockets
*/
package websocketserver

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/tonjun/pubsub"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Options is the options for NewWebSocketServer
type Options struct {

	// ListenAddr is the listening address of the websocket server. e.g. ":7070"
	ListenAddr string

	// Pattern is the http handler pattern for websocket. e.g. "/ws"
	Pattern string

	// TLSCert is the certificate file path for listening on TLS connection
	TLSCert string

	// TLSKey is the private key file path for listening on TLS connection
	TLSKey string
}

// WebSocketServer implements Server interface in pubsub package
type WebSocketServer struct {
	opts        *Options
	cntr        int64
	cntrLck     sync.Mutex
	svr         *http.Server
	connections map[int64]*Connection
	register    chan pubsub.Conn
	unregister  chan pubsub.Conn

	connAddedFn func(c pubsub.Conn)
	connCloseFn func(c pubsub.Conn)
	onMesgFn    func(data []byte, c pubsub.Conn)
}

// NewWebSocketServer creates a new instance of WebSocketServer
func NewWebSocketServer(opts *Options) *WebSocketServer {
	return &WebSocketServer{
		opts:        opts,
		connections: make(map[int64]*Connection),
		register:    make(chan pubsub.Conn),
		unregister:  make(chan pubsub.Conn),
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	EnableCompression: true,
}

// Run runs the websocket server. This calls ListenAndServe / ListenAndServeTLS
func (s *WebSocketServer) Run() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc(s.opts.Pattern, s.serveWS).Methods("GET")

	go func() {
		for {
			select {
			case c := <-s.register:
				s.connections[c.ID()] = c.(*Connection)
				if s.connAddedFn != nil {
					s.connAddedFn(c)
				}

			case c := <-s.unregister:
				if _, ok := s.connections[c.ID()]; ok {
					if s.connCloseFn != nil {
						s.connCloseFn(c)
					}
					delete(s.connections, c.ID())
				}
			}
		}
	}()

	s.svr = &http.Server{
		Addr:    s.opts.ListenAddr,
		Handler: router,
	}
	err := s.svr.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Printf("ListenAndServe: error: %s", err.Error())
		log.Fatal(err)
	}

	//log.Fatal(http.ListenAndServe(s.opts.ListenAddr, router))
}

// Stop stops the websocket server
func (s *WebSocketServer) Stop() {
	err := s.svr.Close()
	if err != nil {
		log.Printf("Stop: ERROR: %s", err.Error())
	}
}

// OnConnectionAdded assigns the connection added callback function
func (s *WebSocketServer) OnConnectionAdded(fn func(c pubsub.Conn)) {
	s.connAddedFn = fn
}

// OnConnectionWillClose assigns the connection close callback function
func (s *WebSocketServer) OnConnectionWillClose(fn func(c pubsub.Conn)) {
	s.connCloseFn = fn
}

// OnMessage assigns the incoming message callback function
func (s *WebSocketServer) OnMessage(fn func(data []byte, c pubsub.Conn)) {
	s.onMesgFn = fn
}

func (s *WebSocketServer) serveWS(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("serveWS: ERROR: %s", err.Error())
		return
	}
	c := NewConnection(s.nextID(), ws)
	s.addConnection(c)
	c.onClose(func(c *Connection) {
		s.delConnection(c)
	})
	c.onMessage(func(c *Connection, b []byte) {
		if s.onMesgFn != nil {
			s.onMesgFn(b, c)
		}
	})
	go c.writePump()
	c.readPump()
}

func (s *WebSocketServer) nextID() int64 {
	s.cntrLck.Lock()
	defer s.cntrLck.Unlock()
	s.cntr++
	return s.cntr
}

func (s *WebSocketServer) addConnection(c pubsub.Conn) {
	select {
	case s.register <- c:
	case <-time.After(1 * time.Second):
		log.Printf("addConnection: ERROR: timeout")
	}
}

func (s *WebSocketServer) delConnection(c pubsub.Conn) {
	select {
	case s.unregister <- c:
	case <-time.After(1 * time.Second):
		log.Printf("delConnection: ERROR: timeout")
	}
}