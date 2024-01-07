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

    requestLines := strings.Split(string(buf), "\r\n")
    requestLine := requestLines[0]
    requestParts := strings.Split(requestLine, " ")
    _, path, _ := requestParts[0], requestParts[1], requestParts[2]

    if path == "/user-agent" {
        userAgent := ""
        for _, line := range requestLines {
            if strings.HasPrefix(line, "User-Agent:") {
                userAgent = strings.TrimSpace(strings.TrimPrefix(line, "User-Agent:"))
                break
            }
        }

        if userAgent != "" {
            response := fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s", len(userAgent), userAgent)
            _, writeErr := conn.Write([]byte(response))
            if writeErr != nil {
                fmt.Println("Error writing:", writeErr.Error())
            }
        } else {
            response := "HTTP/1.1 400 Bad Request\r\n\r\n"
            _, writeErr := conn.Write([]byte(response))
            if writeErr != nil {
                fmt.Println("Error writing:", writeErr.Error())
            }
        }
    } else {
        response := "HTTP/1.1 404 Not Found\r\n\r\n"
        _, writeErr := conn.Write([]byte(response))
        if writeErr != nil {
            fmt.Println("Error writing:", writeErr.Error())
        }
    }
}