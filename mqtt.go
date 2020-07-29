package main

import "mqtt/pkg/server"

func main() {
	s := server.Server{}
	s.Start()
}
