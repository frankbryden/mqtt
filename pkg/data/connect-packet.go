package data

import (
	"log"
	"mqtt/pkg/util"
)

//ConnectPacketFlags represents flags set in CONNECT packet header
type ConnectPacketFlags struct {
	cleanSession bool
	willFlag     bool
	willQos      int
	willRetain   bool
	passwordFlag bool
	usernameFlag bool
}

//ConnectPacket contains data provided by client when connecting
type ConnectPacket struct {
	flags                  *ConnectPacketFlags
	clientID               string
	willTopic, willMessage string
	username, password     string
	retain                 bool
	qos                    int
	cleanSession           bool
}

//LoadConnectPacket creates a ConnectPacket instance from the incoming packet data
func LoadConnectPacket(data []byte) (*ConnectPacket, error) {
	log.Printf("Loading connect packet with %d bytes", len(data))

	packetSize, bytesRead := util.RemainingLengthDecode(data[:4])
	log.Printf("The packet is %d bytes long", packetSize)

	payload := data[bytesRead : packetSize+1]
	log.Printf("Payload: %v, %s", payload, payload)
	for i := 0; i < 10; i++ {
		log.Printf("Byte: %d: %08b (-> %d)", i+1, payload[i], payload[i])
	}

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

	return &connectPacket, nil
}

func loadFlags(payload []byte) *ConnectPacketFlags {
	flags := payload[7]
	flagsInstance := &ConnectPacketFlags{}
	log.Printf("Flags: %08b", flags)

	mask := byte(1)

	//Reserved
	if flags&mask == 0 {
		log.Println("Reserved bit correctly set to 0")
	}
	mask = mask << 1

	//Clean Session
	if flags&mask == mask {
		flagsInstance.cleanSession = true
		log.Println("Clean Session")
	}
	mask = mask << 1

	//Will Flag
	if flags&mask == mask {
		flagsInstance.willFlag = true
		log.Println("Has will message")
	}
	mask = mask << 1

	//Will QoS
	if flagsInstance.willFlag {
		flagsInstance.willQos = int(flags & (mask & mask << 1))
	}

	mask = mask << 2

	//Will Retain
	if flags&mask == mask {
		flagsInstance.willRetain = true
		log.Println("Retain set")
	}
	mask = mask << 1

	//Password Flag
	if flags&mask == mask {
		log.Println("Has password")
		flagsInstance.passwordFlag = true
	}
	mask = mask << 1

	//User Name Flag
	if flags&mask == mask {
		flagsInstance.usernameFlag = true
		log.Println("Has user name")
	}

	return flagsInstance
}
