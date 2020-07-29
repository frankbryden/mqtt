package tests

import (
	"mqtt/pkg/util"
	"testing"
)

func TestDecodeRemainingLength(t *testing.T) {
	input := make([]byte, 4)
	input[0] = 193
	input[1] = 2

	expected := 321

	got := util.RemainingLengthDecode(input)

	if got != expected {
		t.Errorf("Got %d, expected %d", got, expected)
	}
}

func TestDecodeRemainingLengthMaxValue(t *testing.T) {
	input := make([]byte, 4)
	input[0] = 0xFF
	input[1] = 0xFF
	input[2] = 0xFF
	input[3] = 0x7F

	expected := 268435455

	got := util.RemainingLengthDecode(input)

	if got != expected {
		t.Errorf("Got %d, expected %d", got, expected)
	}
}

func TestEncodeRemainingLength(t *testing.T) {

}
