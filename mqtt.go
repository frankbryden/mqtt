package main

import (
	"log"
	"mqtt/pkg/server"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	s := server.Server{}
	s.Start()
}
