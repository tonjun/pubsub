package pubsub

type Conn interface {
	ID() uint64
	Send(data []byte) error
	Close()
}
