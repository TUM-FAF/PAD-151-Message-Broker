package broker

import (
	"PAD-151-Message-Broker/model"
	"fmt"
	"log"
)

// DeadLetterCommand - send broadcast messages
type DeadLetterCommand struct {
	deadUserID int
	message    model.ResponseMessageModel
	broker     *Broker
}

// Execute - send
func (dlc DeadLetterCommand) Execute() {
	senderID := dlc.message.SenderID
	sender := dlc.broker.userMap[senderID]
	log.Printf("Send message to client %v", senderID)
	responseMessageModel := model.ResponseMessageModel{
		-1,
		8,
		-1,
		fmt.Sprintf("Dead Letter: Can't send message to %d", dlc.deadUserID),
	}

	// Send message to specified user

	go sendMessage(sender, responseMessageModel, DeadUserHandler{dlc.broker})
}
