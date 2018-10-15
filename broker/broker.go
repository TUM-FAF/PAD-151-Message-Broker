package broker

import (
	"PAD-151-Message-Broker/model"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
)

// Broker is type representing the message broker
type Broker struct {
	clientCount       int
	userMap           map[int]*User
	newConnections    chan net.Conn
	deadUserIds       chan int
	messages          chan model.SentMessageModel
	broadcastMessages chan model.SentMessageModel
	p2pMessages       chan model.SentMessageModel
	subscribers       map[int][]*User
	postMessages      chan model.SentMessageModel
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
	broker.messages = make(chan model.SentMessageModel)

	// Channel for all messages that will be sent in broadcast
	//
	broker.broadcastMessages = make(chan model.SentMessageModel)

	// Channel for all messages that will be sent to specific clients
	//
	broker.p2pMessages = make(chan model.SentMessageModel)

	//Map of all subscribers
	//
	broker.subscribers = make(map[int][]*User)

	// Channel for all messages that will be sent to subscribers
	broker.postMessages = make(chan model.SentMessageModel)
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
			user.id = rand.Int()
			var responseModel model.ConnectionModel
			responseModel.Rooms = nil
			responseModel.YourID = user.id

			users := make(map[int]string)
			for _, v := range broker.userMap {
				users[v.id] = v.name
			}
			responseModel.Users = nil

			r, _ := model.EncodeYamlMessage(responseModel)

			conn.Write(r)

			broker.userMap[broker.clientCount] = user

			broker.clientCount++

			go getMessages(user, broker.messages, broker.deadUserIds)

		// Accept messages from connected client
		//
		case sentMessageModel := <-broker.messages:
			command := DispatchMessage(sentMessageModel, broker)
			command.Execute()

		// Remove dead clients
		//
		case userID := <-broker.deadUserIds:
			log.Printf("Client %d disconnected", userID)
			delete(broker.userMap, userID)
		}
	}
}
