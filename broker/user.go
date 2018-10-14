package broker

import (
	"PAD-151-Message-Broker/model"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
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

func getUserName(conn net.Conn) string {
	// we create a decoder that reads directly from the socket
	var userModel model.UserModel
	reader := bufio.NewReader(conn)
	message, _ := reader.ReadString('\n')
	// fmt.Println(message)
	json.Unmarshal([]byte(message), &userModel)
	return userModel.Name
}

// Constantly read incoming messages from this
// user and push those onto
// the messages channel for broadcast to others.
//
func getMessages(user *User, messages chan<- string, deadUserIds chan<- int) {
	reader := bufio.NewReader(user.conn)
	for {
		incoming, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		messages <- fmt.Sprintf("Client %s > %s", user.name, incoming)
	}

	// When we encouter `err` reading, send this
	// connection to `deadUserIds` for removal.
	//
	deadUserIds <- user.id
}

// Send message to users connection
func sendMessage(user *User, message string, deadUserIds chan<- int) {
	_, err := user.conn.Write([]byte(message))

	// If there was an error communicating
	// with specified user, the connection is dead
	// and we add user's id to deadUserIds.
	if err != nil {
		deadUserIds <- user.id
	}
}
