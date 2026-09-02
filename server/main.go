// https://gobyexample.com/tcp-server

package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "strings"
    //"sync"
)

func main() {
    listener, err := net.Listen("tcp", ":8090")
    if err != nil {
        log.Fatal("Error listening:", err)
    }

    defer listener.Close()
	//    defer func() {
	// 	if err := listener.Close(); err != nil {
	// 		log.Printf("Error closing listener: %v", err)
	// 	}
	// }()

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
    defer   conn.Close()
	//    defer func() {
	// 	if err := conn.Close(); err != nil {
	// 		log.Printf("Error closing connection: %v", err)
	// 	}
	// }()

    fmt.Fprint(conn, "OK hello proto=1\n")
    reader := bufio.NewReader(conn)

    for {
        line, err := reader.ReadString('\n')
        if err != nil {
            return
        }
    line = strings.TrimSpace(line)
    fmt.Println("Received:", line)
    }
}

// to handle connection state (DISCONNECTED, CONNECTED, AUTHENTIFICATED< TERMINATED)
// type ConnectionState int
// const (
//     Connected ConnectionState = iota
//     Autenticated
//     Terminated
// )

// go map to store the playerbase -> we will need mutex cause all go routine will access this part
// players := make(map[string]net.Conn)

// RFC says : message = command-line / response-line / event-line (its either one)

// client connection should probably control writes (to avoid receiving message from many goroutine simultaneously)
// well need a Player (or client connection strucute)
// type Player struct {
//     Username    string
//     Conn        net.conn
//     Mu          sync.Mutex
// }
