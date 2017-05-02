package pubsub

import (
	"encoding/json"
	"os"
)

type Config struct {
}

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
