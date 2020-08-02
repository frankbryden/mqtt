package data

import (
	"errors"
)

//PacketType of a single control packet
type PacketType int

const (
	//CONNECT client -> server
	CONNECT = 1
	//CONNACK server -> client
	CONNACK = 2
	//PUBLISH client -> server
	PUBLISH = 3
	//PUBACK server -> client
	PUBACK = 4
	//PUBREC packet type
	PUBREC = 5
	//PUBREL packet type
	PUBREL = 6
	//PUBCOMP packet type
	PUBCOMP = 7
	//SUBSCRIBE client -> server
	SUBSCRIBE = 8
	//SUBACK server -> client
	SUBACK = 9
	//UNSUBSCRIBE client -> server
	UNSUBSCRIBE = 10
	//UNSUBACK server -> client
	UNSUBACK = 11
	//PINGREQ server -> client
	PINGREQ = 12
	//PINGRESP server -> client
	PINGRESP = 13
	//DISCONNECT server -> client
	DISCONNECT = 14
)

//ControlPacket contains information about a single control packet
type ControlPacket struct {
	PacketType PacketType
	flags      int
}

//FromPacketType constructs a ControlPacket instance from the packet type header
func FromPacketType(val int) (*ControlPacket, error) {
	var t PacketType
	switch val {
	case 1:
		t = CONNECT
		break
	case 2:
		t = CONNACK
		break
	case 3:
		t = PUBLISH
		break
	case 4:
		t = PUBACK
		break
	case 5:
		t = PUBREC
		break
	case 6:
		t = PUBREL
		break
	case 7:
		t = PUBCOMP
		break
	case 8:
		t = SUBSCRIBE
		break
	case 9:
		t = SUBACK
		break
	case 10:
		t = UNSUBSCRIBE
		break
	case 11:
		t = UNSUBACK
		break
	case 12:
		t = PINGREQ
		break
	case 13:
		t = PINGRESP
		break
	case 14:
		t = DISCONNECT
		break
	default:
		return nil, errors.New("invalid control packet type")
	}
	return &ControlPacket{
		PacketType: t,
		flags:      0,
	}, nil
}

func (p PacketType) String() string {
	return [...]string{"Reserved", "CONNECT", "CONNACK", "PUBLISH", "PUBACK", "PUBREC", "PUBREL", "PUBCOMP",
		"SUBSCRIBE", "SUBACK", "UNSUBSCRIBE", "UNSUBACK", "PINGREQ", "PINGRESP", "DISCONNECT"}[p]
}

func (p ControlPacket) String() string {
	return "CP: " + p.PacketType.String()
}
