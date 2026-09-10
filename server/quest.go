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
	p.Mu.Lock()
	currentRoom := p.CurrentRoom
	p.Mu.Unlock()
	var targetNPC *NPC
	s.Mu.RLock()
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
	var (
		npcQuestGiver bool
		npcQuestID    string
		questDef      Quest
		questExists   bool
	)
	if targetNPC != nil {
		npcQuestGiver = targetNPC.QuestGiver
		npcQuestID = targetNPC.QuestID
		questDef, questExists = s.world.World.Quests[npcQuestID]
	}
	s.Mu.RUnlock()
	if targetNPC == nil {
		s.sendError(p, 404, "NPC_NOT_FOUND")
		return nil
	}
	if !npcQuestGiver {
		s.sendError(p, 406, "NO_QUEST_AVAILABLE")
		return nil
	}
	if !questExists {
		s.sendError(p, 500, "QUEST_DEFINITION_MISSING")
		return nil
	}
	if questDef.Requires != "" {
		p.Mu.Lock()
		prereqStatus, prereqExists := p.PlayerQuests[questDef.Requires]
		p.Mu.Unlock()
		if !prereqExists || prereqStatus != "completed" {
			s.sendError(p, 406, "QUEST_PREREQUISITE_NOT_MET")
			return nil
		}
	}

	p.Mu.Lock()
	status, exists := p.PlayerQuests[npcQuestID]
	var (
		errCode int
		errMsg  string
	)
	if !exists {
		p.PlayerQuests[npcQuestID] = "active"
		status = "active"
	} else if status == "active" {
		errCode, errMsg = 406, "QUEST_ALREADY_ACCEPTED"
	} else if status == "completed" {
		errCode, errMsg = 406, "QUEST_ALREADY_COMPLETED"
	} else if status == "abandoned" {
		errCode, errMsg = 406, "QUEST_ALREADY_ABANDONED"
	} else {
		p.PlayerQuests[npcQuestID] = "active"
		status = "active"
	}
	p.Mu.Unlock()
	if errCode != 0 {
		s.sendError(p, errCode, errMsg)
		return nil
	}
	response := common.QuestResponse{
		Type:   "quest",
		Quest:  questDef,
		Status: status,
	}
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	s.handleQuests(p)
	s.logger.Info("QUEST_ACCEPT player=%s quest=%s npc=%s room=%s", p.Username, npcQuestID, targetNPC.Name, currentRoom)
	s.sendResponse(p, "OK "+string(data))
	return nil
}

func (s *Server) progressQuest(p *Player, eventType, target string) {
	type msg struct{ title, reward string }
	var msgs []msg
	for id, state := range p.PlayerQuests {
		if state != "active" {
			continue
		}
		q, exists := s.world.World.Quests[id]
		if !exists || q.Type != eventType || q.Target != target {
			continue
		}
		p.PlayerQuests[id] = "completed"
		s.logger.Info("QUEST_COMPLETE player=%s quest=%s type=%s target=%s reward=%s", p.Username, id, q.Type, q.Target, q.RewardItem)
		if q.RewardItem != "" {
			s.logger.Info("QUEST_REWARD player=%s quest=%s item=%s", p.Username, id, q.RewardItem)
			p.Inventory = append(p.Inventory, q.RewardItem)
			for i, itemID := range p.Inventory {
				if itemID == q.Target {
					p.Inventory = append(p.Inventory[:i], p.Inventory[i+1:]...)
					break
				}
			}
		}
		msgs = append(msgs, msg{q.Title, q.RewardItem})
		s.broadcastGroupEvent(p.CurrentRoom, fmt.Sprintf("EVT QUEST %s completed! Reward: %s", q.Title, q.RewardItem))
		response := common.QuestResponse{
			Type:   "quest",
			Quest:  q,
			Status: "test",
		}
		s.handleQuests(p)
		if data, err := json.Marshal(response); err == nil {
			s.sendResponse(p, "OK "+string(data))
		}
	}
	if len(msgs) > 0 {
		s.handleInventory(p)
		s.handleQuests(p)
	}
}

func (s *Server) handleQuests(p *Player) error {
	p.Mu.Lock()
	quests := make(map[string]string, len(p.PlayerQuests))
	maps.Copy(quests, p.PlayerQuests)
	p.Mu.Unlock()
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

func (s *Server) handleAbandonQuest(p *Player, questID string) error {
	p.Mu.Lock()
	questState, exists := p.PlayerQuests[questID]
	if !exists {
		p.Mu.Unlock()
		s.sendError(p, 404, "QUEST_NOT_FOUND")
		return nil
	}
	if questState == "completed" {
		p.Mu.Unlock()
		s.sendError(p, 406, "QUEST_ALREADY_COMPLETED")
		return nil
	}
	if questState == "abandoned" {
		p.Mu.Unlock()
		s.sendError(p, 406, "QUEST_ALREADY_ABANDONED")
		return nil
	}
	p.PlayerQuests[questID] = "abandoned"
	p.Mu.Unlock()
	questDef := s.world.World.Quests[questID]
	response := common.QuestResponse{
		Type:   "quests",
		Quest:  questDef,
		Status: "abandoned",
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	s.sendResponse(p, "OK "+string(jsonData))
	//s.sendResponse(p, fmt.Sprintf("OK quest abandoned=%s", questID))
	s.broadcastRoomEvent(p.CurrentRoom, fmt.Sprintf("EVT ROOM QUEST %s abandoned quest: %s", p.Username, questID))
	s.logger.Info("Quest abandoned: player=%s quest=%s", p.Username, questID)
	return nil
}
