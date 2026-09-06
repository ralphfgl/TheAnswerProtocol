package main

func (s *Server) broadcastAll(message string) {

}

func (s *Server) broadcastRoomEvent(roomID string, message string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, player := range s.players {
		if player.CurrentRoom == roomID {
			s.sendResponse(player, message)
		}
	}
}
