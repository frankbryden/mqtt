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
	server       *Server
	conn         net.Conn
	clientID     string
	retain       bool
	qos          int
	cleanSession bool
}

//NewClient constructs a client instance from the associated connection
//and received data contained in the connect packet
func NewClient(server *Server, conn net.Conn, clientData *data.ConnectPacket) *Client {
	return &Client{
		server:       server,
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
		n, e := c.conn.Read(buffer)
		if e != nil {
			log.Print(e)
			log.Printf("%s closed the connection", c.clientID)
			c.conn.Close()
			return
		}

		buffer = buffer[:n]

		log.Printf("Read %d bytes", n)
		if n == 0 {
			log.Printf("Skipping empty read")
			continue
		}
		log.Print(buffer)
		log.Print(string(buffer))

		header := buffer[0]
		//log.Printf("Header: %08b, %08b, %d", header, header>>4, int(header>>4))
		log.Printf(string(buffer))
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
			//Handle the subscription packet
			for _, sub := range subscribePacket.GetSubscriptions() {
				c.server.subManager.Subscribe(sub)
			}
			log.Println(subscribePacket)
			c.server.subManager.ListSubscriptions()
			subackPacket := data.NewSubackPacket(subscribePacket)
			n, err := c.conn.Write(subackPacket.ToByteArray())
			log.Printf("Sent SUBACK: Wrote %d bytes (%v)", n, err)
			break
		case data.PUBLISH:
			log.Print("Publish")
			log.Print(buffer)
			publishPacket, _ := data.LoadPublishPacket(buffer)
			log.Print(publishPacket)
			c.server.DispatchPublish(publishPacket)
			break
		}
	}
}

//Send data to a client
func (c *Client) Send(data []byte) {
	n, err := c.conn.Write(data)
	for i := 0; i < 5; i++ {
		//c.conn.Write([]byte("Hey " + string(i) + " FRANKIE"))
	}
	log.Printf("Wrote %d bytes (%s), sent to %s", n, err, c.clientID)
}
