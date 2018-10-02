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
    "bufio"
    "strconv"
)

func main() {
    var brokerIPName string
    var brokerPort int
    flag.StringVar(&brokerIPName, "ip", "127.0.0.1", "Specify broker's IP")
    flag.IntVar(&brokerPort, "p", 8080, "Specify broker's port" )
    flag.Parse()
    
    var brokerIP = net.ParseIP(brokerIPName)

    if brokerIP == nil {
        fmt.Println("Error! Invalid broker IP.")
        os.Exit(1)
    } else {
        fmt.Printf("%s %d\n", brokerIP, brokerPort)
    }



  conn, err := net.Dial("tcp", brokerIP.String() + ":" + strconv.Itoa(brokerPort))
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
  for { 
    // read in input from stdin
    reader := bufio.NewReader(os.Stdin)
    fmt.Print(">>")
    text, err := reader.ReadString('\n')
    if err != nil {
      fmt.Println(err)
    }
    // send to socket
    fmt.Fprintf(conn, text + "\n")
    // listen for reply
    message, err := bufio.NewReader(conn).ReadString('\n')
        if err != nil {
      fmt.Println(err)
    }
  }
}