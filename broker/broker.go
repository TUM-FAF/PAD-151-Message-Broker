package broker

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

type Broker struct {
	clientCount     int
	allClients      map[net.Conn]int
	newConnections  chan net.Conn
	deadConnections chan net.Conn
	messages        chan string
}

func (broker *Broker) Init() {
	// Number of people whom ever connected
	//
	broker.clientCount = 0

	// All people who are connected; a map wherein
	// the keys are net.Conn objects and the values
	// are client "ids", an integer.
	//
	broker.allClients = make(map[net.Conn]int)

	// Channel into which the TCP server will push
	// new connections.
	//
	broker.newConnections = make(chan net.Conn)

	// Channel into which we'll push dead connections
	// for removal from allClients.
	//
	broker.deadConnections = make(chan net.Conn)

	// Channel into which we'll push messages from
	// connected clients so that we can broadcast them
	// to every connection in allClients.
	//
	broker.messages = make(chan string)
}

func (broker *Broker) StartServer(connHost string, connPort string, connType string) {
	// Start the TCP server
	//
	server, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Server started at port:", connPort)
	}

	//Listen accept connection in another goroutine
	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			broker.newConnections <- conn
		}
	}()
}

func (broker *Broker) Run() {
	for {
		select {
		// Accept new clients
		//
		case conn := <-broker.newConnections:
			log.Printf("Accepted new client, #%d", broker.clientCount)

			// Add this connection to the `allClients` map
			//
			broker.allClients[conn] = broker.clientCount
			broker.clientCount++

			go getMessages(conn, broker.allClients[conn], broker.messages, broker.deadConnections)

		// Accept messages from connected clients
		//
		case message := <-broker.messages:

			// Loop over all connected clients
			//
			for conn := range broker.allClients {

				// Send them a message in a go-routine
				// so that the network operation doesn't block
				//
				go sendMessage(conn, message, broker.deadConnections)
			}
			log.Printf("New message: %s", message)
			log.Printf("Broadcast to %d clients", len(broker.allClients))

		// Remove dead clients
		//
		case conn := <-broker.deadConnections:
			log.Printf("Client %d disconnected", broker.allClients[conn])
			delete(broker.allClients, conn)
		}
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
