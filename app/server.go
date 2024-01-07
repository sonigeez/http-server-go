package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue // Continue to the next iteration of the loop.
		}

		go func(c net.Conn) {
			defer c.Close()
			fmt.Printf("Serving %s\n", c.RemoteAddr().String())
			response := "HTTP/1.1 200 OK\r\n\r\n"
			c.Write([]byte(response))
		}(conn)
	}
}
