package client

import (
	"PAD-151-Message-Broker/model"
	"fmt"

	"gopkg.in/yaml.v2"
)

// ConnectionParser ...
type ConnectionParser struct {
}

// ParseRequest ..
func (cp *ConnectionParser) ParseRequest(response string) (interface{}, error) {
	m := model.UserModel{}
	m.Name = response
	fmt.Println(m)
	// i := connectionRequestHelper(m)
	i, err := model.EncodeJsonMessage(connectionRequestHelper(m))
	return i, err
}

// ParseResponse ..
func (cp *ConnectionParser) ParseResponse(response string) (interface{}, error) {
	m := model.ConnectionModel{}
	err := yaml.Unmarshal([]byte(response), &m)
	fmt.Printf("%T\n", m)
	return m, err
}

func connectionRequestHelper(in model.UserModel) interface{} {
	return in
}

func connectionResponseHelper(in model.ConnectionModel) interface{} {
	return in
}
