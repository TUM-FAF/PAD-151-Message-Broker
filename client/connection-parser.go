package client

import (
	"PAD-151-Message-Broker/model"
)

// ConnectionParser ...
type ConnectionParser struct {
}

// ParseRequest ..
func (cp *ConnectionParser) ParseRequest(response string) (interface{}, error) {
	m := model.UserModel{}
	m.Name = response
	i, err := model.EncodeJsonMessage(connectionRequestHelper(m))
	return i, err
}

// ParseResponse ..
func (cp *ConnectionParser) ParseResponse(response string) (interface{}, error) {
	m := model.ConnectionModel{}
	err := model.DecodeYamlMessage([]byte(response), &m)
	return m, err
}

func connectionRequestHelper(in model.UserModel) interface{} {
	return in
}

func connectionResponseHelper(in model.ConnectionModel) interface{} {
	return in
}
