package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type ConnectionState int

// type and expression is inherited from previous line
// iota is like auto() from enum in python
const (
	Disconnected ConnectionState = iota
	Connected
	Authenticated
	Terminated
)

type Player struct {
	Username	string
	Conn		net.Conn
	State		ConnectionState
	Mu			sync.Mutex
	// bufio.Writer add a buffer on top of an underlying io.Writer
	Writer		*bufio.Writer
}

type Server struct {
	players	map[string]*Player // declare a map with string keys and value of *Player type (pointer to player struct)
	mu		sync.RWMutex
}

// constructor, create a server instance
// mutex has a zero value and is already usable
// we use a struct literal, no malloc is needed
func NewServer() Server {
	return &Server{
		players: make(map[string]*Player),
	}
}

func main() {
	server := NewServer()

	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer listener.Close()
	log.Println("TAP Server starting on :8090")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting conn:", err)
			continue
		}
		go server.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	// create player in CONNECTED state
	player := &Player{
		Conn:	conn,
		State:	Connected,
		Writer:	bufio.NewWriter(conn),
	}
	// clean up on exit
	defer func() {
		conn.Close()
		s.removePlayer(player)
	}()
	// send greetings
	s.sendResponse(player, "OK hello proto=1")

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Connection closed for %s: %v", player.Username, err)
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		s.handleCommand(player, line)
	}
}

func (s *Server) handleCommand(player *Player, line string) {
	// func SplitN(s, sep string, n int) []string -> s the string to split, n max number of piece to return (if negative split all occurences. Returns a slice of substring
	parts := strings.SplitN(line, " ", 2)
	command := strings.ToUpper(parts[0])
	var args string
	if len(parts) > 1 {
		args = parts[1]
	}
	switch command {
	case "CONNECT":
		s.handleConnect(player, args)
	case "QUIT":
		s.handleQuit(player)
	default:
		if player.State != Authenticated {
			s.sendError(player, 401, "NOT_AUTENTICATED")
			return
		}
		s.handleCommandAuthenticated(player, command, args)
	}
}

func (s *Server) handleConnect(player *Player, username string) {
	if player.State != Connected {
		s.sendError(player, 400, "INVALID_STATE")
		return
	}
	username = strings.TrimSpace(username)
	if username == "" {
		s.sendError(player, 400, "USERNAME_REQUIRED")
		return
	}
	// check if username in use
	s.mu.Lock()
	defer s.mu.Unlock()
	// map lookup in go returns 2 value, the actual value and a boolean hat tell if the key exist
	// comma separate the assignement from the condition
	if _, exists := s.players[username]; exists {
		s.sendError(player, 201, "NAME_IN_USE")
		return
	}
	// registration
	player.Username = username
	player.State = Authenticated
	s.player[username] = player
	s.sendResponse(player, "OK connected")
	log.Printf("Player %s connected", username)
}

func (s *Server) handleQuit(player *Player) {
	s.sendResponse(player, "OK goodbye")
	log.Printf("Player %s quit", player.Username)
	player.Conn.Close()
}

func (s *Server) handleCommandAuthenticated(player *Player, command string, args string) {
	switch command {
	case "LOOK":
	case "MOVE":
	case "CHAT":
	case "WHO":
	case "GROUP CREATE":
	case "GROUP INVITE":
	case "GROUP JOIN":
	case "GROUP LEAVE":
	case "TAKE":
	case "DROP":
	case "INVENTORY":
	case "TALK":
	case "ATTACK":
	case "STATUS":
	case "QUEST":
	case "QUESTS":

	}
}
