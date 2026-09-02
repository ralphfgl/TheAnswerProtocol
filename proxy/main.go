package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer wsConn.Close()
	fmt.Println("React frontend connected to proxy")
	tcpConn, err := net.Dial("tcp", "localhost:8090")
	if err != nil {
		fmt.Println("Error connecting to TCP server:", err)
		wsConn.WriteMessage(websocket.TextMessage, []byte("Server is offline.\n"))
		return
	}
	defer tcpConn.Close()
	fmt.Println("Proxy connected to TCP server")
	go func() {
		tcpReader := bufio.NewReader(tcpConn)
		for {
			serverMessage, err := tcpReader.ReadString('\n')
			if err != nil {
				fmt.Println("TCP server disconnected.")
				wsConn.Close()
				return
			}
			err = wsConn.WriteMessage(websocket.TextMessage, []byte(serverMessage))
			if err != nil {
				return
			}
		}
	}()
	for {
		_, reactMessage, err := wsConn.ReadMessage()
		if err != nil {
			fmt.Println("React frontend disconnected")
			break
		}
		msgForServer := string(reactMessage) + "\n"
		_, err = tcpConn.Write([]byte(msgForServer))
		if err != nil {
			fmt.Println("Error writing to TCP server:", err)
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("Proxy WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting proxy server:", err)
	}
}