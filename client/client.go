package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// Client ...
type Client struct {
	Protocol    string
	Address     string
	connection  net.Conn
	incomingMsg chan string
	outcominMsg chan string
}

// Init ...
func (c *Client) Init(protocol string, address string) {
	c.Protocol = protocol
	c.Address = address
	conn, err := net.Dial(c.Protocol, c.Address)
	c.connection = conn
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c.incomingMsg = make(chan string)
	c.outcominMsg = make(chan string)
}

// Run ...
func (c *Client) Run() {
	go c.listenIncomingMessages()
	go c.listenOutcomingMessages()

	for {
		select {
		case newIncomingMessage := <-c.incomingMsg:
			c.handleIncomingMessage(newIncomingMessage)
		case newOutcomindMessage := <-c.outcominMsg:
			c.handleOutcomingMessage(newOutcomindMessage)
		}
	}

}

func (c *Client) handleIncomingMessage(message string) {
	fmt.Println("incoming msg: " + message)
}

func (c *Client) handleOutcomingMessage(message string) {
	fmt.Println("outcoming msg: " + message)
	c.connection.Write([]byte(message))
}

func (c *Client) send(data string) {

}

func (c *Client) get() string {
	return ""
}

func (c *Client) listenIncomingMessages() {
	reader := bufio.NewReader(c.connection)
	c.listen(reader, c.incomingMsg)
}

func (c *Client) listenOutcomingMessages() {
	reader := bufio.NewReader(os.Stdin)
	c.listen(reader, c.outcominMsg)
}

func (c *Client) listen(reader *bufio.Reader, channel chan string) {
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err, message)
		}
		channel <- message
	}
}
