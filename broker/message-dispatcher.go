package broker

import (
	"PAD-151-Message-Broker/model"
)

// DispatchMessage - dispatch all incoming messages.
func DispatchMessage(message model.SentMessageModel, broker *Broker) Command {
	switch message.Type {
	case 0:
		return BcastCommand{message, broker}
	case 1:
		return P2pCommand{message, broker}
	case 3:
		return PublishCommand{message, broker}
	case 6:
		return SubscribeCommand{message, broker}
	}
	return IgnoreCommand{}
}
