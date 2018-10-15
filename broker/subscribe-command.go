package broker

import (
	"PAD-151-Message-Broker/model"
)

// SubscribeCommand - send broadcast messages
type SubscribeCommand struct {
	message model.SentMessageModel
	broker  *Broker
}

// Execute - send
func (pc SubscribeCommand) Execute() {

}
