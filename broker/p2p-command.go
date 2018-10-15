package broker

import (
	"PAD-151-Message-Broker/model"
)

// P2pCommand - send broadcast messages
type P2pCommand struct {
	message model.SentMessageModel
	broker  *Broker
}

// Execute - send
func (pc P2pCommand) Execute() {

}
