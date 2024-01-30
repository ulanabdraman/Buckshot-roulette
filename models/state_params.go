package models

type UserState int

const (
	Idle UserState = iota
	InLobby
	GivingLevel
	GivingCode
	InGame
)
