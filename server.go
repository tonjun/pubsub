package pubsub

type Server interface {
	Run()
	Stop()

	OnConnectionAdded(fn func(conn Conn))
	OnConnectionWillClose(fn func(conn Conn))
	OnMessage(fn func(data []byte, conn Conn))
}
