package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
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
	p.Mu.Lock()
	if p.InCombat {
		if p.CombatTarget != npcRef {
			p.Mu.Unlock()
			s.sendError(p, 409, "ALREADY_IN_COMBAT_WITH_ANOTHER NPC")
			return nil
		}
	}
	p.Mu.Unlock()

	s.Mu.Lock()
	defer s.Mu.Unlock()
	var currentLocation *Location
	for i := range s.world.World.Locations {
		if s.world.World.Locations[i].Id == p.CurrentRoom {
			currentLocation = &s.world.World.Locations[i]
			break
		}
	}
	var targetNpcID string
	var targetNPC *NPC
	var targetNPCIndex int
	for _, spawn := range currentLocation.Spawns {
		npcType := spawn.NpcType
		for i, npc := range s.world.World.NPCs {
			if npc.Id == npcType {
				if npc.Id == npcRef || strings.EqualFold(npc.Name, npcRef) {
					targetNpcID = spawn.NpcType
					targetNPC = &s.world.World.NPCs[i]
					targetNPCIndex = i
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
	}
	if !targetNPC.Hostile {
		s.sendError(p, 405, "NPC_NOT_HOSTILE")
		return nil
	}
	p.Mu.Lock()
	if !p.InCombat {
		p.InCombat = true
		p.CombatTarget = targetNpcID
		p.Status = "combat"
	}
	npcHP := targetNPC.Stats["hp"]
	playerDamage := subtractOrZero(p.Attack, targetNPC.Stats["defense"])
	if playerDamage == 0 {
		playerDamage = 1
	}
	targetNPC.Stats["hp"] -= playerDamage
	npcDefeated := targetNPC.Stats["hp"] <= 0
	npcDamage := 0
	if !npcDefeated {
		npcDamage = subtractOrZero(targetNPC.Stats["attack"], p.Defense)
		if npcDamage == 0 {
			npcDamage = 1
		}
		p.HP -= npcDamage
	}
	playerDefeated := p.HP <= 0
	s.world.World.NPCs[targetNPCIndex].Stats["hp"] = targetNPC.Stats["hp"]
	response := common.CombatResponse{
		Type:       "combat",
		AttackerHP: p.HP,
		TargetHP:   npcHP,
		Atk:        playerDamage,
		CounterAtk: npcDamage,
		Status:     "combat",
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		p.Mu.Unlock()
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	s.sendResponse(p, "OK "+string(jsonData))
	s.broadcastRoomEvent(p.CurrentRoom, fmt.Sprintf("EVT ROOM COMBAT %s attacked %s for %d damage", p.Username, targetNPC.Name, playerDamage))
	if npcDamage > 0 && !playerDefeated {
		s.broadcastRoomEvent(p.CurrentRoom, fmt.Sprintf("EVT ROOM COMBAT %s counter-attacked %s for %d damage", targetNPC.Name, p.Username, npcDamage))

	}
	s.logger.Info("Combat: player=%s npc=%s damage=%d counter=%d npc_hp=%d player_hp=%d", p.Username, targetNPC.Name, playerDamage, npcDamage, targetNPC.Stats["hp"], p.HP)
	if npcDefeated {
		s.broadcastRoomEvent(p.CurrentRoom, fmt.Sprintf("EVT ROOM COMBAT %s defeated %s!", p.Username, targetNPC.Name))
		s.removeNPCFromRoom(p.CurrentRoom, targetNpcID)
		p.InCombat = false
		p.CombatTarget = ""
		p.Status = "healthy"
		p.Mu.Unlock()
		s.sendResponse(p, fmt.Sprintf("OK You defeated %s", targetNPC.Name))
		s.handleLook(p)
		return nil
	}
	if playerDefeated {
		p.HP = p.MaxHP / 2
		oldRoom := p.CurrentRoom
		p.CurrentRoom = "start"
		p.InCombat = false
		p.CombatTarget = ""
		p.Status = "healthy"
		p.Mu.Unlock()

		s.broadcastRoomEvent(oldRoom, fmt.Sprintf("EVT ROOM PRESENCE LEAVE %s", p.Username))
		s.broadcastRoomEvent("start", fmt.Sprintf("EVT ROOM PRESENCE ENTER %s", p.Username))
		s.sendResponse(p, "OK You lost and respawned at the start.")
		s.logger.Info("Combat ended: player=%s defeated by %s, respawned at start", p.Username, targetNPC.Name)
		s.handleLook(p)
		return nil
	}
	p.Mu.Unlock()
	s.sendResponse(p, "OK Combat in progress. Attack again.")
	return nil
}

func (s *Server) removeNPCFromRoom(roomID, npcType string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	for i := range s.world.World.Locations {
		if s.world.World.Locations[i].Id == roomID {
			for j, spawn := range s.world.World.Locations[i].Spawns {
				if spawn.NpcType == npcType {
					s.world.World.Locations[i].Spawns = append(
						s.world.World.Locations[i].Spawns[:j],
						s.world.World.Locations[i].Spawns[j+1:]...,
					)
					return
				}
			}
			break
		}
	}

}

func (s *Server) handleFlee(p *Player) error {
	p.Mu.Lock()
	if !p.InCombat {
		s.sendError(p, 400, "NOT_IN_COMBAT")
		return nil
	}
	if rand.Intn(100) < 75 {
		p.InCombat = false
		p.CombatTarget = ""
		p.Status = "healthy"
		room := p.CurrentRoom
		username := p.Username
		p.Mu.Unlock()
		s.sendResponse(p, "OK flee from combat")
		s.broadcastRoomEvent(room, fmt.Sprintf("EVT ROOM COMBAT %s fled from combat!", username))
		s.logger.Info("Player %s fled from combat", username)
		return nil
	}
	p.Mu.Unlock()
	s.sendResponse(p, "OK failed fleeing from combat")
	return nil
}
