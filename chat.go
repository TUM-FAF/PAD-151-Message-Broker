/*
	chat.go -- A minimal Go TCP chat example.
	Run this like
		> go run chat.go
	That will run a TCP chat server at localhost:8080.
	You can connect to that chat server like
		> telnet localhost 8080
	And, of course, others can connect using your IP
	address like
		> telnet YOUR-IP-HERE 8080
	assuming your firewall allows it.
*/

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	// ConnHost means connection host
	ConnHost = "0.0.0.0"
	// ConnPort means connection port
	ConnPort = "8080"
	// ConnType means connection type
	ConnType = "tcp"
)

func main() {

	// Number of people whom ever connected
	//
	clientCount := 0

	// All people who are connected; a map wherein
	// the keys are net.Conn objects and the values
	// are client "ids", an integer.
	//
	allClients := make(map[net.Conn]int)

	// Channel into which the TCP server will push
	// new connections.
	//
	newConnections := make(chan net.Conn)

	// Channel into which we'll push dead connections
	// for removal from allClients.
	//
	deadConnections := make(chan net.Conn)

	// Channel into which we'll push messages from
	// connected clients so that we can broadcast them
	// to every connection in allClients.
	//
	messages := make(chan string)

	// Start the TCP server
	//
	server, err := net.Listen(ConnType, ConnHost+":"+ConnPort)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Server started at port:", ConnPort)
	}

	// Accept connections in a separate go-routine
	// and add them to new connection
	go acceptConnection(server, newConnections)

	// Loop endlessly
	//
	for {

		// Handle 1) new connections; 2) dead connections;
		// and, 3) broadcast messages.
		//
		select {

		// Accept new clients
		//
		case conn := <-newConnections:

			log.Printf("Accepted new client, #%d", clientCount)

			// Add this connection to the `allClients` map
			//
			allClients[conn] = clientCount
			clientCount++

			go getMessages(conn, allClients[conn], messages, deadConnections)

		// Accept messages from connected clients
		//
		case message := <-messages:

			// Loop over all connected clients
			//
			for conn := range allClients {

				// Send them a message in a go-routine
				// so that the network operation doesn't block
				//
				go sendMessage(conn, message, deadConnections)
			}
			log.Printf("New message: %s", message)
			log.Printf("Broadcast to %d clients", len(allClients))

		// Remove dead clients
		//
		case conn := <-deadConnections:
			log.Printf("Client %d disconnected", allClients[conn])
			delete(allClients, conn)
		}
	}
}

// Tell the server to accept connections forever
// and push new connections into the newConnections channel.
//
func acceptConnection(server net.Listener, newConnections chan<- net.Conn) {
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		newConnections <- conn
	}
}

// Constantly read incoming messages from this
// client in a goroutine and push those onto
// the messages channel for broadcast to others.
//
func getMessages(conn net.Conn, clientID int, messages chan<- string, deadConnections chan<- net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		incoming, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		messages <- fmt.Sprintf("Client %d > %s", clientID, incoming)
	}

	// When we encouter `err` reading, send this
	// connection to `deadConnections` for removal.
	//
	deadConnections <- conn
}

// Send message to connection
func sendMessage(conn net.Conn, message string, deadConnections chan<- net.Conn) {
	_, err := conn.Write([]byte(message))

	// If there was an error communicating
	// with them, the connection is dead.
	if err != nil {
		deadConnections <- conn
	}
}
