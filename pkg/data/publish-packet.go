package data

import (
	"log"
	"mqtt/pkg/util"
)

//PublishPacket represents a single subscribe packet
type PublishPacket struct {
	topic          Topic
	data           []byte
	dup            byte
	retain         byte
	qos            int
	originalPacket []byte
}

//LoadPublishPacket creates a ConnectPacket instance from the incoming packet data
func LoadPublishPacket(data []byte) (*PublishPacket, error) {
	log.Printf("Loading connect packet with %d bytes", len(data))

	header := data[0]
	//DUP flag
	dup := (header >> 3) & 1
	//Qos flag
	qos := int((header >> 1) & 3)
	//retain flag
	retain := header & 1

	log.Printf("dup %d, qos %d, retain %d", dup, qos, retain)

	packetSize, bytesRead := util.RemainingLengthDecode(data[1:5])

	variableHeaderLen := 0
	variableHeader := data[bytesRead+1:]

	topicName, n := util.GetUTFString(variableHeader)
	variableHeaderLen += n
	topic := SplitFilter(topicName)
	log.Printf("The packet is %d bytes long", packetSize)

	variableHeader = data[n+1:]

	//check qos to see if packet id will be present
	packetID := -1
	byteShift := 0
	if qos > 0 {
		byteShift = 2
		variableHeaderLen += byteShift
		packetID = int(variableHeader[0])*256 + int(variableHeader[1])
	}

	payload := variableHeader[bytesRead+byteShift : packetSize+1]
	publishData := payload[:packetSize-variableHeaderLen]

	log.Printf("[%d] %s -> %s", packetID, topicName, publishData)

	return &PublishPacket{
		topic:          topic,
		data:           publishData,
		dup:            dup,
		retain:         retain,
		qos:            qos,
		originalPacket: data,
	}, nil
}

func (pp *PublishPacket) ToByteArray() []byte {
	var resp []byte

	//Insert all the necesssary fields

	//Packet Type: PUBLISH
	header := PUBLISH & pp.dup << 3 & byte(pp.qos) << 1 & pp.retain
	resp = append(resp, header)

	//Remaining length is the 2 bytes for packet ID + payload length,
	//which is 1 byte per return code
	/*
		remainingLength := 2 + len(cp.returnCodes)
		remainingLengthBytes := util.RemainingLengthEncode(remainingLength)
		for _, remainingLengthByte := range remainingLengthBytes {
			resp = append(resp, remainingLengthByte)
		}

		//Packet ID
		bytePacketID := int16(cp.packetID)
		resp = append(resp, byte(bytePacketID>>8))
		resp = append(resp, byte(bytePacketID&128))

		for _, b := range resp {
			log.Printf("%08b", b)
		}
	*/
	return resp
}

func (pp *PublishPacket) GetOriginalPacket() []byte {
	return pp.originalPacket
}

func (pp *PublishPacket) GetTopic() Topic {
	return pp.topic
}
