package pubsub

import (
	"encoding/json"
	"os"
)

// Config is the pubsub server's configuration
type Config struct {
	Addr string     `json:"addr"` // TCP address to listen on
	Path string     `json:"path"` // websocket path. e.g. "/ws"
	TLS  *TLSConfig `json:"tls"`  // optional TLS config
}

// TLSConfig is used to configure a TLS server
type TLSConfig struct {
	Addr string `json:"addr"` // TLS listen address
	Cert string `json:"cert"` // certificate file path
	Key  string `json:"key"`  // private key file path
}

// NewConfig returns a new Config given the config file path
func NewConfig(file string) *Config {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	cfg := &Config{}
	decoder := json.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
