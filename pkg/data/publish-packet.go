package data

import (
	"log"
	"mqtt/pkg/util"
)

//PublishPacket represents a single subscribe packet
type PublishPacket struct {
	topic                Topic
	topicRaw             string
	data                 []byte
	dup                  byte
	retain               byte
	qos                  int
	packetID             int
	variableHeaderBeyond []byte
}

//LoadPublishPacket creates a ConnectPacket instance from the incoming packet data
func LoadPublishPacket(packetData []byte) (*PublishPacket, error) {
	log.Printf("Loading publish packet with %d bytes", len(packetData))
	header := packetData[0]
	//DUP flag
	dup := (header >> 3) & 1
	//Qos flag
	qos := int((header >> 1) & 3)
	//retain flag
	retain := header & 1

	log.Printf("dup %d, qos %d, retain %d", dup, qos, retain)

	_, bytesRead := util.RemainingLengthDecode(packetData[1:5])

	variableHeaderLen := 0
	variableHeader := packetData[bytesRead+1:]

	topicName, n := util.GetUTFString(variableHeader)
	variableHeaderLen += n
	topic := SplitFilter(topicName)

	variableHeader = variableHeader[n:]
	variableHeaderBeyond := variableHeader[:]

	//check qos to see if packet id will be present
	packetID := -1
	byteShift := 0
	if qos > 0 {
		log.Printf("Adding PacketID required by QoS %d", qos)
		byteShift = 2
		packetID = int(variableHeader[0])*256 + int(variableHeader[1])
	}
	variableHeaderLen += byteShift
	variableHeader = variableHeader[byteShift:]

	publishData := variableHeader

	log.Printf("[%d] %s -> %s", packetID, topicName, publishData)

	return &PublishPacket{
		topic:                topic,
		topicRaw:             topicName,
		data:                 publishData,
		dup:                  dup,
		retain:               retain,
		qos:                  qos,
		packetID:             packetID,
		variableHeaderBeyond: variableHeaderBeyond,
	}, nil
}

//ToByteArray returns the necessary bytes to send the packet over the wire
func (pp *PublishPacket) ToByteArray() []byte {
	var resp []byte

	//Insert all the necesssary fields

	//Packet Type: PUBLISH
	header := PUBLISH<<4 | pp.dup<<3 | byte(pp.qos)<<1 | pp.retain
	resp = append(resp, header)

	var variableContent []byte

	//Topic Name
	variableContent = append(variableContent, util.EncodeUTFString(pp.topicRaw)...)

	//Packet ID
	if pp.qos > 0 {
		bytePacketID := util.Get2ByteInt(pp.packetID)
		variableContent = append(variableContent, bytePacketID[0])
		variableContent = append(variableContent, bytePacketID[1])
		log.Print(variableContent)
	}

	variableContent = append(variableContent, pp.data...)

	//Remaining length is the 2 bytes for packet ID + payload length,
	//which is 1 byte per return code
	remainingLength := len(variableContent)
	remainingLengthBytes := util.RemainingLengthEncode(remainingLength)
	for _, remainingLengthByte := range remainingLengthBytes {
		resp = append(resp, remainingLengthByte)
	}

	resp = append(resp, variableContent...)

	return resp
}

func (pp *PublishPacket) GetOriginalPacket() []byte {
	return pp.ToByteArray()
}

func (pp *PublishPacket) GetTopic() Topic {
	return pp.topic
}
