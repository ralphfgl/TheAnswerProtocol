package main

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
