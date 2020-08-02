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

//GetUTFString returns a single (and first encountered) UTF string from data
//along with the number of bytes read
func GetUTFString(data []byte) (string, int) {
	lengthData := data[:2]
	length := int(lengthData[0])*256 + int(lengthData[1])
	log.Printf("From %08b %08b, we get %d (%s)", lengthData[0], lengthData[1], length, data[2:length+2])
	return string(data[2 : length+2]), length + 2
}
