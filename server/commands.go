package main

import (
	"fmt"
)

type Command struct {
	Name         string
	MinArgs      int
	MaxArgs      int
	Validator    func([]string) error
	Handler      func(*Player, []string) error
	RequiresAuth bool
}

// registry with all the command available
type CommandRegistry struct {
	commands map[string]*Command
}

// command registry constructor
func NewCommandRegistry(s *Server) *CommandRegistry {
	cr := &CommandRegistry{
		commands: make(map[string]*Command),
	}
	cr.registerCommands(s)
	return cr
}

func (cr *CommandRegistry) registerCommands(s *Server) {
	cr.commands["CONNECT"] = &Command{
		Name:         "CONNECT",
		MinArgs:      1,
		MaxArgs:      1,
		RequiresAuth: false,
		Validator: func(args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("CONNECT needs exactly one username.")
			}
			return nil
		},
		Handler: func(p *Player, args []string) error {
			s.handleConnect(p, args[0])
			return nil
		},
	}
	cr.commands["QUIT"] = &Command{
		Name:         "QUIT",
		MinArgs:      0,
		MaxArgs:      0,
		RequiresAuth: false,
		Validator: func(args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("QUIT takes no arguments")
			}
			return nil
		},
		Handler: func(p *Player, args []string) error {
			s.handleQuit(p)
			return nil
		},
	}
	cr.commands["LOOK"] = &Command{
		Name:         "LOOK",
		MinArgs:      0,
		MaxArgs:      0,
		RequiresAuth: true,
		Validator: func(args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("LOOK takes no arguments")
			}
			return nil
		},
		Handler: func(p *Player, args []string) error {
			return s.handleLook(p)
		},
	}
	cr.commands["MOVE"] = &Command{
		Name:         "MOVE",
		MinArgs:      1,
		MaxArgs:      1,
		RequiresAuth: true,
		Validator: func(args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("MOVE need a direction")
			}
			return nil
		},
		Handler: func(p *Player, args []string) error {
			return s.handleMove(p, args[0])
		},
	}
	// and so on
}
