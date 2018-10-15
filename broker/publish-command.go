package broker

import (
	"PAD-151-Message-Broker/model"
)

// PublishCommand - send broadcast messages
type PublishCommand struct {
	message model.SentMessageModel
	broker  *Broker
}

// Execute - send
func (pc PublishCommand) Execute() {

}
