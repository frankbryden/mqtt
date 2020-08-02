package data

import (
	"log"
)

type SubackPacket struct {
	packetID    int
	returnCodes []int
}

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

func (cp *SubackPacket) ToByteArray() []byte {
	resp := make([]byte, 4)
	resp[0] = SUBACK

	bytePacketID := byte(cp.packetID)
	resp[2] = bytePacketID >> 8
	resp[3] = bytePacketID & 128

	for _, b := range resp {
		log.Printf("%08b", b)
	}
	return resp
}
