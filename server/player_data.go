package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// Persitent player data
// Go struc with JSON tags that defines how struct fields map to JSON field names during encoding/decoding
// it convert Go PacalCase/camelCase to snake_case for json
type PlayerData struct {
	Username	string		`json:"username"`
	Inventory	[]string	`json:"inventory"`
	//CurrentRoom
	//Stats
	//LastLog
}

// handle saving/loading player data
type PlayerDataManager struct {
	dataDir	string
	mu		sync.RWMutex
}

func NewPlayerDataManager(dataDir string) (*PlayerDataManager, error) {
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data dir: %w", err)
	}
	return &PlayerDataManager{dataDir: dataDir}, nil
}

func (pdm *PlayerDataManager) SavePlayer(player *Player) error {
	pdm.mu.Lock()
	defer pdm.mu.Unlock()
	data := &PlayerData{
		Username:	player.Username,
		Inventory:		player.Inventory,
		// add more 
	}
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal player data: %w", err)
	}
	filename := filepath.Join(pdm.dataDir, player.Username+".json")
	if err := ioutil.WriteFile(filename, jsonData, 0644): err != nil {
		return fmt.Errorf("failed to write player data: %w", err)
	}
	return nil
}
