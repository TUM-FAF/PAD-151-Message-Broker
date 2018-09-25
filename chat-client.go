/*
	chat-client.go -- A minimal chat client.
	Run this like
		> go run chat-client.go [-ip <broker-ip>] [-p <broker-port>]
*/

package main

import (
    "net"
	"fmt"
	"flag"
	"os"
)

func main() {
	var brokerIPName string
	var brokerPort int
	flag.StringVar(&brokerIPName, "ip", "127.0.0.1", "Specify broker's IP")
	flag.IntVar(&brokerPort, "p", 8080, "Specify broker's port" )
	flag.Parse()
	
	fmt.Println(brokerIPName)

	var brokerIP = net.ParseIP(brokerIPName)

	if brokerIP == nil {
		fmt.Println("Error! Invalid broker IP.")
		os.Exit(1)
	} else {
		fmt.Printf("%s %d\n",brokerIP, brokerPort)
	}
}