package broker

import (
	"PAD-151-Message-Broker/broker/command"
	"PAD-151-Message-Broker/model"
	"log"
)

// DispatchMessage - dispatch all incoming messages.
func DispatchMessage(message model.SentMessageModel, broker *Broker) command.Command {

	// Loop over all connected clients
	//
	responseMessageModel := model.ResponseMessageModel{
		message.SenderID,
		message.Type,
		-1,
		message.Message,
	}

	for id := range broker.userMap {

		// Send them a message in a go-routine
		// so that the network operation doesn't block
		//
		user := broker.userMap[id]

		// Send message to specified user

		go sendMessage(user, responseMessageModel, broker.deadUserIds)
	}
	log.Printf("New message: Client %s -> %s", broker.userMap[message.SenderID].name, message.Message)
	log.Printf("Broadcast to %d clients", len(broker.userMap))

	return command.IgnoreCommand{}
}
