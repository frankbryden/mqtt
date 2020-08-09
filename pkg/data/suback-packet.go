package data

import (
	"log"
	"mqtt/pkg/util"
)

//SubackPacket is the response to a SUBSCRIBE packet
type SubackPacket struct {
	packetID    int
	returnCodes []int
}

//NewSubackPacket instantiates a SUBACK packet from a SUBSCRIBE packet
func NewSubackPacket(s *SubscribePacket) *SubackPacket {
	returnCodes := make([]int, len(s.subscriptions))
	for i, sub := range s.subscriptions {
		if sub.GetQos() < 1 {
			returnCodes[i] = 0
		} else {
			returnCodes[i] = 1
		}
	}
	return &SubackPacket{
		packetID:    s.packetID,
		returnCodes: returnCodes,
	}
}

//ToByteArray returns the necessary bytes to send the packet over the wire
func (cp *SubackPacket) ToByteArray() []byte {
	var resp []byte

	//Insert all the necesssary fields

	//Packet Type: SUBACK
	resp = append(resp, SUBACK<<4)

	//Remaining length is the 2 bytes for packet ID + payload length,
	//which is 1 byte per return code
	remainingLength := 2 + len(cp.returnCodes)
	remainingLengthBytes := util.RemainingLengthEncode(remainingLength)
	for _, remainingLengthByte := range remainingLengthBytes {
		resp = append(resp, remainingLengthByte)
	}

	//Packet ID
	bytePacketID := int16(cp.packetID)
	resp = append(resp, byte(bytePacketID>>8))
	resp = append(resp, byte(bytePacketID&255))
	log.Printf("Packet ID %d becomes %d", cp.packetID, bytePacketID)

	for _, returnCode := range cp.returnCodes {
		resp = append(resp, byte(returnCode))
	}

	for _, b := range resp {
		log.Printf("%08b", b)
	}
	return resp
}
