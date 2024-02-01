package main

import (
	"log"
	"net"
	"strings"
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

		if value.typ != "array" {
			log.Printf("Invalid type: %s. Expected array", value.typ)
		}
		if len(value.array) == 0 {
			log.Printf("Invalid array length: %d. Expected > 0", len(value.array))
		}

		cmd := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]
		writer := NewWriter(conn)

		handler, ok := Handlers[cmd]
		if !ok {
			log.Printf("Unknown command: %s", cmd)
			writer.Write(RedisMessage{typ: "string", str: "Unknown command"})
			continue
		}
		results := handler(args)
		writer.Write(results)
	}

}
