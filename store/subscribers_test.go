package store

import (
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type myConn struct {
	id uint64
}

func newMyConn(id uint64) *myConn {
	return &myConn{id}
}

func (c *myConn) ID() uint64 {
	return c.id
}

func (c *myConn) Send(d []byte) error {
	return nil
}

func (c *myConn) Close() {
}

func TestSubscribers(t *testing.T) {
	subs := NewSubscribers()
	subs.Init()
	c1 := newMyConn(1)
	c2 := newMyConn(2)
	subs.Add("t1", c1)
	subs.Add("t1", c2)
	s := subs.GetSubscribers("t1")
	assert.Equal(t, 2, len(s))
	for i, c := range s {
		log.Printf("conn id[%d]: %d %v", i, c.ID(), c)
	}
	subs.Close()
	time.Sleep(10 * time.Millisecond)
}
