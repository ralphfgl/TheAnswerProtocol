package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"

	"the_answer_protocol/common"
)

func (s *Server) broadcastAll(message string) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	for _, player := range s.players {
		s.sendResponse(player, message)
	}
}

func (s *Server) broadcastRoomEvent(roomID string, message string) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	for _, player := range s.players {
		if player.CurrentRoom == roomID {
			s.sendResponse(player, message)
		}
	}
}

func (s *Server) broadcastGroupEvent(groupID string, message string) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	if groupID == "" {
		return
	}
	for _, username := range s.groups[groupID] {
		if player, exists := s.players[username]; exists {
			s.sendResponse(player, message)
		}
	}
}

func (s *Server) broadcastGroupList() error {
	groupEvent := common.GroupInfo{
		Type:      "group",
		GroupList: slices.Collect(maps.Keys(s.groups)),
	}
	jsonData, err := json.Marshal(groupEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal jsonData: %w", err)
	}
	s.broadcastAll(string(jsonData))
	return nil
}
