package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"strings"

	"the_answer_protocol/common"
)

func (s *Server) handleQuest(p *Player, npcRef string) error {
	npcRef = strings.TrimSpace(npcRef)
	var targetNPC *NPC
	s.Mu.RLock()
	currentRoom := p.CurrentRoom
	s.Mu.RUnlock()

	for _, loc := range s.world.World.Locations {
		if loc.Id == currentRoom {
			for _, sp := range loc.Spawns {
				for i, npc := range s.world.World.NPCs {
					if npc.Id == sp.NpcType {
						if strings.EqualFold(npc.Id, npcRef) || strings.EqualFold(npc.Name, npcRef) {
							targetNPC = &s.world.World.NPCs[i]
							break
						}
					}
				}
				if targetNPC != nil {
					break
				}
			}
			break
		}
	}
	if targetNPC == nil {
		s.sendError(p, 404, "NPC_NOT_FOUND")
		return nil
	}
	if !targetNPC.QuestGiver {
		s.sendError(p, 406, "NO_QUEST_AVAILABLE")
		return nil
	}
	fmt.Println("target NPC : ", targetNPC)
	questID := targetNPC.QuestID
	status, exists := p.PlayerQuests[questID]
	if !exists {
		p.PlayerQuests[questID] = "active"
		status = "active"
	} else if status == "active" {
		s.sendError(p, 406, "QUEST_ALREADY_ACCEPTED")
		return nil
	} else if status == "completed" {
		s.sendError(p, 406, "QUEST_ALREADY_COMPLETED")
		return nil
	}
	response := common.QuestResponse{
		Type:   "quest",
		Quest:  s.world.World.Quests[questID],
		Status: status,
	}
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	s.sendResponse(p, "OK "+string(data))
	return nil
}

func (s *Server) progressQuest(p *Player, eventType, target string) {
	// p.Mu.Lock()
	// defer p.Mu.Unlock()
	for id, state := range p.PlayerQuests {
		if state == "" {
			continue
		}
		q := s.world.World.Quests[id]
		if q.Type != eventType || q.Target != target {
			continue
		}
		p.PlayerQuests[id] = "completed"
		if q.RewardItem != "" {
			p.Inventory = append(p.Inventory, q.RewardItem)
			for i, id := range p.Inventory {
				if id == q.Target {
					p.Inventory = append(p.Inventory[:i], p.Inventory[i+1:]...)
					break
				}
			}
		}
		s.sendResponse(p, fmt.Sprintf("EVT QUEST %s completed! Reward: %s", q.Title, q.RewardItem))
	}
}

func (s *Server) handleQuests(p *Player) error {
	quests := make(map[string]string, len(p.PlayerQuests))
	maps.Copy(quests, p.PlayerQuests)
	response := common.QuestsResponse{
		Type:     "quests",
		QuestMap: quests,
		Count:    len(quests),
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	s.sendResponse(p, "OK "+string(jsonData))
	return nil
}
