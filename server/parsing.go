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
