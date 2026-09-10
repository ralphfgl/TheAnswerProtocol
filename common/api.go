// Package common: API for client and server
package common

type RoomInfo struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Exits       map[string]string `json:"exits"`
}

type LookResponse struct {
	Type    string   `json:"type"`
	Room    RoomInfo `json:"room"`
	Players []string `json:"players"`
	Items   []string `json:"items"`
	NPCs    []string `json:"npcs"`
}

type InventoryInfo struct {
	Type  string   `json:"type"`
	Items []string `json:"items"`
}

type GroupInfo struct {
	Type      string   `json:"type"`
	GroupList []string `json:"group_list"`
	MyGroup string `json:"my_group"`
}

type StatusInfo struct {
	Type   string `json:"type"`
	HP     int    `json:"hp"`
	MaxHP  int    `json:"max_hp"`
	Status string `json:"status"`
}

type CombatResponse struct {
	Type       string `json:"type"`
	AttackerHP int    `json:"attacker_hp"`
	TargetHP   int    `json:"target_hp"`
	Atk        int    `json:"player_damage"`
	CounterAtk int    `json:"npc_damage"`
	Status     string `json:"status"`
}

type TalkResponse struct {
	Type     string `json:"type"`
	NPC      string `json:"npc"`
	Dialogue string `json:"dialogue"`
}

type QuestResponse struct {
	Type   string `json:"type"`
	Quest  any    `json:"quest"`
	Status string `json:"status"`
}

type QuestsResponse struct {
	Type     string            `json:"type"`
	QuestMap map[string]string `json:"quest_map"`
	Count    int               `json:"count"`
}
