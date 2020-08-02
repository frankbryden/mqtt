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

//Subscribe a client to a topic using a Subscription instance
func (s *SubscriptionManager) Subscribe(sub data.Subscription) {
	s.subscriptions = append(s.subscriptions, sub)
}

//ListSubscriptions prints out all active subscriptions
func (s *SubscriptionManager) ListSubscriptions() {
	for _, sub := range s.subscriptions {
		log.Printf("Client [%s] subscribed to %s", sub.GetClientID(), sub.GetFilter())
	}
}
