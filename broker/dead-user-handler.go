package broker

import (
	"PAD-151-Message-Broker/model"
)

// DeadUserHandler ..
type DeadUserHandler struct {
	broker *Broker
}

// HandleExistingUser dead users ..
func (duh DeadUserHandler) HandleExistingUser(deadUser *User, deadMessage model.ResponseMessageModel) {
	duh.broker.deadUserIds <- deadUser.id
	DeadLetterCommand{deadUser.id, deadMessage, duh.broker}.Execute()
}

// HandleUnexistingUser ...
func (duh DeadUserHandler) HandleUnexistingUser(deadMessage model.ResponseMessageModel) {
	DeadLetterCommand{-1, deadMessage, duh.broker}.Execute()
}
