package server

import (
	"log"
	"mqtt/pkg/data"
	"net"
)

/*
willTopic, willMessage string
	username, password     string
	retain                 bool
	qos                    int
	cleanSession           bool
*/

//Client is an active connection to a client.
//Will eventually be used to keep track of state in persistent storage
type Client struct {
	conn         net.Conn
	clientID     string
	retain       bool
	qos          int
	cleanSession bool
}

//NewClient constructs a client instance from the associated connection
//and received data contained in the connect packet
func NewClient(conn net.Conn, clientData *data.ConnectPacket) *Client {
	return &Client{
		conn:         conn,
		clientID:     clientData.GetClientID(),
		retain:       clientData.ShouldRetain(),
		qos:          clientData.GetQos(),
		cleanSession: clientData.ShouldCleanSession(),
	}
}

//Run the client goroutine to serve incoming requests
func (c *Client) Run() {
	go c.serveRequests()
}

func (c *Client) serveRequests() {
	for {
		//Read(b []byte) (n int, err error)
		buffer := make([]byte, 256)
		n, _ := c.conn.Read(buffer)

		log.Printf("Read %d bytes", n)
		log.Print(buffer)
		log.Print(string(buffer))

		header := buffer[0]
		log.Printf("Header: %08b, %08b, %d", header, header>>4, int(header>>4))
		packetType := int(header >> 4)
		controlPacket, err := data.FromPacketType(packetType)
		log.Printf("Received control packet of type %v", controlPacket.PacketType)
		if err != nil {
			c.conn.Close()
		}
		switch controlPacket.PacketType {
		case data.SUBSCRIBE:
			log.Print("Subscribe")
			subscribePacket, _ := data.LoadSubscribePacket(buffer, c.clientID)
			log.Println(subscribePacket)
			subackPacket := data
			break
		case data.PUBLISH:
			log.Print("Publish")
			break
		}
	}
}
