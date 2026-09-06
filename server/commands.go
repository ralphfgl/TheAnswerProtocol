package main

import (
	"fmt"
	"strings"
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
	cr.commands["CHAT"] = &Command{
		Name:         "CHAT",
		MinArgs:      2,
		MaxArgs:      0,
		RequiresAuth: true,
		Validator: func(args []string) error {
			if len(args) < 2 {
				return fmt.Errorf("CHAT needs a scope and a message")
			}
			scope := args[0]
			if scope != "GLOBAL" && scope != "ROOM" && scope != "GROUP" {
				return fmt.Errorf("invalid scope: %s (use GLOBAL, ROOM or GROUP)", scope)
			}
			return nil
		},
		// NOTE: return directly scope and message
		Handler: func(p *Player, args []string) error {
			scope := args[0]
			message := strings.Join(args[1:], " ")
			return s.handleChat(p, scope, message)
		},
	}
	cr.commands["WHO"] = &Command{
		Name:         "WHO",
		MinArgs:      0,
		MaxArgs:      0,
		RequiresAuth: true,
		Validator: func(args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("WHO takes no arguments")
			}
			return nil
		},
		Handler: func(p *Player, args []string) error {
			return s.handleWho(p)
		},
	}
	cr.commands["GROUP"] = &Command{
		Name:         "GROUP",
		MinArgs:      1,
		MaxArgs:      0,
		RequiresAuth: true,
		Validator: func(args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("GROUP need a subcommand (CREATE, INVITE, JOIN, LEAVE)")
			}
			return nil
		},
		Handler: func(p *Player, args []string) error {
			return s.handleGroup(p, args)
		},
	}
	cr.commands["TAKE"] = &Command{
		Name:         "TAKE",
		MinArgs:      1,
		MaxArgs:      0,
		RequiresAuth: true,
		Validator: func(args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("TAKE needs an item ID")
			}
			return nil
		},
		Handler: func(p *Player, args []string) error {
			return s.handleTake(p, strings.Join(args, " "))
		},
	}
	cr.commands["INVENTORY"] = &Command{
		Name:         "INVENTORY",
		MinArgs:      0,
		MaxArgs:      0,
		RequiresAuth: true,
		Validator: func(args []string) error {
			if len(args) > 0 {
				return fmt.Errorf("INVENTORY takes no arguments")
			}
			return nil
		},
		Handler: func(p *Player, args []string) error {
			return s.handleInventory(p)
		},
	}
	cr.commands["DROP"] = &Command{
		Name:         "DROP",
		MinArgs:      1,
		MaxArgs:      0,
		RequiresAuth: true,
		Validator: func(args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("DROP needs an item ID")
			}
			return nil
		},
		Handler: func(p *Player, args []string) error {
			return s.handleDrop(p, strings.Join(args, " "))
		},
	}
}
