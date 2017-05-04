package pubsub

type Handler interface {
	ProcessMessage(s Server, c Conn, mesg *Message)
}
