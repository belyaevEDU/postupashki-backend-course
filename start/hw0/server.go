package main

import (
	"fmt"
	"net"
)

const (
	address  = ":8080"
	response = "OK\n"
)

func main() {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("error raised creating a listener: ", err)
		return
	}

	defer func() {
		err := listener.Close()
		if err != nil {
			fmt.Println("error closing listener: ", err)
		}
	}()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("error raised accepting a connection: ", err)
			continue
		}

		go handleConnection(connection)
	}
}

func handleConnection(connection net.Conn) {
	defer func() {
		err := connection.Close()
		if err != nil {
			fmt.Println("error closing a connection: ", err)
		}
	}()

	_, err := connection.Write([]byte(response))
	if err != nil {
		fmt.Println("error raised when writing a response: ", err)
	}
}
