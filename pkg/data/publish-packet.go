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

	variableHeader := data[bytesRead+1:]

	topicName, n := util.GetUTFString(variableHeader)
	log.Printf("The packet is %d bytes long", packetSize)

	variableHeader = data[n+1:]

	//check qos to see if packet id will be present
	packetID := -1
	byteShift := 0
	if qos > 0 {
		byteShift = 2
		packetID = int(variableHeader[0])*256 + int(variableHeader[1])
	}

	payload := variableHeader[bytesRead+byteShift+1 : packetSize+1]
	log.Printf("Payload: %v, %s", payload, payload)
	for i := 0; i < 10; i++ {
		log.Printf("Byte: %d: %08b (-> %d)", i+1, payload[i], payload[i])
	}

	publishData, n := util.GetUTFString(payload)

	log.Printf("[%d] %s -> %s", packetID, topicName, publishData)

	return &PublishPacket{}, nil
}
