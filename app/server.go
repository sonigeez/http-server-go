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
    n, err := conn.Read(buf)
    if err != nil {
        fmt.Println("Error reading:", err.Error())
        return
    }
    buf = buf[:n] // Trim the buffer to actual read size

    requestLine := strings.Split(string(buf), "\r\n")[0]
    requestParts := strings.Split(requestLine, " ")

    _, path, _ := requestParts[0], requestParts[1], requestParts[2]

    var response string
    switch {
    case strings.HasPrefix(path, "/echo/"):
        message := strings.TrimPrefix(path, "/echo/")
        response = buildResponse(200, "text/plain", message)

    case path == "/":
        response = buildResponse(200, "text/plain", "Welcome!")

    case strings.HasPrefix(path, "/user-agent"):
        userAgent := strings.Split(string(buf), "\r\n")[2]
        response = buildResponse(200, "text/plain", userAgent)

    default:
        response = buildResponse(404, "text/plain", "Not Found")
    }

    _, writeErr := conn.Write([]byte(response))
    if writeErr != nil {
        fmt.Println("Error writing:", writeErr.Error())
    }
}

func buildResponse(statusCode int, contentType string, body string) string {
	return fmt.Sprintf("HTTP/1.1 %d OK\r\n"+
		"Content-Type: %s\r\n"+
		"Content-Length: %d\r\n"+
		"Connection: close\r\n"+
		"\r\n"+
		"%s", statusCode, contentType, len(body), body)
}