package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
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
	Username        string
	Conn            net.Conn
	State           ConnectionState
	Mu              sync.Mutex
	Writer          *bufio.Writer
	CurrentRoom     string
	GroupID         string
	Inventory       []string
	Attack          int
	Defense         int
	HP              int
	MaxHP           int
	Status          string
	InCombat        bool
	CombatTarget    string
	CmdWindowStart  time.Time
	CmdInWindow     int
	PlayerQuests    map[string]string
	ChatWindowStart time.Time
	ChatInWindow    int
	maxInventory    int
}

type Server struct {
	players           map[string]*Player
	Mu                sync.RWMutex
	cmdRegistry       *CommandRegistry
	playerLocations   map[string]string
	world             *GameWorld
	groups            map[string][]string
	nextGroupID       int
	logger            *Logger
	connectionMu      sync.Mutex
	recentConnections []time.Time
	maxTotalConn      int
}

func NewServer(worldFile string, logger *Logger) (*Server, error) {
	world, err := parsing(worldFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load world: %w", err)
	}
	// FIX: add validation
	s := &Server{
		players:      make(map[string]*Player),
		groups:       make(map[string][]string),
		world:        &world,
		logger:       logger,
		maxTotalConn: 100,
	}
	if err := validate(world); err != nil {
		s.logger.Error("Validation error: %v", err)
		return nil, fmt.Errorf("failed to validate the world: %w", err)
	}
	s.cmdRegistry = NewCommandRegistry(s)
	return s, nil
}

func (s *Server) handleCommand(player *Player, line string) {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return
	}
	commandName := parts[0]
	args := parts[1:]
	s.logger.Info("Command received: player=%s command=%s args=%s", player.Username, commandName, args)
	s.checkCommandFlood(player)
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
	logger := NewLogger()
	server, err := NewServer("../data.json", logger)
	if err != nil {
		logger.Error("Failed to initialize the server: %v", err)
		os.Exit(1)
	}
	port := ":8090"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		//log.Fatal("Error listening:", err)
		logger.Error("Failed to start TCP listener: %v", err)
		os.Exit(1)
	}
	defer listener.Close()
	logger.Info("TAP Server starting on port %s", port[1:])
	for {
		conn, err := listener.Accept()
		if err != nil {
			//log.Println("Error accepting conn:", err)
			logger.Error("Error accepting connection", err)
			continue
		}
		go server.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	s.Mu.Lock()
	if len(s.players) >= s.maxTotalConn {
		s.Mu.Unlock()
		s.logger.Warn("LIMIT: server full")
		conn.Write([]byte("ERR 503 SERVER_FULL\n"))
		conn.Close()
		return
	}
	s.Mu.Unlock()
	s.checkRapidConnections()
	address := conn.RemoteAddr().String()
	s.logger.Info("Client connection opened from %s", address)
	player := &Player{
		Conn:         conn,
		State:        Connected,
		Writer:       bufio.NewWriter(conn),
		HP:           100,
		MaxHP:        100,
		Attack:       12,
		Defense:      10,
		Status:       "healthy",
		InCombat:     false,
		PlayerQuests: make(map[string]string),
		maxInventory: 10,
	}
	defer func() {
		s.logger.Info("Client disconected: player=%s address=%s", player.Username, address)
		conn.Close()
		s.removePlayer(player)
	}()
	s.sendResponse(player, "OK hello proto=1")

	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			s.logger.Info("Connection closed for %s: %v", player.Username, err)
			return
		}
		if len(line) > 104 {
			s.sendError(player, 413, "MESSAGE_TOO_LONG")
			s.logger.Warn("LIMIT: message too long, player=%s length=%d", player.Username, len(line))
			return
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		s.handleCommand(player, line)
	}
}

func (s *Server) sendResponse(player *Player, message string) {
	s.logger.Info("Response sent: player=%s address=%s response=%s", player.Username, player.Conn.RemoteAddr(), message)
	player.Mu.Lock()
	defer player.Mu.Unlock()
	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}
	if _, err := player.Writer.WriteString(message); err != nil {
		s.logger.Error("Failed to write response: player=%s error=%v", player.Username, err)
		return
	}
	if err := player.Writer.Flush(); err != nil {
		s.logger.Error("Failed to flush response: player=%s error=%v", player.Username, err)
		return
	}
}

func (s *Server) sendError(player *Player, code int, message string) {
	s.logger.Warn("Error response sent: player=%s code=%d message=%s", player.Username, code, message)
	s.sendResponse(player, fmt.Sprintf("ERR %03d %s", code, message))
}

func (s *Server) removePlayer(player *Player) {
	if player.Username != "" {
		s.Mu.Lock()
		delete(s.players, player.Username)
		playerCount := len(s.players)
		s.Mu.Unlock()
		s.broadcastAll(fmt.Sprintf("EVT STATS players=%d", playerCount))
		s.logger.Info("Player %s removed", player.Username)
	}
}

func (s *Server) checkCommandFlood(p *Player) {
	now := time.Now()
	p.Mu.Lock()
	defer p.Mu.Unlock()
	if now.Sub(p.CmdWindowStart) >= time.Second {
		p.CmdWindowStart = now
		p.CmdInWindow = 0
	}
	p.CmdInWindow++
	if p.CmdInWindow > 20 {
		s.logger.Warn("Possible command flooding: player=%s command_per_minutes=%d", p.Username, p.CmdInWindow)
	}
}

func (s *Server) checkRapidConnections() {
	s.connectionMu.Lock()
	defer s.connectionMu.Unlock()
	now := time.Now()
	s.recentConnections = append(s.recentConnections, now)
	var recent []time.Time
	for _, t := range s.recentConnections {
		if now.Sub(t) <= time.Minute {
			recent = append(recent, t)
		}
	}
	s.recentConnections = recent
	if len(recent) > 2 {
		s.logger.Warn("Possible rapid connection pattern: connections_last_minut=%d", len(recent))
	}
}

func (s *Server) checkChatFlood(p *Player) bool {
	now := time.Now()
	p.Mu.Lock()
	defer p.Mu.Unlock()

	if now.Sub(p.ChatWindowStart) >= time.Minute {
		p.ChatWindowStart = now
		p.ChatInWindow = 0
	}
	p.ChatInWindow++
	if p.ChatInWindow > 10 {
		s.logger.Warn("LIMIT: Chat rate limit exceeded: player=%s count=%d", p.Username, p.ChatInWindow)
		return false
	}
	return true
}
