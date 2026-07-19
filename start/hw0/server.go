package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	address  = ":8080"
	response = "OK\n"
)

func main() {
	ctx, stopFunc := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopFunc()

	var workGroup sync.WaitGroup

	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("error raised creating a listener:", err)
		return
	}

	go func() {
		<-ctx.Done()
		err := listener.Close()
		if err != nil {
			fmt.Println("error closing listener:", err)
		}
	}()

	for {
		connection, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():

			default:
				fmt.Println("error raised accepting a connection:", err)
			}
			break
		}

		workGroup.Go(func() {
			handleConnection(connection)
		})
	}

	fmt.Println("Waiting..")
	workGroup.Wait()
	fmt.Println("Closed")
}

func handleConnection(connection net.Conn) {
	defer func() {
		err := connection.Close()
		if err != nil {
			fmt.Println("error closing a connection:", err)
		}
	}()

	_, err := connection.Write([]byte(response))
	if err != nil {
		fmt.Println("error raised when writing a response:", err)
	}
}
