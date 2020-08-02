package data

import (
	"log"
	"strings"
)

//Topic represents a single topic, which is a collection of strings
type Topic []string

const (
	//SingleWildCard is a wildcard for a single level
	SingleWildCard = "+"
	//MultiWildCard is a wildcard matching multiple ending levels
	MultiWildCard = "#"
	//SysWildCard is for server-specific messages
	SysWildCard = "$"
	//Separator is found between topic levels
	Separator = "/"
)

//Subscription represents a single subscription. A topic can be matched against .
type Subscription struct {
	filter   Topic
	clientID string
	qos      int
}

//NewSubscriptionFromTopic constructs a new subscription from a Topic instance
func NewSubscriptionFromTopic(raw string) *Subscription {
	filter := SplitFilter(raw)
	if ValidateFilter(filter) {
		return &Subscription{
			filter: filter,
		}
	}
	return nil
}

//MatchesTopic returns true if the provided topic is contained within
//the underlying subscription
func (s *Subscription) MatchesTopic(t Topic) bool {
	log.Println(s.filter)
	if len(t) < len(s.filter) {
		//The only way this can match is if the next level is the final level
		//and it contains a multi wild card
		if !(len(s.filter) == len(t)+2 && s.filter[len(s.filter)-1] == MultiWildCard) {
			return false
		}
	}
	for i, token := range t {
		currentFilterLevel := s.filter[i]
		if i == len(s.filter)-1 {
			//The final level must either match the token, or be a wildcard
			if currentFilterLevel == MultiWildCard && !strings.Contains(token, SysWildCard) {
				return true
			}
			if i == len(t)-1 {
				//this is the last level of the topic
				return ((currentFilterLevel == SingleWildCard && !strings.Contains(token, SysWildCard)) ||
					currentFilterLevel == token)
			}
			return false
		}
		if token != currentFilterLevel && currentFilterLevel != MultiWildCard {
			return false
		}
	}
	return true
}

//ValidateFilter makes sure the filter follows all the rules defined by the spec
func ValidateFilter(filter Topic) bool {
	filterLen := len(filter)
	if filterLen < 1 {
		return false
	}
	for i, token := range filter {
		log.Println(token)
		if token == MultiWildCard {
			/*
				MultiWildCard must obey 3 rules:
					- All on its own (single item)
					- After a separator if not single
					- Must ALWAYS be final character
			*/
			if filterLen > 1 {
				//The filter is more than one char, need to check it is after a separator
				if i == 0 {
					//Length is more than 1 and WildCard is first char, cannot pass
					log.Println("1")
					return false
				}
				if filter[i-1] != Separator {
					//Previous char is not a Separator
					log.Println("2")
					return false
				}
				if i != filterLen-1 {
					//Wildcard is not at the end
					log.Println("3")
					return false
				}
			}
		} else if token == SingleWildCard {
			if len(token) > 1 {
				//Wildcard is not alone in the level
				log.Println("5")
				return false
			}
		} else if strings.Contains(token, SingleWildCard) || strings.Contains(token, MultiWildCard) {
			//Wildcard is not on its own in the current level
			return false
		}
	}
	return true
}

//SplitFilter takes a string and returns a topic (string split along Separator)
func SplitFilter(filter string) Topic {
	var final Topic
	filter = strings.Trim(filter, " ")
	split := strings.Split(filter, Separator)

	//Interleave split with Separator tokens
	for _, token := range split {
		final = append(final, token, Separator)
	}

	//Remove possibly extra Separator at the end
	//if string(filter[len(filter)-1]) != Separator {
	final = final[:len(final)-1]
	//}
	return final
}

func (t Topic) String() string {
	var b strings.Builder
	for _, token := range t {
		b.WriteString("'")
		b.WriteString(token)
		b.WriteString("'")
	}
	return b.String()
}

//GetClientID returns client id
func (s *Subscription) GetClientID() string {
	return s.clientID
}

//GetQos returns qos
func (s *Subscription) GetQos() int {
	return s.qos
}

//GetFilter returns Topic (filter)
func (s *Subscription) GetFilter() Topic {
	return s.filter
}
