package common

// API types that both server and client use
// These are the "contract" between client and server

type RoomInfo struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Exits       map[string]string `json:"exits"`
}

type LookResponse struct {
	Room    RoomInfo `json:"room"`
	Players []string `json:"players"`
	Items   []string `json:"items"`
	NPCs    []string `json:"npcs"`
}

// Future API types can go here
type InventoryInfo struct {
	Items       []string `json:"items"`
	Gold        int      `json:"gold"`
	MaxCapacity int      `json:"max_capacity"`
}

type StatusInfo struct {
	HP         int `json:"hp"`
	MaxHP      int `json:"max_hp"`
	Level      int `json:"level"`
	Experience int `json:"experience"`
}
