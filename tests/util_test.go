package tests

import (
	"log"
	"mqtt/pkg/data"
	"mqtt/pkg/util"
	"reflect"
	"testing"
)

func TestDecodeRemainingLength(t *testing.T) {
	input := make([]byte, 4)
	input[0] = 193
	input[1] = 2

	expected := 321

	got, _ := util.RemainingLengthDecode(input)

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

	got, _ := util.RemainingLengthDecode(input)

	if got != expected {
		t.Errorf("Got %d, expected %d", got, expected)
	}
}

func TestEncodeRemainingLength(t *testing.T) {

}

func TestSplitFilterStartEndSep(t *testing.T) {
	input := "/bathroom/temp/"
	expected := data.Topic{"", "/", "bathroom", "/", "temp", "/", ""}

	got := data.SplitFilter(input)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Got %v, expected %v", got, expected)
	}
}

func TestSplitFilterNoSurroundingSep(t *testing.T) {
	input := "bathroom/temp"
	expected := data.Topic{"bathroom", "/", "temp"}

	got := data.SplitFilter(input)

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Got %v, expected %v", got, expected)
	}
}

func TestValidateFilterEndingMutilWildcard(t *testing.T) {
	input := data.SplitFilter("sport/tennis/#")
	expected := true

	log.Println(input)

	got := data.ValidateFilter(input)

	if got != expected {
		t.Errorf("Got %t, expected %t", got, expected)
	}
}

func TestValidateFilterBurriedWildcard(t *testing.T) {
	input := data.SplitFilter("sport/tennis#")
	expected := false

	got := data.ValidateFilter(input)

	if got != expected {
		t.Errorf("Got %t, expected %t", got, expected)
	}
}

func TestValidateFilterNonEndingMutilWildcard(t *testing.T) {
	input := data.SplitFilter("sport/tennis/#/ranking")
	expected := false

	got := data.ValidateFilter(input)

	if got != expected {
		t.Errorf("Got %t, expected %t", got, expected)
	}
}

func TestMatchesTopicMultiWildCard(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.Ltime)
	filter := data.NewSubscriptionFromTopic("sport/tennis/player1/#")
	rawInputs := []string{"sport/tennis/player1", "sport/tennis/player1/ranking", "sport/tennis/player1/score/wimbledon"}

	for _, rawInput := range rawInputs {
		input := data.SplitFilter(rawInput)
		if !filter.MatchesTopic(input) {
			t.Errorf("%s was supposed to match the filter, but it didn't", rawInput)
		}
	}
}

func TestMatchesTopicSingleWildCard(t *testing.T) {
	filter := data.NewSubscriptionFromTopic("sport/tennis/+")
	rawInputs := []string{"sport/tennis/player1", "sport/tennis/player2"}

	for _, rawInput := range rawInputs {
		input := data.SplitFilter(rawInput)
		if !filter.MatchesTopic(input) {
			t.Errorf("%s was supposed to match the filter, but it didn't", rawInput)
		}
	}
}

func TestNonMatchesTopicSingleWildCard(t *testing.T) {
	filter := data.NewSubscriptionFromTopic("sport/tennis/+")
	rawInputs := []string{"sport/tennis/player1/ranking", "sport/tennis/player2/hey", "sport"}

	for _, rawInput := range rawInputs {
		input := data.SplitFilter(rawInput)
		if filter.MatchesTopic(input) {
			t.Errorf("%s was not supposed to match the filter, but it did", rawInput)
		}
	}
}

func TestMatchesAll(t *testing.T) {
	filter := data.NewSubscriptionFromTopic("#")
	rawInputs := []string{"sport/tennis/player1/ranking", "sport/tennis/player2/hey", "sport"}

	for _, rawInput := range rawInputs {
		input := data.SplitFilter(rawInput)
		if !filter.MatchesTopic(input) {
			t.Errorf("%s was supposed to match the filter, but it didn't", rawInput)
		}
	}
}

func TestSysWildcard(t *testing.T) {
	rawInput := "$SYS/hi"
	input := data.SplitFilter(rawInput)
	matchingRawFilters := []string{"$SYS/+", "$SYS/#", "$SYS/hi"}
	nonMatchingRawFilters := []string{"#", "SYS", "+/hi"}

	for _, matchingRawFilter := range matchingRawFilters {
		matchingFilter := data.NewSubscriptionFromTopic(matchingRawFilter)
		if !matchingFilter.MatchesTopic(input) {
			t.Errorf("%s was supposed to match the filter %s, but it didn't", rawInput, matchingRawFilter)
		}
	}

	for _, nonMatchingRawFilter := range nonMatchingRawFilters {
		nonMatchingFilter := data.NewSubscriptionFromTopic(nonMatchingRawFilter)
		if nonMatchingFilter.MatchesTopic(input) {
			t.Errorf("%s was not supposed to match the filter %s, but it did", rawInput, nonMatchingRawFilter)
		}
	}
}
