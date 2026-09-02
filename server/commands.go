package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type Command struct {
	Name		string
	// requiresAuth
	// MinArgs
	// MaxArgs
	// validator
	Handler		func(*Player, []string) error
}

// registry with all the command available
type CommandRegistry struct {
	commands map[string]*Command
}

func NewCommandRegistry() *CommandRegistry {
	cr := &CommandRegistry {
		commands: make(map[string]*Command),
	}
	cr.registerCommands()
	return cr
}

func (cr *CommandRegistry) registerCommands() {
	// LOOK command
	cr.commands["LOOK"] = &Command{
		Name:		"LOOK",
	}
	Handler: func(p *Player, args []string) error {
		// logic
		return nil
	}
	cr.commands["MOVE"] = &Command{
		// 
	}
	// and so on
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
		s.sendError(player)
	}
}
