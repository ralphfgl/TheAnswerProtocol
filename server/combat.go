package main

import (
	"fmt"

	"the_answer_protocol/common"
)

func subtractOrZero(a, b int) int {
	if res := a - b; res > 0 {
		return res
	}
	return 0
}

func (s *Server) handleAttack(p *Player, npcRef string) error {
	var targetNPC *NPC
	var targetNPCSpawn *Spawn

	for _, spawn := range s.world.World.Locations {
		if spawn.Id == p.CurrentRoom {
			for _, s := range spawn.Spawns {
				for _, npc := range s.world.World.NPCs {
					if npc.Id == s.NpcType {
						if npc.Id == npcRef || strings.EqualFold(npc.Name, npcRef) {
							targetNPC = &npc
							targetNPCSpawn = &s
							break
						}
					}
				}
				if targetNPC != nil {
					break
				}
			}
		}
		if targetNPC != nil {
			break
		}
	}

	if targetNPC == nil {
		return fmt.Errorf("NPC not found in room: %s", npcRef)
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
		if p.HP < 0 {
			// NOTE: restore 50 hp and at start loc
			p.CurrentRoom = "start"
			p.HP = 50
		}
	}
	response := common.CombatResponse{
		AttackerHP: p.HP,
		TargetHP:   npcHP,
		playerDmg:  playerDamage,
		Status:     "combat",
	}
	if npcDamage > 0 {
		response["npc_damage"] = npcDamage
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	s.sendResponse(p, "OK"+string(jsonData))

	// Broadcast combat event to room
	s.broadcastRoomEvent(p.CurrentRoom, fmt.Sprintf("EVT ROOM COMBAT %s attacked %s", p.Username, targetNPC.Name))

	// Check if NPC is defeated
	if npcHP <= 0 {
		s.broadcastRoomEvent(p.CurrentRoom, fmt.Sprintf("EVT ROOM COMBAT %s defeated %s!", p.Username, targetNPC.Name))
	}

	// Check if player is defeated
	if p.HP <= 0 {
		// Respawn at start with half HP
		p.HP = p.MaxHP / 2
		oldRoom := p.CurrentRoom
		p.CurrentRoom = "start"

		s.broadcastRoomEvent(oldRoom, fmt.Sprintf("EVT ROOM PRESENCE LEAVE %s", p.Username))
		s.broadcastRoomEvent("start", fmt.Sprintf("EVT ROOM PRESENCE ENTER %s", p.Username))
		s.sendResponse(p, "OK You have been defeated and respawned at the start.")
	}

	return nil
}
