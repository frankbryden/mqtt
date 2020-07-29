package data

import (
	"errors"
)

//PacketType of a single control packet
type PacketType int

const (
	//CONNECT client -> server
	CONNECT = iota + 1
	//CONNACK server -> client
	CONNACK = iota
	//PUBLISH client -> server
	PUBLISH = iota
	//PUBACK server -> client
	PUBACK = iota
	//PUBREC packet type
	PUBREC = iota
	//PUBREL packet type
	PUBREL = iota
	//PUBCOMP packet type
	PUBCOMP = iota
	//SUBSCRIBE client -> server
	SUBSCRIBE = iota
	//SUBACK server -> client
	SUBACK = iota
	//UNSUBSCRIBE client -> server
	UNSUBSCRIBE = iota
	//UNSUBACK server -> client
	UNSUBACK = iota
	//PINGREQ server -> client
	PINGREQ = iota
	//PINGRESP server -> client
	PINGRESP = iota
	//DISCONNECT server -> client
	DISCONNECT = iota
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
