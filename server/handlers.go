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
	s.mu.RLock()
	for _, player := range s.players {
		if player.CurrentRoom == room && player.Username != p.Username {
			playersInRoom = append(playersInRoom, player.Username)
		}
	}
	s.mu.RUnlock()
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
