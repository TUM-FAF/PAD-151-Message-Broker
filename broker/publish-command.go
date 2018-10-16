package broker

import (
	"PAD-151-Message-Broker/model"
	"log"
)

// PublishCommand - send broadcast messages
type PublishCommand struct {
	message model.SentMessageModel
	broker  *Broker
}

// Execute - send
func (pc PublishCommand) Execute() {
	publisherID := pc.message.SenderID
	log.Printf("Clien %v posted", publisherID)
	responseMessageModel := model.ResponseMessageModel{
		pc.message.SenderID,
		pc.broker.userMap[pc.message.SenderID].name,
		pc.message.Type,
		-1,
		pc.message.Message,
	}

	for _, v := range pc.broker.subscribers[publisherID] {
		user := pc.broker.userMap[v]
		go sendMessage(user, responseMessageModel, DeadUserHandler{pc.broker})
	}
}
