package client

import (
	"PAD-151-Message-Broker/model"
	"PAD-151-Message-Broker/notification"
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/fatih/color"
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
	// fmt.Print("Enter your name: ")
	infoC := color.New(color.FgCyan, color.Bold)
	comC := color.New(color.FgYellow)
	infoC.Printf("Enter your name: ")
	reader := bufio.NewReader(os.Stdin)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err, message)
	}
	message = strings.TrimSuffix(message, "\n")
	c.sendConnectionRequest(message)
	data := c.getConnectionResponse()
	c.id = data.YourID
	// fmt.Printf("Your ID: %v\nRooms: %v\nUsers: %v\n", data.YourID, data.Rooms, data.Users)
	magenta := color.New(color.FgMagenta).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	infoC.Printf("Your ID: %v\n", magenta(data.YourID))
	infoC.Printf("Users: %v\n", red(data.Users))
	infoC.Printf("Commands:\n")
	comC.Println("broadcast:	/b [msg]")
	comC.Println("publish:	/p [msg]")
	comC.Println("subscribe:	/s [user ID]")
	comC.Println("p2p:		/u [user ID] [msg]")
	comC.Println("Connected:	/c")
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
	notification.PlayNotification()
	m := model.ResponseMessageModel{}
	err := model.DecodeYamlMessage([]byte(message), &m)
	if err != nil {
		fmt.Println(err)
	}
	senderC := color.New().SprintFunc()
	messageC := color.New(color.FgYellow, color.Italic).SprintFunc()
	if m.SenderID%2 == 0 {
		senderC = color.New(color.FgGreen, color.Bold, color.Italic).SprintFunc()
	} else {
		senderC = color.New(color.FgRed, color.Bold, color.Italic).SprintFunc()
	}
	fmt.Printf("%s %s\n", senderC(m.SenderName+":"), messageC(m.Message))
}

func (c *Client) handleOutcomingMessage(message string) {
	mp := MessageParser{}
	m, err := mp.Parse(message, c)
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
	dangerC := color.New(color.FgRed, color.Bold)
	for {
		s, err := reader.Read()
		if err != nil {
			dangerC.Println("Broker shutdown!!!")
			c.connection.Close()
			os.Exit(1)
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
