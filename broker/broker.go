package broker

import (
	"fmt"
	"log"
	"net"
	"os"
)

// Broker is type representing the message broker
type Broker struct {
	clientCount    int
	userMap        map[int]*User
	newConnections chan net.Conn
	deadUserIds    chan int
	messages       chan string
}

// Init initiates broker data
func (broker *Broker) Init() {
	// Number of people whom ever connected
	//
	broker.clientCount = 0

	broker.userMap = make(map[int]*User)

	// Channel into which the TCP server will push
	// new connections.
	//
	broker.newConnections = make(chan net.Conn)

	// Channel into which we'll push dead connections
	// for removal from allClients.
	//
	broker.deadUserIds = make(chan int)

	// Channel into which we'll push messages from
	// connected clients so that we can broadcast them
	// to every connection in allClients.
	//
	broker.messages = make(chan string)
}

// StartServer creates server, accepts connetions and runs broker
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
	defer server.Close()

	//Listen accepted connection in another goroutine
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

	broker.Run()
}

// Run handle 1) new connections; 2) dead connections;
// and, 3) broadcast messages.
//
func (broker *Broker) Run() {
	for {
		select {
		// Accept new clients
		//
		case conn := <-broker.newConnections:
			// Create user and add him to the `userMap`
			//
			user := new(User)
			user.Init(conn, broker.clientCount)

			broker.userMap[broker.clientCount] = user

			broker.clientCount++

			go getMessages(user, broker.messages, broker.deadUserIds)

		// Accept messages from connected client
		//
		case message := <-broker.messages:

			// Loop over all connected clients
			//
			for id := range broker.userMap {

				// Send them a message in a go-routine
				// so that the network operation doesn't block
				//
				user := broker.userMap[id]

				// Send message to specified user
				go sendMessage(user, message, broker.deadUserIds)
			}
			log.Printf("New message: %s", message)
			log.Printf("Broadcast to %d clients", len(broker.userMap))

		// Remove dead clients
		//
		case userID := <-broker.deadUserIds:
			log.Printf("Client %d disconnected", userID)
			delete(broker.userMap, userID)
		}
	}
}
