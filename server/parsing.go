package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type GameWorld struct {
	World World `json:"world"`
}

type World struct {
	Locations []Location       `json:"locations"`
	Items     []Item           `json:"items"`
	NPCs      []NPC            `json:"npcs"`
	Quests    map[string]Quest `json:"quests"`
}

type Location struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Exits       map[string]string `json:"exits"`
	Spawns      []Spawn           `json:"spawns,omitempty"`
	Items       []string          `json:"items,omitempty"`
}

type Spawn struct {
	NpcType string `json:"npc_type"`
	Count   int    `json:"count"`
}

type Item struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Obtainable  bool   `json:"obtainable"`
}

type NPC struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Dialogue    []string       `json:"dialogue"`
	Stats       map[string]int `json:"stats"`
	Hostile     bool           `json:"hostile"`
	QuestGiver  bool           `json:"quest_giver"`
	QuestID     string         `json:"quest_id,omitempty"`
}

type Quest struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Giver            string `json:"giver"`
	Type             string `json:"type"`
	Target           string `json:"target"`
	Description      string `json:"description"`
	RewardItem       string `json:"reward_item"`
	DialogueStart    string `json:"dialogue_start,omitempty"`
	DialogueComplete string `json:"dialogue_complete"`
	Requires         string `json:"requires,omitempty"`
}

type ValidationErrors []string

func (e ValidationErrors) Error() string {
	return strings.Join(e, "\n")
}

func validate(g GameWorld) error {
	var errors ValidationErrors
	// Build lookup maps
	locationMap := make(map[string]Location)
	for _, loc := range g.World.Locations {
		locationMap[loc.Id] = loc
	}
	itemMap := make(map[string]Item)
	for _, item := range g.World.Items {
		itemMap[item.Id] = item
	}
	npcMap := make(map[string]NPC)
	for _, npc := range g.World.NPCs {
		npcMap[npc.Id] = npc
	}
	// 1. Validate all exits
	for _, loc := range g.World.Locations {
		for direction, targetId := range loc.Exits {
			if _, exists := locationMap[targetId]; !exists {
				errors = append(errors, fmt.Sprintf("Location '%s' has exit '%s' pointing to non-existent location '%s'",
					loc.Id, direction, targetId))
			}
		}
	}
	// 2. Validate spawns in locations
	for _, loc := range g.World.Locations {
		for _, spawn := range loc.Spawns {
			if _, exists := npcMap[spawn.NpcType]; !exists {
				errors = append(errors, fmt.Sprintf("Location '%s' spawns NPC '%s' which doesn't exist in NPC definitions",
					loc.Id, spawn.NpcType))
			}
		}
	}
	// 3. Validate items in locations
	for _, loc := range g.World.Locations {
		for _, itemId := range loc.Items {
			if _, exists := itemMap[itemId]; !exists {
				errors = append(errors, fmt.Sprintf("Location '%s' contains item '%s' which doesn't exist in item definitions",
					loc.Id, itemId))
			}
		}
	}
	// 4. Validate NPC references
	usedItems := make(map[string]bool) // Track where items are used (location or quest reward)
	usedNPCs := make(map[string]bool)  // Track where NPCs are used (location spawns or quest giver)
	for _, loc := range g.World.Locations {
		for _, itemId := range loc.Items {
			usedItems[itemId] = true
		}
		for _, spawn := range loc.Spawns {
			usedNPCs[spawn.NpcType] = true
		}
	}
	// 5. Validate quests
	for questId, quest := range g.World.Quests {
		// Check quest ID matches map key
		if questId != quest.ID {
			errors = append(errors, fmt.Sprintf("Quest map key '%s' doesn't match quest ID '%s'", questId, quest.ID))
		}
		// Check quest giver exists
		if _, exists := npcMap[quest.Giver]; !exists {
			errors = append(errors, fmt.Sprintf("Quest '%s' references non-existent NPC giver '%s'", questId, quest.Giver))
		} else {
			// Check that NPC has quest_giver flag set
			if !npcMap[quest.Giver].QuestGiver {
				errors = append(errors, fmt.Sprintf("NPC '%s' is referenced as quest giver for quest '%s' but has quest_giver=false",
					quest.Giver, questId))
			}
			// Check that NPC's quest_id matches
			if npcMap[quest.Giver].QuestID != questId {
				errors = append(errors, fmt.Sprintf("NPC '%s' has quest_id='%s' but is listed as giver for quest '%s'",
					quest.Giver, npcMap[quest.Giver].QuestID, questId))
			}
		}
		// Check target exists for quest type
		switch quest.Type {
		case "fetch_item":
			if _, exists := itemMap[quest.Target]; !exists {
				errors = append(errors, fmt.Sprintf("Quest '%s' has fetch_item target '%s' which doesn't exist in item definitions",
					questId, quest.Target))
			}
		case "defeat_npc":
			if _, exists := npcMap[quest.Target]; !exists {
				errors = append(errors, fmt.Sprintf("Quest '%s' has defeat_npc target '%s' which doesn't exist in NPC definitions",
					questId, quest.Target))
			}
		default:
			errors = append(errors, fmt.Sprintf("Quest '%s' has unknown type '%s'", questId, quest.Type))
		}
		// Check reward item exists
		if quest.RewardItem != "" {
			if _, exists := itemMap[quest.RewardItem]; !exists {
				errors = append(errors, fmt.Sprintf("Quest '%s' rewards item '%s' which doesn't exist in item definitions",
					questId, quest.RewardItem))
			} else {
				usedItems[quest.RewardItem] = true
			}
		}
	}
	// 6. Validate NPC quest references
	for _, npc := range g.World.NPCs {
		if npc.QuestGiver && npc.QuestID != "" {
			if _, exists := g.World.Quests[npc.QuestID]; !exists {
				errors = append(errors, fmt.Sprintf("NPC '%s' references quest '%s' which doesn't exist in quest definitions",
					npc.Id, npc.QuestID))
			}
		}
		if npc.QuestGiver && npc.QuestID == "" {
			errors = append(errors, fmt.Sprintf("NPC '%s' has quest_giver=true but no quest_id specified", npc.Id))
		}
	}
	// 7. Check for duplicate item IDs
	itemIdSet := make(map[string]bool)
	for _, item := range g.World.Items {
		if itemIdSet[item.Id] {
			errors = append(errors, fmt.Sprintf("Duplicate item ID found: '%s'", item.Id))
		}
		itemIdSet[item.Id] = true
	}
	// 8. Check for duplicate NPC IDs
	npcIdSet := make(map[string]bool)
	for _, npc := range g.World.NPCs {
		if npcIdSet[npc.Id] {
			errors = append(errors, fmt.Sprintf("Duplicate NPC ID found: '%s'", npc.Id))
		}
		npcIdSet[npc.Id] = true
	}
	// 9. Check for duplicate location IDs
	locIdSet := make(map[string]bool)
	for _, loc := range g.World.Locations {
		if locIdSet[loc.Id] {
			errors = append(errors, fmt.Sprintf("Duplicate location ID found: '%s'", loc.Id))
		}
		locIdSet[loc.Id] = true
	}
	// 10. Ensure quest reward items are unique (not already in world)
	for _, quest := range g.World.Quests {
		if quest.RewardItem != "" {
			// Check if reward item is already placed in a location
			itemInLocation := false
			for _, loc := range g.World.Locations {
				for _, itemId := range loc.Items {
					if itemId == quest.RewardItem {
						itemInLocation = true
						break
					}
				}
				if itemInLocation {
					break
				}
			}
			if itemInLocation {
				errors = append(errors, fmt.Sprintf("Quest reward item '%s' from quest '%s' is already present in a location (should be unique)",
					quest.RewardItem, quest.ID))
			}
		}
	}
	// 11. Check for bidirectional exits (optional but recommended)
	for _, loc := range g.World.Locations {
		for direction, targetId := range loc.Exits {
			// Check if the target location has an exit back
			if targetLoc, exists := locationMap[targetId]; exists {
				// Check if any exit from target points back to source
				hasReverseExit := false
				reverseDirection := ""
				for srcDir, srcTarget := range targetLoc.Exits {
					if srcTarget == loc.Id {
						hasReverseExit = true
						reverseDirection = srcDir
						break
					}
				}
				if !hasReverseExit {
					errors = append(errors, fmt.Sprintf("Location '%s' has exit '%s' to '%s' but no reverse exit exists",
						loc.Id, direction, targetId))
				} else {
					// Check if direction makes sense (optional warning)
					// This is just an informational check, not a hard error
					if direction == "north" && reverseDirection != "south" {
						fmt.Printf("Warning: Location '%s' exit north to '%s' but reverse is '%s' (expected 'south')\n",
							loc.Id, targetId, reverseDirection)
					}
				}
			}
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

func parsing(fileName string) (GameWorld, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return GameWorld{}, fmt.Errorf("Error reading file: %w", err)
	}
	var game GameWorld
	err = json.Unmarshal(data, &game)
	if err != nil {
		return GameWorld{}, fmt.Errorf("Error parsing: %w", err)
	}
	return game, nil
}
