package models

type Lobby struct {
	Level   int
	Code    string
	Players []Player
	Game    bool
}

type Player struct {
	ChatID   int64
	Username string
	UserID   int
}
