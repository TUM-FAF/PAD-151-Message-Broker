package broker

import (
	"PAD-151-Message-Broker/model"
	"log"
	"strconv"
	"strings"
)

// CusersCommand - send broadcast messages
type CusersCommand struct {
	message model.SentMessageModel
	broker  *Broker
}

// Execute - send
func (cc CusersCommand) Execute() {
	log.Printf("Get list of active users for client %v", cc.message.SenderID)
	responseMessageModel := model.ResponseMessageModel{
		cc.message.SenderID,
		cc.message.Type,
		-1,
		cc.message.Message,
	}

	users := make([]string, 0)
	user := cc.broker.userMap[cc.message.SenderID]
	for _, v := range cc.broker.userMap {
		idName := strconv.Itoa(v.id)
		users = append(users, v.name+" "+idName)
	}
	responseMessageModel.Message = strings.Join(users, ", ")

	go sendMessage(user, responseMessageModel, DeadUserHandler{cc.broker})

}
