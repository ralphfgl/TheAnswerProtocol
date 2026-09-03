package main

import (
	"fmt"
)

type Command struct {
	Name    string
	MinArgs int
	MaxArgs int
	// Description	string
	Validator func([]string) error
	Handler   func(*Player, []string) error
	// NOTE: maybe change the design
	RequiresAuth bool
}

// registry with all the command available
type CommandRegistry struct {
	commands map[string]*Command
}

// command registry constructor
func NewCommandRegistry() *CommandRegistry {
	cr := &CommandRegistry{
		commands: make(map[string]*Command),
	}
	cr.registerCommands()
	return cr
}

func (cr *CommandRegistry) registerCommands() {
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
		// NOTE: no handler cause connect is a special function
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
	}
	cr.commands["LOOK"] = &Command{
		Name:         "LOOK",
		MinArgs:      0,
		MaxArgs:      0,
		RequiresAuth: true,
		Validator: func(args []string) error {
			if len(args) > 0 {
				// NOTE: logg
				return fmt.Errorf("LOOK takes no arguments")
			}
			return nil
		},
		Handler: func(p *Player, args []string) error {
			// place holder for look logic
			return nil
		},
	}
	cr.commands["MOVE"] = &Command{
		// place holder
	}
	// and so on
}
