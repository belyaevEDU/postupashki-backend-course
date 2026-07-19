package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

const (
	address    = "localhost:8080"
	response   = "OK\n"
	trimCutset = "\n "
)

func main() {
	connection, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("error raised when dialing the address:", err)
		return
	}

	defer func() {
		err := connection.Close()
		if err != nil {
			fmt.Println("error raised when closing the connection:", err)
		}
	}()

	reader := bufio.NewReader(connection)
	buffer, err := io.ReadAll(reader)

	if err != nil {
		fmt.Println("error raised when reading: ")
	}

	if string(buffer) == response {
		fmt.Println("The server returned the expected response:", strings.TrimRight(response, trimCutset))
	} else {
		fmt.Println("The server returned an unexpected response:", strings.TrimRight(string(buffer), trimCutset))
	}
}
