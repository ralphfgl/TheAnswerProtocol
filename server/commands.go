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
		// place holder
	}
	// and so on
}

// Handler implementation
func (s *Server) handleLook(p *Player) error {
	room := p.CurrentRoom
	if room == "" {
		room = "start"
	}
	var currentRoom *Location
	for i, location := range s.world.World.Locations {
		if location.Id == room {
			currentRoom = &s.world.World.Locations[i]
			break
		}
	}
	var playersInRoom []string
	s.mu.RLock()
	for _, player := range s.players {
		if player.CurrentRoom == room && player.Username != p.Username {
			playersInRoom = append(playersInRoom, player.Username)
		}
	}
	s.mu.RUnlock()
	response := fmt.Sprintf(`{"room":{"id":"%s","name":"%s","description":"%s","exits":{`,
		currentRoom.Id, currentRoom.Name, currentRoom.Description)
	first := true
	for dir, target := range currentRoom.Exits {
		if !first {
			response += ","
		}
		response += fmt.Sprintf(`"%s":"%s"`, dir, target)
		first = false
	}
	response += `},"players":[`
	for i, name := range playersInRoom {
		if i > 0 {
			response += ","
		}
		response += fmt.Sprintf(`"%s"`, name)
	}
	for i, item := range currentRoom.Items {
		if i > 0 {
			response += ","
		}
		response += fmt.Sprintf(`"%s"`, item)
	}
	response += `], "npcs":[`
	for i, spawn := range currentRoom.Spawns {
		if i > 0 {
			response += ","
		}
		response += fmt.Sprintf(`"%s"`, spawn.NpcType)
	}
	response += `]}}`
	s.sendResponse(p, "OK "+ response)
	return nil
}
