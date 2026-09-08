package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"the_answer_protocol/common"
)

type QuestState struct {
	Status string `json:"status"`
}

type PlayerQuestData struct {
	Quests map[string]*QuestState `json:"quests"`
	mu     sync.RWMutex
}

func (s *Server) questByID(id string) *Quest {
	for i := range s.world.World.Quests {
		if s.world.World.Quests[i].ID == id {
			return &s.world.World.Quests[i]
		}
	}
	return nil
}

func (s *Server) questForNPC(npcID string) *Quest {
	for i := range s.world.World.Quests {
		if s.world.World.Quests[i].Giver == npcID {
			return &s.world.World.Quests[i]
		}
	}
	return nil
}

func (s *Server) handleQuest(p *Player, npcRef string) error {
	npc, err := s.findNPCInRoom(p.CurrentRoom, npcRef)
	if err != nil {
		s.sendError(p, 404, "NPC_NOT_FOUND")
		return nil
	}

	quest := s.questForNPC(npc.Id)
	if quest == nil {
		s.sendError(p, 406, "NO_QUEST_AVAILABLE")
		return nil
	}

	p.QuestData.mu.Lock()
	if p.QuestData.Quests == nil {
		p.QuestData.Quests = make(map[string]*QuestState)
	}
	if state := p.QuestData.Quests[quest.ID]; state != nil {
		p.QuestData.mu.Unlock()
		if state.Status == "completed" {
			s.sendError(p, 406, "QUEST_ALREADY_COMPLETED")
		} else {
			s.sendError(p, 406, "QUEST_ALREADY_ACCEPTED")
		}
		return nil
	}
	p.QuestData.Quests[quest.ID] = &QuestState{Status: "active"}
	p.QuestData.mu.Unlock()

	s.sendResponse(p, "OK "+mustJSON(common.QuestResponse{Type: "quest", QuestList: []common.Quests{toAPIQuest(*quest)}, Count: 1}))
	return nil
}

func (s *Server) handleQuests(p *Player, _ string) error {
	p.QuestData.mu.RLock()
	defer p.QuestData.mu.RUnlock()
	list := make([]common.Quests, 0, len(p.QuestData.Quests))
	for id := range p.QuestData.Quests {
		if q := s.questByID(id); q != nil {
			list = append(list, toAPIQuest(*q))
		}
	}
	s.sendResponse(p, "OK "+mustJSON(common.QuestResponse{Type: "quests", QuestList: list, Count: len(list)}))
	return nil
}

func (s *Server) handleAbandonQuest(p *Player, id string) error {
	p.QuestData.mu.Lock()
	defer p.QuestData.mu.Unlock()
	state := p.QuestData.Quests[id]
	if state == nil {
		s.sendError(p, 404, "QUEST_NOT_FOUND")
		return nil
	}
	if state.Status == "completed" {
		s.sendError(p, 406, "QUEST_ALREADY_COMPLETED")
		return nil
	}
	delete(p.QuestData.Quests, id)
	s.sendResponse(p, fmt.Sprintf("OK quest=%s abandoned", id))
	return nil
}

// progressQuest is the only progression path. TAKE and ATTACK call it with the event type and target.
func (s *Server) progressQuest(p *Player, eventType, target string) {
	p.QuestData.mu.Lock()
	defer p.QuestData.mu.Unlock()
	for id, state := range p.QuestData.Quests {
		if state.Status != "active" {
			continue
		}
		q := s.questByID(id)
		if q == nil || q.Type != eventType || q.Target != target {
			continue
		}
		state.Status = "completed"
		if q.RewardItem != "" {
			p.Inventory = append(p.Inventory, q.RewardItem)
		}
		s.sendResponse(p, fmt.Sprintf("EVT QUEST %s completed! Reward: %s", q.Title, q.RewardItem))
	}
}

func toAPIQuest(q Quest) common.Quests {
	return common.Quests{ID: q.ID, Title: q.Title, Giver: q.Giver, Type: q.Type, Target: q.Target, Description: q.Description, RewardItem: q.RewardItem, DialogueStart: q.DialogueStart, DialogueComplete: q.DialogueComplete}
}

func mustJSON(v any) string { b, _ := json.Marshal(v); return string(b) }
