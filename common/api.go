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
	Items []string `json:"items"`
}

type StatusInfo struct {
	HP    int `json:"hp"`
	MaxHP int `json:"max_hp"`
	// NOTE: look a enum in combat.go
	// status need to be defined
	Status int
}

// NOTE: for example purpose
type CombatResponse struct {
	AttackerHP int `json:"attacker_hp"`
	TargetHP   int `json:"target_hp"`
	Damage     int `json:"npc_damage,omitempty`
}
