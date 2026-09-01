// https://gobyexample.com/tcp-server

package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "strings"
)

func main() {
    parsing("data.json")
    listener, err := net.Listen("tcp", ":8090")
    if err != nil {
        log.Fatal("Error listening:", err)
    }

    defer func() {
		if err := listener.Close(); err != nil {
			log.Printf("Error closing listener: %v", err)
		}
	}()

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println("Error accepting conn:", err)
            continue
        }
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}()

    reader := bufio.NewReader(conn)
    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            log.Printf("Read error: %v", err)
            return
        }

        answer := strings.TrimSpace(message)
        response := fmt.Sprintf("Response: %s\n", answer)
        _, err = conn.Write([]byte(response))
        if err != nil {
            log.Printf("Server write error: %v", err)
        }
    }

}

//connect to server "netcat localhost 8090"
//free the port "kill -9 $(lsof -t -i:8090)"