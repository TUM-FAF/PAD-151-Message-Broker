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
	id          int
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
	c.id = data.YourID
	fmt.Printf("Your ID: %v\nRooms: %v\nUsers: %v\n", data.YourID, data.Rooms, data.Users)
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
	mp := ResponseParser{}
	m, err := mp.Parse(message)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(m)
}

func (c *Client) handleOutcomingMessage(message string) {
	mp := MessageParser{}
	m, err := mp.Parse(message)
	if err != nil {
		return
	}
	b, err := model.EncodeJsonMessage(m)

	if err != nil {
		fmt.Print(err)
	}
	newline := make([]byte, 1)
	newline[0] = '\n'
	c.connection.Write(append(b, newline...))
}

func (c *Client) listenIncomingMessages() {
	reader := c.connection
	wrapper := BufferedReader{reader}
	c.listen(wrapper, c.incomingMsg)
}

func (c *Client) listenOutcomingMessages() {
	reader := bufio.NewReader(os.Stdin)
	wrapper := NewlineReader{reader}
	c.listen(wrapper, c.outcominMsg)
}

func (c *Client) listen(reader Reader, channel chan string) {
	for {
		s, err := reader.Read()
		if err != nil {
			fmt.Println(err)
		}
		channel <- string(s)
	}
}

func (c *Client) sendConnectionRequest(data string) {
	cp := ConnectionParser{}
	msg, err := cp.ParseRequest(data)
	if err != nil {
		fmt.Println(err)
	}
	newline := make([]byte, 1)
	newline[0] = '\n'
	c.connection.Write(append(msg.([]byte), newline...))
}

func (c *Client) getConnectionResponse() model.ConnectionModel {
	recvBuf := make([]byte, 1024)
	n, err := c.connection.Read(recvBuf[:])
	if err != nil {
		fmt.Println(err)
	}
	cp := ConnectionParser{}
	msg, err := cp.ParseResponse(string(recvBuf[:n]))
	if err != nil {
		fmt.Println(err)
	}
	return msg.(model.ConnectionModel)
}
