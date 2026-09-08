package main

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"math/rand"
	"slices"
	"strings"

	"the_answer_protocol/common"
)

func (s *Server) handleConnect(player *Player, username string) {
	if player.State != Connected {
		s.sendError(player, 400, "INVALID_STATE")
		return
	}
	username = strings.TrimSpace(username)
	if username == "" {
		s.sendError(player, 400, "USERNAME_REQUIRED")
		return
	}
	// check if username in use
	s.Mu.Lock()
	defer s.Mu.Unlock()
	// map lookup in go returns 2 value, the actual value and a boolean hat tell if the key exist
	// comma separate the assignement from the condition
	if _, exists := s.players[username]; exists {
		s.sendError(player, 201, "NAME_IN_USE")
		return
	}
	// registration
	player.Username = username
	player.State = Authenticated
	player.CurrentRoom = "start"
	//player.Inventory
	//player.HP = 100

	s.players[username] = player

	s.sendResponse(player, "OK connected")
	s.logger.Info("Player authenticated: username=%s address=%s", username, player.Conn.RemoteAddr())
}

func (s *Server) handleQuit(player *Player) {
	s.sendResponse(player, "OK bye")
	log.Printf("Player %s quit", player.Username)
	// NOTE: the defer will clean up, this is redundunt
	player.Conn.Close()
}

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
	s.Mu.RLock()
	for _, player := range s.players {
		if player.CurrentRoom == room && player.Username != p.Username {
			playersInRoom = append(playersInRoom, player.Username)
		}
	}
	s.Mu.RUnlock()
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
		Type:    "room",
		Room:    roomInfo,
		Players: playersInRoom,
		Items:   itemIDs,
		NPCs:    npcIDs,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	s.sendResponse(p, "OK "+string(jsonData))
	return nil
}

func (s *Server) handleMove(p *Player, direction string) error {
	// NOTE: added against moving in combat
	p.Mu.Lock()
	inCombat := p.InCombat
	p.Mu.Unlock()
	if inCombat {
		s.sendError(p, 403, "CANNOT_MOVE_IN_COMBAT")
		return nil
	}
	var currentLocation *Location
	for i := range s.world.World.Locations {
		if s.world.World.Locations[i].Id == p.CurrentRoom {
			currentLocation = &s.world.World.Locations[i]
			break
		}
	}
	targetRoomID, exists := currentLocation.Exits[direction]
	if !exists {
		s.sendError(p, 301, "NO_EXIT")
		return nil
	}
	oldRoom := p.CurrentRoom
	s.Mu.Lock()
	p.CurrentRoom = targetRoomID
	s.Mu.Unlock()
	// FIX: add logg?
	// s.logger.Info("World state changed: player=%s move from=%s to=%s", p.Username, oldRoom, targetRoomID)
	s.broadcastRoomEvent(oldRoom, "EVT ROOM PRESENCE LEAVE "+p.Username)
	s.broadcastRoomEvent(targetRoomID, "EVT ROOM PRESENCE ENTER "+p.Username)
	s.sendResponse(p, fmt.Sprintf("OK room=%s", p.CurrentRoom))
	s.handleLook(p)
	return nil
}

func (s *Server) handleChat(p *Player, scope string, message string) error {
	event := fmt.Sprintf("EVT %s CHAT %s %s", scope, p.Username, message)
	switch scope {
	case "GLOBAL":
		s.broadcastAll(event)
	case "ROOM":
		s.broadcastRoomEvent(p.CurrentRoom, event)
	case "GROUP":
		if p.GroupID == "" {
			return fmt.Errorf("not in a group")
		}
		s.broadcastGroupEvent(p.GroupID, event)
	}
	s.sendResponse(p, "OK")
	return nil
}

func (s *Server) handleWho(p *Player) error {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	s.sendResponse(p, fmt.Sprintf("OK players=%d", len(s.players)))
	return nil
}

func (s *Server) handleGroup(p *Player, args []string) error {
	subCmd := args[0]
	rest := args[1:]
	switch subCmd {
	case "CREATE":
		return s.handleGroupCreate(p)
	case "INVITE":
		return s.handleGroupInvite(p, rest)
	case "JOIN":
		return s.handleGroupJoin(p, rest)
	case "LEAVE":
		return s.handleGroupLeave(p)
	case "DISPLAY":
		return s.handleGroupDisplay(p)
	default:
		return fmt.Errorf("unknown group subcommand: %s", subCmd)
	}
}

func (s *Server) handleGroupCreate(p *Player) error {
	if p.GroupID != "" {
		return fmt.Errorf("already in a group")
	}
	s.Mu.Lock()
	defer s.Mu.Unlock()
	groupID := fmt.Sprintf("group_%d", s.nextGroupID)
	s.nextGroupID++
	s.groups[groupID] = []string{p.Username}
	p.GroupID = groupID
	s.sendResponse(p, fmt.Sprintf("OK group=%s", groupID))
	if err := s.broadcastGroupList(); err != nil {
		return err
	}
	return nil
}

func (s *Server) handleGroupInvite(p *Player, args []string) error {
	if p.GroupID == "" {
		return fmt.Errorf("not in a group")
	}
	if len(args) < 1 {
		return fmt.Errorf("GROUP INVITE needs a username")
	}
	target := args[0]
	s.Mu.RLock()
	targetPlayer, exists := s.players[target]
	s.Mu.RUnlock()
	if !exists {
		return fmt.Errorf("player not found: %s", target)
	}
	if targetPlayer.GroupID != "" {
		return fmt.Errorf("player already in a group")
	}
	s.sendResponse(targetPlayer, fmt.Sprintf("EVT GROUP INVITE %s invited you to group %s", p.Username, p.GroupID))
	s.sendResponse(p, "OK")
	return nil
}

func (s *Server) handleGroupJoin(p *Player, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("GROUP JOIN needs a group ID")
	}
	if p.GroupID != "" {
		return fmt.Errorf("already in a group")
	}
	groupID := args[0]
	s.Mu.Lock()
	if _, exists := s.groups[groupID]; !exists {
		s.Mu.Unlock()
		return fmt.Errorf("group not found: %s", groupID)
	}
	s.groups[groupID] = append(s.groups[groupID], p.Username)
	p.GroupID = groupID
	s.Mu.Unlock()
	s.sendResponse(p, fmt.Sprintf("OK group=%s", groupID))
	s.broadcastGroupEvent(groupID, fmt.Sprintf("EVT JOIN %s joined the group", p.Username))
	if err := s.broadcastGroupList(); err != nil {
		return err
	}
	return nil
}

func (s *Server) handleGroupLeave(p *Player) error {
	if p.GroupID == "" {
		return fmt.Errorf("not in a group")
	}
	groupID := p.GroupID
	s.Mu.Lock()
	for i, name := range s.groups[groupID] {
		if name == p.Username {
			s.groups[groupID] = append(s.groups[groupID][:i], s.groups[groupID][i+1:]...)
			break
		}
	}
	if len(s.groups[groupID]) == 0 {
		delete(s.groups, groupID)
	}
	p.GroupID = ""
	s.Mu.Unlock()
	s.sendResponse(p, "OK")
	s.broadcastGroupEvent("EVT GROUP LEAVE %s left the group", p.Username)
	if err := s.broadcastGroupList(); err != nil {
		return err
	}
	return nil
}

func (s *Server) handleGroupDisplay(p *Player) error {
	groupEvent := common.GroupInfo{
		Type:      "group",
		GroupList: slices.Collect(maps.Keys(s.groups)),
	}
	jsonData, err := json.Marshal(groupEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal group list: %w", err)
	}
	s.sendResponse(p, string(jsonData))
	return nil
}

func (s *Server) handleTake(p *Player, itemRef string) error {
	// NOTE: could change the world structure to map instead of slice to access directly with key
	var currentLocation *Location
	for i := range s.world.World.Locations {
		if s.world.World.Locations[i].Id == p.CurrentRoom {
			currentLocation = &s.world.World.Locations[i]
			break
		}
	}
	var targetItemID string
	var targetItem Item
	// search in the room and in the world, itemID is the name while item is the object
	for _, itemID := range currentLocation.Items {
		for _, item := range s.world.World.Items {
			if item.Id == itemID {
				if item.Id == itemRef || strings.EqualFold(item.Name, itemRef) {
					targetItemID = itemID
					targetItem = item
					break
				}
			}
		}
		if targetItemID != "" {
			break
		}
	}
	// item not in the room
	if targetItemID == "" {
		s.sendError(p, 404, "ITEM_NOT_FOUND")
		return nil
		//return fmt.Errorf("item not found in room: %s", itemRef)
	}
	if !targetItem.Obtainable {
		return fmt.Errorf("item cannot be taken: %s", targetItem.Name)
	}
	s.Mu.Lock()
	// remove from the room
	for i, id := range currentLocation.Items {
		if id == targetItemID {
			currentLocation.Items = append(currentLocation.Items[:i], currentLocation.Items[i+1:]...)
			break
		}
	}
	p.Inventory = append(p.Inventory, targetItemID)
	s.Mu.Unlock()
	// FIX :add logging?
	// s.logger.Info("World state changed: player=%s picked up item=%s room=%s", p.Username, targetItemID, p.CurrentRoom)
	s.sendResponse(p, fmt.Sprintf("OK taken=%s", targetItemID))
	s.broadcastRoomEvent(p.CurrentRoom, fmt.Sprintf("EVT ROOM ITEM_TAKEN %s %s", p.Username, targetItemID))
	return nil
}

func (s *Server) handleInventory(p *Player) error {
	itemNames := []string{}
	for _, itemID := range p.Inventory {
		for _, item := range s.world.World.Items {
			if item.Id == itemID {
				itemNames = append(itemNames, item.Name)
				break
			}
		}
	}
	response := common.InventoryInfo{
		Type:  "inventory",
		Items: itemNames,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	s.sendResponse(p, "OK "+string(jsonData))
	return nil
}

func (s *Server) handleDrop(p *Player, itemRef string) error {
	var targetItemID string
	for _, itemID := range p.Inventory {
		for _, item := range s.world.World.Items {
			if item.Id == itemID {
				if item.Id == itemRef || strings.EqualFold(item.Name, itemRef) {
					targetItemID = itemID
					break
				}
			}
		}
		if targetItemID != "" {
			break
		}
	}
	if targetItemID == "" {
		s.sendError(p, 404, "ITEM_NOT_IN_INVENTORY")
		return nil
	}
	var currentLocation *Location
	for i := range s.world.World.Locations {
		if s.world.World.Locations[i].Id == p.CurrentRoom {
			currentLocation = &s.world.World.Locations[i]
			break
		}
	}
	s.Mu.Lock()
	for i, id := range p.Inventory {
		if id == targetItemID {
			p.Inventory = append(p.Inventory[:i], p.Inventory[i+1:]...)
			break
		}
	}
	currentLocation.Items = append(currentLocation.Items, targetItemID)
	s.Mu.Unlock()
	// FIX :add logging?
	// s.logger.Info("World state changed: player=%s dropped item=%s room=%s", p.Username, targetItemID, p.CurrentRoom)
	s.sendResponse(p, fmt.Sprintf("OK dropped=%s", targetItemID))
	s.broadcastRoomEvent(p.CurrentRoom, fmt.Sprintf("EVT ROOM ITEM_DROP %s %s", p.Username, targetItemID))
	return nil
}

func (s *Server) handleStatus(p *Player) error {
	response := common.StatusInfo{
		Type:   "status",
		HP:     p.HP,
		MaxHP:  p.MaxHP,
		Status: p.Status,
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	s.sendResponse(p, "OK "+string(jsonData))
	return nil
}

// func (s *Server) handleAttack(p *Player, npcRef string) error {
// 	return nil
// }

func (s *Server) handleTalk(p *Player, npcRef string) error {
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

	dialogueText := ""
	if len(targetNPC.Dialogue) > 0 {
		dialogueText = targetNPC.Dialogue[rand.Intn(len(targetNPC.Dialogue))]
	}

	response := common.TalkResponse{
		Type:     "talk",
		NPC:      targetNPC.Name,
		Dialogue: dialogueText,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	s.sendResponse(p, "OK "+string(jsonData))
	return nil
}
