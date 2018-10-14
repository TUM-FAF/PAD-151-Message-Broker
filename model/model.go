package model

// UserModel refer to message transfered by user at connection
type UserModel struct {
	Name string `json:"name"`
}

// ConnectionModel - on connection (or on request), broker send connected clients
type ConnectionModel struct {
	YourID int `json:"yourId"`
	Users  []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"users"`
	Rooms []struct {
		ID   int    `json:"id"`
		Room string `json:"room"`
	} `json:"rooms"`
}

// SendMessageModel - on new message, client send
type SendMessageModel struct {
	Type      int    `json:"type"`
	Receivers []int  `json:"receivers"`
	Message   string `json:"message"`
}

// ResponseMessageModel - on new message, broker send
type ResponseMessageModel struct {
	SenderUser int    `json:"senderUser"`
	Type       int    `json:"type"`
	Room       int    `json:"room"`
	Message    string `json:"message"`
}
