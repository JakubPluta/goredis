package main

import (
	"fmt"
	"log"
	"net"
)

const (
	serverPort = ":6379"
	okMsg      = "+OK\r\n"
)

func main() {
	log.Printf("Starting server on %s", serverPort)
	listener, err := net.Listen("tcp", serverPort)
	if err != nil {
		log.Printf("Error starting server: %s", err.Error())
		return

	}

	conn, err := listener.Accept()
	if err != nil {
		log.Printf("Error accepting connection: %s", err.Error())
		return
	}
	defer conn.Close()

	for {
		resp := NewResp(conn)
		value, err := resp.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %s", err.Error())
			return
		}
		fmt.Println(value)
		conn.Write([]byte(okMsg))
	}

}
