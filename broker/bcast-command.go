package broker

import (
	"PAD-151-Message-Broker/model"
	"log"
)

// BcastCommand - send broadcast messages
type BcastCommand struct {
	message model.SentMessageModel
	broker  *Broker
}

// Execute - send
func (bc BcastCommand) Execute() {
	log.Printf("Broadcast to %d clients", len(bc.broker.userMap))
	responseMessageModel := model.ResponseMessageModel{
		bc.message.SenderID,
		bc.message.Type,
		-1,
		bc.message.Message,
	}
	for id := range bc.broker.userMap {

		// Send them a message in a go-routine
		// so that the network operation doesn't block
		//
		user := bc.broker.userMap[id]

		// Send message to specified user

		go sendMessage(user, responseMessageModel, bc.broker.deadUserIds)
	}
}
