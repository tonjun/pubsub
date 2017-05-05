package store

import (
	//"log"
	"testing"

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
	assert.Equal(t, uint64(1), s[0].ID())
	assert.Equal(t, uint64(2), s[1].ID())
	//for i, c := range s {
	//	log.Printf("conn id[%d]: %d %v", i, c.ID(), c)
	//}
	subs.Remove("t1", c1)

	s = subs.GetSubscribers("t1")
	assert.Equal(t, 1, len(s))
	assert.Equal(t, uint64(2), s[0].ID())

	c3 := newMyConn(3)
	subs.Add("t1", c3)
	subs.Add("t2", c3)
	subs.Add("t3", c3)
	subs.Add("t3", c1)

	s = subs.GetSubscribers("t1")
	assert.Equal(t, 2, len(s))
	assert.Equal(t, uint64(2), s[0].ID())
	assert.Equal(t, uint64(3), s[1].ID())

	s = subs.GetSubscribers("t2")
	assert.Equal(t, 1, len(s))
	assert.Equal(t, uint64(3), s[0].ID())

	s = subs.GetSubscribers("t3")
	assert.Equal(t, 2, len(s))
	assert.Equal(t, uint64(1), s[0].ID())
	assert.Equal(t, uint64(3), s[1].ID())

	subs.Close()
}
