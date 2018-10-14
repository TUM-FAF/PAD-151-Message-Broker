package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

// Client ...
type Client struct {
	Protocol   string
	Address    string
	Connection net.Conn
}

// Init ...
func (c *Client) Init(protocol string, address string) {
	c.Protocol = protocol
	c.Address = address
	var err error
	c.Connection, err = net.Dial(c.Protocol, c.Address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Run ...
func (c *Client) Run() {
	for {
		// read in input from stdin
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">>")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		// send to socket
		fmt.Fprintf(c.Connection, text+"\n")
		// listen for reply
		message, err := bufio.NewReader(c.Connection).ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(message)
	}
}
