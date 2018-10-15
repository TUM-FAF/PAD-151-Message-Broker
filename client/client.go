package client

import (
	"PAD-151-Message-Broker/model"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
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

//Connect ...
func (c *Client) Connect() {
	fmt.Print("Enter your name: ")
	reader := bufio.NewReader(os.Stdin)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err, message)
	}
	message = strings.TrimSuffix(message, "\n")
	c.sendConnectionRequest(message)
	data := c.getConnectionResponse()
	fmt.Printf("data: %s", data)
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
			newOutcomindMessage = strings.TrimSuffix(newOutcomindMessage, "\n")
			c.handleOutcomingMessage(newOutcomindMessage)
		}
	}

}

func (c *Client) handleIncomingMessage(message string) {
	// go notification.PlayNotification()
	mp := MessageParser{}
	m, err := mp.Parse(message)
	fmt.Println(err)
	fmt.Println(m)

}

func (c *Client) handleOutcomingMessage(message string) {
	mp := MessageParser{}
	m, err := mp.Parse(message)
	if err != nil {
		return
	}
	b, err := model.EncodeJsonMessage(m)
	fmt.Println(string(b))

	if err != nil {
		fmt.Print(err)
	}
	c.connection.Write(b)
}

func (c *Client) sendConnectionRequest(data string) {
	cp := ConnectionParser{}
	msg, err := cp.ParseRequest(data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(msg.([]byte)))
	c.connection.Write(msg.([]byte))
}

func (c *Client) getConnectionResponse() string {
	reader := bufio.NewReader(c.connection)
	data, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
	}
	cp := ConnectionParser{}
	msg, err := cp.ParseResponse(data)
	if err != nil {
		fmt.Println(err)
	}
	return msg.(string)
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
