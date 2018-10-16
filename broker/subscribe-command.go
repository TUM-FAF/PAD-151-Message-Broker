package broker

import (
	"PAD-151-Message-Broker/model"
	"log"
)

// SubscribeCommand - send broadcast messages
type SubscribeCommand struct {
	message model.SentMessageModel
	broker  *Broker
}

// Execute - send
func (pc SubscribeCommand) Execute() {
	subscriberID := pc.message.SenderID
	publisherID := pc.message.Receivers[0]

	log.Printf("Subscribe to client %v", publisherID)
	pc.broker.subscribers[publisherID] = append(pc.broker.subscribers[publisherID], subscriberID)
}
