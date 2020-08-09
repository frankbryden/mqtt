package data

//ConnackRemainingLength is the value to be put in the Remaining Length field
const ConnackRemainingLength = 2

//ConnackPacket is the response to a connect packet
type ConnackPacket struct {
	connectPacket *ConnectPacket
}

//NewConnackPacket returns a new connack packet instance based on a connect packet
func NewConnackPacket(connectPacket *ConnectPacket) *ConnackPacket {
	return &ConnackPacket{
		connectPacket: connectPacket,
	}
}

//ToByteArray returns the necessary bytes to send the packet over the wire
func (cp *ConnackPacket) ToByteArray() []byte {
	data := make([]byte, 4)

	//Fixed header
	data[0] = byte(2) << 4
	data[1] = byte(2)

	//Variable header
	data[2] = byte(1)
	data[3] = byte(0)

	return data
}
