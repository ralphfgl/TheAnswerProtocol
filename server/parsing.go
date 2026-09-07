package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GameWorld struct {
	World World `json:"world"`
}

type World struct {
	Locations []Location `json:"locations"`
	Items     []Item     `json:"items"`
	NPCs      []NPC      `json:"npcs"`
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
