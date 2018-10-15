package broker

import (
	"PAD-151-Message-Broker/model"
	"log"
)

// P2pCommand - send broadcast messages
type P2pCommand struct {
	message model.SentMessageModel
	broker  *Broker
}

// Execute - send
func (pc P2pCommand) Execute() {
	id := pc.message.Receivers[0]
	log.Printf("Send message to client %v", pc.message.Receivers)
	responseMessageModel := model.ResponseMessageModel{
		pc.message.SenderID,
		pc.broker.userMap[pc.message.SenderID].name,
		pc.message.Type,
		-1,
		pc.message.Message,
	}

	user := pc.broker.userMap[id]

	// Send message to specified user

	go sendMessage(user, responseMessageModel, pc.broker.deadUserIds)

}
