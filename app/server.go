package main

import (
	"fmt"
	"net"
	"os"
	"strings"
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
			continue
		}
		go handleRequest(conn)
	}
}


func handleRequest(conn net.Conn) {
    defer conn.Close()

    buf := make([]byte, 1024)
    _, err := conn.Read(buf)
    if err != nil {
        fmt.Println("Error reading:", err.Error())
        return
    }

    requestLine := strings.Split(string(buf), "\r\n")[0]
    requestParts := strings.Split(requestLine, " ")
    if len(requestParts) < 3 {
        fmt.Println("Invalid request line")
        return
    }

    method, path, _ := requestParts[0], requestParts[1], requestParts[2]
    if method == "GET" && path == "/" {
        response := "HTTP/1.1 200 OK\r\n\r\n"
		fmt.Println("200")
        conn.Write([]byte(response))
    } else {
		fmt.Println("404")
        response := "HTTP/1.1 404 Not Found\r\n\r\n"
        conn.Write([]byte(response))
    }
}
