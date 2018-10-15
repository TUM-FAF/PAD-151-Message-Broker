package broker

import (
	"PAD-151-Message-Broker/model"
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// User is type representing the client that connects to server
type User struct {
	conn net.Conn
	id   int
	name string
}

// Init initiates data in user
func (user *User) Init(conn net.Conn, id int) {
	user.conn = conn
	user.id = id
	user.name = getUserName(conn)
	log.Printf("Accepted new client, #%s", user.name)
}

func (user *User) sendConnectResponse(response interface{}) {
	model.EncodeYamlMessage(response.([]byte))
}

func getUserName(conn net.Conn) string {
	// we create a decoder that reads directly from the socket
	var userModel model.UserModel
	reader := bufio.NewReader(conn)
	message, _ := reader.ReadString('\n')
	// fmt.Println(message)
	message = strings.TrimSuffix(message, "\n")
	model.DecodeJsonMessage([]byte(message), &userModel)
	return userModel.Name
}

// Constantly read incoming messages from this
// user and push those onto
// the messages channel for broadcast to others.
//
func getMessages(user *User, messages chan<- model.SentMessageModel, deadUserIds chan<- int) {
	reader := bufio.NewReader(user.conn)
	for {
		incoming, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		// messages <- fmt.Sprintf("Client %s > %s", user.name, incoming)
		incoming = strings.TrimSuffix(incoming, "\n")
		sentMessageModel := model.SentMessageModel{}

		model.DecodeJsonMessage([]byte(incoming), &sentMessageModel)
		sentMessageModel.SenderID = user.id
		messages <- sentMessageModel
	}

	// When we encouter `err` reading, send this
	// connection to `deadUserIds` for removal.
	//
	deadUserIds <- user.id
}

// Send message to users connection
func sendMessage(user *User, response model.ResponseMessageModel, deadUserIds chan<- int) {

	data, error := model.EncodeYamlMessage(response)
	if error != nil {
		fmt.Println(error)
	}
	_, err := user.conn.Write(data)

	// If there was an error communicating
	// with specified user, the connection is dead
	// and we add user's id to deadUserIds.
	if err != nil {
		deadUserIds <- user.id
	}
}
