package main

import "PAD-151-Message-Broker/broker"

func main() {
	myBroker := broker.Broker{}
	myBroker.Init()
	myBroker.StartServer("0.0.0.0", "8080", "tcp")
}
