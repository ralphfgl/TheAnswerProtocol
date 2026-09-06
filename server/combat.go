package main

type PlayerStatus int

const (
	Healthy PlayerStatus = iota
	Poisoned
	Buff
	Debuff
)

type CombatState int

const (
	InCombat CombatState = iota
	Idle
)
