package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GameWorld struct {
	World World
}

type World struct {
	Locations []Location
	Items     []Item
	NPCs      []NPC
}

type Spawn struct {
	NpcType string
	Count   int
}

type Location struct {
	Id          string
	Name        string
	Description string
	Exits       map[string]string
	Spawns      []Spawn
	Items       []string
}

type Item struct {
	Id          string
	Name        string
	Description string
	Obtainable  bool
}

type NPC struct {
	Id          string
	Name        string
	Description string
	Dialogue    []string
	Stats       map[string]int
}

func parsing(fileName string) (GameWorld, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return GameWorld{}, fmt.Errorf("Error reading file:", err)
	}
	var game GameWorld
	err = json.Unmarshal(data, &game)
	if err != nil {
		return GameWorld{}, fmt.Errorf("Error parsing:", err)
	}
	// for _, loc := range game.World.Locations {
	// 	fmt.Printf("[%s] %s: %s\n", loc.Id, loc.Name, loc.Description)
	// }
	// for _, item := range game.World.Items {
	// 	fmt.Printf("[%s] %s (%t)\n", item.Id, item.Name, item.Obtainable)
	// }
	// for _, npc := range game.World.NPCs {
	// 	fmt.Printf("[%s] %s (HP: %d)\n", npc.Id, npc.Name, npc.Stats["hp"])
	// }
	return game, nil
}
