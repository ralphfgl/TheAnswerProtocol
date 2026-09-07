package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
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
	Username string
	Conn     net.Conn
	State    ConnectionState
	Mu       sync.Mutex
	// bufio.Writer add a buffer on top of an underlying io.Writer
	Writer      *bufio.Writer
	CurrentRoom string
	GroupID     string
	Inventory   []string
	Attack      int
	Defense     int
	HP          int
	MaxHP       int
	Status      string
	InCombat    bool
}

type Server struct {
	players         map[string]*Player
	Mu              sync.RWMutex
	cmdRegistry     *CommandRegistry
	playerLocations map[string]string
	world           *GameWorld
	groups          map[string][]string // each key a string, each value a slice
	nextGroupID     int
}

// constructor, create a server instance
// mutex has a zero value and is already usable
// we use a struct literal, no malloc is needed

func NewServer(worldFile string) (*Server, error) {
	world, err := parsing(worldFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load world: %w", err)
	}
	s := &Server{
		players: make(map[string]*Player),
		groups:  make(map[string][]string),
		world:   &world,
	}
	s.cmdRegistry = NewCommandRegistry(s)
	// NOTE: could add some logging about loading success
	return s, nil
}

func (s *Server) handleCommand(player *Player, line string) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return
	}
	commandName := strings.ToUpper(parts[0])
	args := parts[1:]
	cmd, exists := s.cmdRegistry.commands[commandName]
	if !exists {
		s.sendError(player, 400, fmt.Sprintf("UNKNOWN_COMMAND: %s", commandName))
		return
	}
	if len(args) < cmd.MinArgs {
		s.sendError(player, 400, fmt.Sprintf("TOO_FEW_ARGS: Need at least %d arguments", cmd.MinArgs))
		return
	}
	if cmd.MaxArgs > 0 && len(args) > cmd.MaxArgs {
		s.sendError(player, 400, fmt.Sprintf("TOO_MANY_ARGS: Maximum %d arguments allowed", cmd.MaxArgs))
		return
	}
	if cmd.Validator != nil {
		if err := cmd.Validator(args); err != nil {
			s.sendError(player, 400, fmt.Sprintf("INVALID_ARGS: %s", err.Error()))
			return
		}
	}
	if cmd.RequiresAuth && player.State != Authenticated {
		s.sendError(player, 401, "NOT_AUTHENTICATED")
		return
	}
	if err := cmd.Handler(player, args); err != nil {
		s.sendError(player, 500, fmt.Sprintf("COMMAND_ERROR: %s", err.Error()))
		return
	}
}

func main() {
	server, err := NewServer("../data.json")
	if err != nil {
		log.Println("Failed to initialize the server: ", err)
		os.Exit(1)
	}
	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer listener.Close()
	// NOTE: change print to a dynamic value
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
		Conn:     conn,
		State:    Connected,
		Writer:   bufio.NewWriter(conn),
		HP:       100,
		MaxHP:    100,
		Attack:   55,
		Defense:  40,
		Status:   "healthy",
		InCombat: false,
	}
	// clean up on exit
	defer func() {
		conn.Close()
		s.removePlayer(player)
	}()
	// send greetings
	s.sendResponse(player, "OK hello proto=1")
	//s.sendResponse(player, jsonTest)

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
		// NOTE: add parsing of the command here before handling
		s.handleCommand(player, line)
	}
}

// NOTE: general sendResponse and then wrapper for error, event
func (s *Server) sendResponse(player *Player, message string) {
	player.Mu.Lock()
	defer player.Mu.Unlock()
	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}
	// NOTE: add error checking on both write string and flush
	player.Writer.WriteString(message)
	player.Writer.Flush()
}

func (s *Server) sendError(player *Player, code int, message string) {
	s.sendResponse(player, fmt.Sprintf("ERR %03d %s", code, message))
}

func (s *Server) removePlayer(player *Player) {
	if player.Username != "" {
		s.Mu.Lock()
		delete(s.players, player.Username)
		s.Mu.Unlock()
		// NOTE: add timestamp and ip address
		log.Printf("Player %s removed", player.Username)
	}
}
