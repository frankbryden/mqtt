package server

import (
	"log"
	"mqtt/pkg/data"
	"net"
)

//Server is the core of the mqtt broker
type Server struct {
	subManager *SubscriptionManager
	clients    map[string]*Client
}

//NewServer instantiates a server
func NewServer() *Server {
	subManager := NewSubscriptionManager()
	return &Server{
		subManager: subManager,
		clients:    make(map[string]*Client),
	}
}

//Start the server
func (s *Server) Start() {

	ln, err := net.Listen("tcp", "0.0.0.0:1883")
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

		c := NewClient(s, conn, connectPacket)
		s.clients[connectPacket.GetClientID()] = c
		c.Run()
	}

}

//DispatchPublish retransmists a publish packet to subscribed clients
func (s *Server) DispatchPublish(pp *data.PublishPacket) {
	for _, clientID := range s.subManager.GetMatchingClients(pp.GetTopic()) {
		s.clients[clientID].Send(pp.ToByteArray())
	}
}
