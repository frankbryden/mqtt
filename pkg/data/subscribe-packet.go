package data

import (
	"log"
	"mqtt/pkg/util"
)

//SubscribePacket represents a single subscribe packet
type SubscribePacket struct {
	subscriptions []Subscription
	packetID      int
}

//LoadSubscribePacket creates a ConnectPacket instance from the incoming packet data
func LoadSubscribePacket(data []byte, clientID string) (*SubscribePacket, error) {
	packetSize, bytesRead := util.RemainingLengthDecode(data[1:5])

	variableHeader := data[bytesRead+1:]

	log.Print(variableHeader)

	packetID := 256*int(variableHeader[0]) + int(variableHeader[1])
	log.Printf("Subscribe packet has id %d", packetID)

	payload := variableHeader[2:packetSize]
	packetSize -= 2
	var topics []Subscription

	for packetSize > 0 {
		s, bytesRead := util.GetUTFString(payload)
		packetSize -= bytesRead
		payload = payload[bytesRead:]
		qos := int(payload[0])
		subscription := Subscription{
			filter:   SplitFilter(s),
			clientID: clientID,
			qos:      qos,
		}
		topics = append(topics, subscription)
		if packetSize == 0 {
			break
		}
		packetSize--
		payload = payload[1:]
	}

	return &SubscribePacket{
		subscriptions: topics,
		packetID:      packetID,
	}, nil
}

//GetSubscriptions returns subscriptions associated with subscribe packet
func (s *SubscribePacket) GetSubscriptions() []Subscription {
	return s.subscriptions
}
