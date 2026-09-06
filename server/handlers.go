package main

import (
	"encoding/json"
	"fmt"
	"the_answer_protocol/common"
)

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
	var playersInRoom []string
	s.Mu.RLock()
	for _, player := range s.players {
		if player.CurrentRoom == room && player.Username != p.Username {
			playersInRoom = append(playersInRoom, player.Username)
		}
	}
	s.Mu.RUnlock()
	var itemIDs []string
	for _, itemID := range currentLocation.Items {
		found := false
		for _, item := range s.world.World.Items {
			if item.Id == itemID {
				found = true
				break
			}
		}
		if found {
			itemIDs = append(itemIDs, itemID)
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
			npcIDs = append(npcIDs, spawn.NpcType)
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

func (s *Server) handleMove(p *Player, direction string) error {
	var currentLocation *Location
	for i := range s.world.World.Locations {
		if s.world.World.Locations[i].Id == p.CurrentRoom {
			currentLocation = &s.world.World.Locations[i]
			break
		}
	}
	targetRoomID, exists := currentLocation.Exits[direction]
	if !exists {
		return fmt.Errorf("no exit in direction: %s", direction)
	}
	oldRoom := p.CurrentRoom
	s.Mu.Lock()
	p.CurrentRoom = targetRoomID
	s.Mu.Unlock()
	// NOTE: verify the ABNF
	s.broadcastRoomEvent(oldRoom, "EVT ROOM PRESENCE LEAVE "+p.Username)
	s.broadcastRoomEvent(targetRoomID, "EVT ROOM PRESENCE ENTER "+p.Username)
	s.handleLook(p)
	return nil
}

func (s *Server) handleChat(p *Player, scope string, message string) error {
	event := fmt.Sprintf("EVT %s CHAT %s %s", scope, p.Username, message)
	switch scope {
	case "GLOBAL":
		s.broadcastAll(event)
	case "ROOM":
		s.broadcastRoomEvent(p.CurrentRoom, event)
	case "GROUP":
		if p.GroupID == "" {
			return fmt.Errorf("not in a group")
		}
		s.broadcastGroupEvent(p.GroupID, event)
	}
	s.sendResponse(p, "OK")
	return nil
}

func (s *Server) handleWho(p *Player) error {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	s.sendResponse(p, fmt.Sprintf("OK players=%d", len(s.players)))
	return nil
}
