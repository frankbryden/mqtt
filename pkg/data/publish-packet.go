package data

import (
	"log"
	"mqtt/pkg/util"
)

//PublishPacket represents a single subscribe packet
type PublishPacket struct {
}

//LoadPublishPacket creates a ConnectPacket instance from the incoming packet data
func LoadPublishPacket(data []byte) (*PublishPacket, error) {
	log.Printf("Loading connect packet with %d bytes", len(data))

	header := data[0]
	//DUP flag
	dup := (header >> 3) & 1
	//Qos flag
	qos := (header >> 1) & 3
	//retain flag
	retain := header & 1

	log.Printf("dup %d, qos %d, retain %d", dup, qos, retain)

	packetSize, bytesRead := util.RemainingLengthDecode(data[1:5])
	log.Printf("The packet is %d bytes long", packetSize)

	payload := data[bytesRead : packetSize+1]
	log.Printf("Payload: %v, %s", payload, payload)
	for i := 0; i < 10; i++ {
		log.Printf("Byte: %d: %08b (-> %d)", i+1, payload[i], payload[i])
	}
	versionNumber := payload[6]
	log.Printf("Client using mqtt version %d", versionNumber)

	connectPacket := ConnectPacket{}

	flagsInstance := loadFlags(payload)
	connectPacket.flags = flagsInstance

	payload = payload[10:]

	//Client Identifier, Will Topic, Will Message, User Name, Password

	//Client Identifier
	clientID, n := util.GetUTFString(payload)
	connectPacket.clientID = clientID
	payload = payload[n:]

	//Will Topic/Message
	if flagsInstance.willFlag {
		willTopic, n := util.GetUTFString(payload)
		connectPacket.willTopic = willTopic
		payload = payload[n:]

		willMessage, n := util.GetUTFString(payload)
		connectPacket.willMessage = willMessage
		payload = payload[n:]
	}

	//User name
	if flagsInstance.usernameFlag {
		username, n := util.GetUTFString(payload)
		connectPacket.username = username
		payload = payload[n:]
	}

	//Password
	if flagsInstance.usernameFlag {
		password, n := util.GetUTFString(payload)
		connectPacket.password = password
		payload = payload[n:]
	}

	return &PublishPacket{}, nil
}
