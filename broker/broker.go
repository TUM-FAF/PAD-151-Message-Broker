package broker

import (
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
	messages        string
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

func (broker *Broker) HandleChannels() {
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
		}
	}
}
