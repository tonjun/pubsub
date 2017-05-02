package pubsub

type Conn interface {
	ID() int64
	Send(data []byte) error
	Close()
}
