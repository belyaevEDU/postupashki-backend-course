package main

import (
	"fmt"
	"net"
)

const (
	address  = "localhost:8080"
	response = "OK\n"
)

func main() {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("error raised when dialing the address:", err)
		return
	}

	defer connection.Close()

	buffer := make([]byte, 3)
	connection.Read(buffer)

	if string(buffer) == response {
		fmt.Println("The server returned the expected response:", response)
	} else {
		fmt.Println("The server returned an unexpected response:", string(buffer))
	}
}
