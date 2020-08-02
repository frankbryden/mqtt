package server

import (
	"log"
	"mqtt/pkg/data"
	"net"
)

type Server struct {
}

type Limbo struct {
	conn net.Conn
}

//Start the server
func (s *Server) Start() {

	ln, err := net.Listen("tcp", "127.0.0.1:1883")
	log.Printf("Listening on port %d", 1883)
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go s.HandleConn(conn)
	}
}

//HandleConn sets up an incoming connection from a potential client
func (s *Server) HandleConn(conn net.Conn) {
	//Read(b []byte) (n int, err error)
	buffer := make([]byte, 256)
	n, _ := conn.Read(buffer)

	log.Printf("Read %d bytes", n)

	header := buffer[0]
	packetType := int(header >> 4)
	controlPacket, err := data.FromPacketType(packetType)
	if err != nil {
		conn.Close()
	}
	log.Println(controlPacket)

	if controlPacket.PacketType == data.CONNECT {
		connectPacket, err := data.LoadConnectPacket(buffer[1:])
		if err != nil {
			conn.Close()
		}
		log.Println(connectPacket)
		connackPacket := data.NewConnackPacket(connectPacket)
		conn.Write(connackPacket.ToByteArray())

		c := NewClient(conn, connectPacket)
		c.Run()
	}

}
