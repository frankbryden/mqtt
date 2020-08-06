package server

import (
	"log"
	"mqtt/pkg/data"
)

//SubscriptionManager manages active subscriptions
type SubscriptionManager struct {
	//subscriptions is a list of subscriptions
	subscriptions []data.Subscription
}

//NewSubscriptionManager constructs an empty subscription manager
func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		subscriptions: []data.Subscription{},
	}
}

//Subscribe a client to a topic using a Subscription instance
func (s *SubscriptionManager) Subscribe(sub data.Subscription) {
	s.subscriptions = append(s.subscriptions, sub)
}

//ListSubscriptions prints out all active subscriptions
func (s *SubscriptionManager) ListSubscriptions() {
	for _, sub := range s.subscriptions {
		log.Print(sub)
		log.Printf("Client [%s] subscribed to %s", sub.GetClientID(), sub.GetFilter())
	}
}

//GetMatchingClients returns all client IDs which of clients subscribed
//to the provided Topic `t`
func (s *SubscriptionManager) GetMatchingClients(t data.Topic) []string {
	var clientIDs []string
	for _, sub := range s.subscriptions {
		if sub.MatchesTopic(t) {
			clientIDs = append(clientIDs, sub.GetClientID())
		}
	}
	return clientIDs
}
