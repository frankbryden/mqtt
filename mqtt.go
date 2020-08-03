package main

import (
	"fmt"
	"log"
	"mqtt/pkg/server"
)

const (
	C1 = iota + 1
	C2
	C3
)

func main() {
	fmt.Println(C1, C2, C3) // "1 2 3"
	log.SetFlags(log.Lshortfile | log.Ltime)
	s := server.NewServer()
	s.Start()
}
