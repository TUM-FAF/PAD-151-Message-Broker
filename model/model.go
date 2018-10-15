package model

import yaml "gopkg.in/yaml.v2"

// UserModel refer to message transfered by user at connection
type UserModel struct {
	Name string `json:"name"`
}

// ConnectionModel - on connection (or on request), broker send connected clients
type ConnectionModel struct {
	YourID int `yaml:"yourId"`
	Users  []struct {
		ID   int    `yaml:"id"`
		Name string `yaml:"name"`
	} `yaml:"users"`
	Rooms []struct {
		ID   int    `yaml:"id"`
		Room string `yaml:"room"`
	} `yaml:"rooms"`
}

// SentMessageModel - on new message, client send
type SentMessageModel struct {
	SenderID  int    `json:"senderID"`
	Type      int    `json:"type"`
	Receivers []int  `json:"receivers"`
	Message   string `json:"message"`
}

// ResponseMessageModel - on new message, broker send
type ResponseMessageModel struct {
	SenderID int    `json:"senderID"`
	Type     int    `json:"type"`
	Room     int    `json:"room"`
	Message  string `json:"message"`
}

// Parse d..
func (c *ConnectionModel) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}
