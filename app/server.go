package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"path/filepath"
	"io/ioutil"
)

var dirFlag = flag.String("directory", ".", "directory to serve files from")



func main() {
	flag.Parse()
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
    buf = buf[:n]

    requestLine := strings.Split(string(buf), "\r\n")[0]
    requestParts := strings.Split(requestLine, " ")

    _, path, _ := requestParts[0], requestParts[1], requestParts[2]

    var response string
    switch {
    case strings.HasPrefix(path, "/echo/"):
        message := strings.TrimPrefix(path, "/echo/")
        response = buildResponse(200, "text/plain", message)
	case strings.HasPrefix(path, "/files/"):
		filePath := filepath.Join(*dirFlag, strings.TrimPrefix(path, "/files/"))
		//check if method is GET
		if requestParts[0] == "POST" {
			// get file from body and save it to the directory
			file, err := os.Create(filePath)
			if err != nil {
				response = buildResponse(500, "text/plain", "Internal Server Error")
			} else {
				defer file.Close()
				_, err := file.Write(buf)
				if err != nil {
					response = buildResponse(500, "text/plain", "Internal Server Error")
				} else {
					response = buildResponse(201, "text/plain", "OK")
				}
			}
			break
		}
		// Use os.Stat to check if the file exists and is not a directory before attempting to open it.
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				response = buildResponse(404, "text/plain", "File not found")
			} else {
				response = buildResponse(500, "text/plain", "Internal Server Error")
			}
		} else if fileInfo.IsDir() {
			response = buildResponse(400, "text/plain", "Bad Request")
		} else {
			file, err := os.Open(filePath)
			if err != nil {
				response = buildResponse(500, "text/plain", "Internal Server Error")
			} else {
				defer file.Close()
			// If the file exists, read its contents with ioutil.ReadFile or by creating a buffer and using os.ReadFile.
			fileContents, err := ioutil.ReadFile(filePath)
			if err != nil {
				response = buildResponse(500, "text/plain", "Internal Server Error")
			} else {
				response = buildResponse(200, "application/octet-stream", string(fileContents))
			}
		}
	}
    case path == "/":
        response = buildResponse(200, "text/plain", "Welcome!")

	case strings.HasPrefix(path, "/user-agent"):
	userAgent := strings.Split(string(buf), "\r\n")[2]
	userAgent = strings.TrimPrefix(userAgent, "User-Agent: ")
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