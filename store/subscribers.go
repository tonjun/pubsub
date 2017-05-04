package store

import (
	"log"
	//"time"

	"github.com/google/btree"
	"github.com/tonjun/pubsub"
)

// Subscribers is responsible for storing the list of subscribers per topic
type Subscribers struct {
	subs map[string]*btree.BTree
	add  chan req
	del  chan req
	get  chan getReq
}

type req struct {
	topic string
	conn  pubsub.Conn
}

type getReq struct {
	topic string
	resp  chan []pubsub.Conn
}

// NewSubscribers returns a new instance of Subscribers store
func NewSubscribers() *Subscribers {
	return &Subscribers{
		subs: make(map[string]*btree.BTree),
		add:  make(chan req),
		del:  make(chan req),
		get:  make(chan getReq),
	}
}

// Init initializes the go routine that handles requests
func (s *Subscribers) Init() {
	go s.run()
}

// Close closes the Subscribers
func (s *Subscribers) Close() {
	close(s.add)
}

// Add adds the given connection to the list of subscribers to the topic
func (s *Subscribers) Add(topic string, c pubsub.Conn) {
	r := req{
		topic: topic,
		conn:  c,
	}
	s.add <- r
}

// Remove removes the given connection from the list of subscribers of the topic
func (s *Subscribers) Remove(topic string, c pubsub.Conn) {
	r := req{
		topic: topic,
		conn:  c,
	}
	s.del <- r
}

// GetSubscribers returns the list of subscribers to the given topic
func (s *Subscribers) GetSubscribers(topic string) []pubsub.Conn {
	ch := make(chan []pubsub.Conn)
	r := getReq{
		topic: topic,
		resp:  ch,
	}
	s.get <- r
	subs := <-ch
	return subs
}

func (s *Subscribers) getTree(topic string) *btree.BTree {
	var tree *btree.BTree
	if t, ok := s.subs[topic]; !ok {
		tree = btree.New(32)
		s.subs[topic] = tree
	} else {
		tree = t
	}
	return tree
}

func (s *Subscribers) run() {
	defer log.Printf("run done")
	for {
		select {
		case r, ok := <-s.add:
			if !ok {
				log.Printf("add channel closed")
				return
			}
			//log.Printf("add: topic: \"%s\" connID: %d", r.topic, r.conn.ID())
			i := treeItem{
				Key:   r.conn.ID(),
				Value: r.conn,
			}
			s.getTree(r.topic).ReplaceOrInsert(i)

		case r := <-s.del:
			//log.Printf("del: topic: \"%s\" connID: %d", r.topic, r.conn.ID())
			i := treeItem{
				Key:   r.conn.ID(),
				Value: r.conn,
			}
			s.getTree(r.topic).Delete(i)

		case r := <-s.get:
			//log.Printf("get: topic: \"%s\"", r.topic)
			pivot := treeItem{Key: 0}
			conns := make([]pubsub.Conn, 0)
			s.getTree(r.topic).AscendGreaterOrEqual(pivot, func(a btree.Item) bool {
				conns = append(conns, a.(treeItem).Value)
				return true
			})
			r.resp <- conns

		}
	}
}
