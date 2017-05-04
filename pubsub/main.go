package main

import ()

func main() {
	srv := NewPubSubServer(nil)
	srv.Main()
}
