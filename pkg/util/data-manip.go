package util

import "log"

//RemainingLengthDecode decodes up to 4 remaining length bytes,
//returning the decoded integer along with the number of bytes read.
func RemainingLengthDecode(data []byte) (int, int) {
	multiplier := 1
	value := 0
	bytesRead := 0

	for _, encodedByte := range data {
		bytesRead++

		value += int(encodedByte&127) * multiplier

		if multiplier > 128*128*128 {
			log.Fatal("Malformed Remaining Length")
		}

		multiplier *= 128

		if (encodedByte & 128) == 0 {
			break
		}
	}
	return value, bytesRead
}

//RemainingLengthEncode encodes the size following the remaining length encoding
//defined in the mqtt standard
func RemainingLengthEncode(x int) []byte {
	var output []byte
	for {
		encodedByte := x % 128

		x = x / 128

		// if there are more data to encode, set the top bit of this byte

		if x > 0 {
			encodedByte = encodedByte | 128
		}
		output = append(output, byte(encodedByte))
		if x <= 0 {
			break
		}
	}
	return output
}

//GetUTFString returns a single (and first encountered) UTF string from data
//along with the number of bytes read
func GetUTFString(data []byte) (string, int) {
	lengthData := data[:2]
	length := int(lengthData[0])*256 + int(lengthData[1])
	log.Printf("From %08b %08b, we get %d (%s)", lengthData[0], lengthData[1], length, data[2:length+2])
	return string(data[2 : length+2]), length + 2
}

//EncodeUTFString takes a string and returns a slice of bytes with
//the string prepended by its length (2 bytes)
func EncodeUTFString(data string) []byte {
	res := make([]byte, len(data)+2)
	byteLen := Get2ByteInt(len(data))
	res[0] = byteLen[0]
	res[1] = byteLen[1]
	for i, v := range []byte(data) {
		res[i+2] = v
	}
	return res
}

//Get2ByteInt returns 2 bytes making up the int
func Get2ByteInt(data int) []byte {
	res := make([]byte, 2)
	byteInt := int16(data)
	res[0] = byte(byteInt >> 8)
	res[1] = byte(byteInt & 255)
	return res
}
