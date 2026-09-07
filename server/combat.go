package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"the_answer_protocol/common"
)

func subtractOrZero(a, b int) int {
	if res := a - b; res > 0 {
		return res
	}
	return 0
}

func (s *Server) handleAttack(p *Player, npcRef string) error {
	var currentLocation *Location
	for i := range s.world.World.Locations {
		if s.world.World.Locations[i].Id == p.CurrentRoom {
			currentLocation = &s.world.World.Locations[i]
			break
		}
	}
	var targetNpcID string
	var targetNPC NPC
	for _, spawn := range currentLocation.Spawns {
		npcType := spawn.NpcType
		for _, npc := range s.world.World.NPCs {
			if npc.Id == npcType {
				if npc.Id == npcRef || strings.EqualFold(npc.Name, npcRef) {
					targetNpcID = spawn.NpcType
					targetNPC = npc
					break
				}
			}
		}
		if targetNpcID != "" {
			break
		}
	}
	if targetNpcID == "" {
		s.sendError(p, 404, "NPC_NOT_FOUND")
		return nil
		//return fmt.Errorf("item not found in room: %s", itemRef)
	}
	if !targetNPC.Hostile {
		s.sendError(p, 405, "NPC_NOT_HOSTILE")
		return nil
	}
	npcHP := targetNPC.Stats["hp"]
	if npcHP <= 0 {
		return fmt.Errorf("NPC is already defeated")
	}
	p.InCombat = true
	playerDamage := subtractOrZero(p.Attack, targetNPC.Stats["defense"])
	npcHP -= playerDamage
	targetNPC.Stats["hp"] = npcHP
	npcDamage := 0
	if npcHP > 0 {
		npcDamage = subtractOrZero(targetNPC.Stats["attack"], p.Defense)
		p.HP -= npcDamage
	}
	response := common.CombatResponse{
		AttackerHP: p.HP,
		TargetHP:   npcHP,
		Atk:        playerDamage,
		CounterAtk: npcDamage,
		Status:     "combat",
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	s.sendResponse(p, "OK "+string(jsonData))
	// s.broadcastRoomEvent(p.CurrentRoom, fmt.Sprintf("EVT ROOM COMBAT %s attacked %s", p.Username, targetNPC.Name))
	//
	// // Check if NPC is defeated
	// if npcHP <= 0 {
	// 	s.broadcastRoomEvent(p.CurrentRoom, fmt.Sprintf("EVT ROOM COMBAT %s defeated %s!", p.Username, targetNPC.Name))
	// }
	//
	// Check if player is defeated
	if p.HP <= 0 {
		p.HP = p.MaxHP / 2
		oldRoom := p.CurrentRoom
		p.CurrentRoom = "start"

		s.broadcastRoomEvent(oldRoom, fmt.Sprintf("EVT ROOM PRESENCE LEAVE %s", p.Username))
		s.broadcastRoomEvent("start", fmt.Sprintf("EVT ROOM PRESENCE ENTER %s", p.Username))
		s.sendResponse(p, "OK You lost and respawned at the start.")
	}

	return nil
}
