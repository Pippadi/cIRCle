package persistence

import (
	"log"

	"github.com/Pippadi/cIRCle/src/connection"
)

func LoadConnConfig() *connection.Config {
	return &connection.Config{Host: "localhost", Port: 6667, Nick: "mario", Channels: []string{"#op"}}
}

func DumpConnConfig(conf *connection.Config) { log.Println(conf) }
