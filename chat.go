package main

import (
	"fmt"
	"net"
	"os"
)

const (
    CONN_HOST      = "0.0.0.0"
    CONN_PORT      = "8080"
    CONN_TYPE      = "tcp"
)

func main() {
	// Start the TCP server
	//
	server, err := net.Listen(CONN_TYPE, CONN_HOST + ":" + CONN_PORT)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Server started at", server.Addr())
	}
}
