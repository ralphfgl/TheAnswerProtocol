package main

import (
	"encoding/json"
	"fmt"
	"the_answer_protocol/common"
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

func (s *Server) handleLook(p *Player) error {
	room := p.CurrentRoom
	if room == "" {
		room = "start"
	}
	var currentLocation *Location
	for i, location := range s.world.World.Locations {
		if location.Id == room {
			currentLocation = &s.world.World.Locations[i]
			break
		}
	}
	// if currentLocation == nil {
	// 	return fmt.Errorf("room not found: %s", room)
	// }
	var playersInRoom []string
	s.mu.RLock()
	for _, player := range s.players {
		if player.CurrentRoom == room && player.Username != p.Username {
			playersInRoom = append(playersInRoom, player.Username)
		}
	}
	s.mu.RUnlock()
	var itemIDs []string
	for _, itemID := range currentLocation.Items {
		// Verify item exists
		found := false
		for _, item := range s.world.World.Items {
			if item.Id == itemID {
				found = true
				break
			}
		}
		if found {
			itemIDs = append(itemIDs, itemID) // Send ID, not name
		}
	}
	var npcIDs []string
	for _, spawn := range currentLocation.Spawns {
		// Verify NPC exists
		found := false
		for _, npc := range s.world.World.NPCs {
			if npc.Id == spawn.NpcType {
				found = true
				break
			}
		}
		if found {
			npcIDs = append(npcIDs, spawn.NpcType) // Send ID, not name
		}
	}
	roomInfo := common.RoomInfo{
		Id:          currentLocation.Id,
		Name:        currentLocation.Name,
		Description: currentLocation.Description,
		Exits:       currentLocation.Exits,
	}
	response := common.LookResponse{
		Room:    roomInfo,
		Players: playersInRoom,
		Items:   itemIDs,
		NPCs:    npcIDs,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	s.sendResponse(p, "OK"+string(jsonData))
	return nil
}
