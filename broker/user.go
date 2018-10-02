package broker

import (
	"bufio"
	"fmt"
	"net"
)

// User is type representing the client that connects to server
type User struct {
	conn net.Conn
	id   int
}

// NewUser returns a new user and return pointer ot it
func NewUser(conn net.Conn, id int) *User {
	return &User{
		conn: conn,
		id:   id,
	}
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
		messages <- fmt.Sprintf("Client %d > %s", user.id, incoming)
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
