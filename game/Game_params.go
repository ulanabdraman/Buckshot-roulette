package game

type GamePlayer struct {
	Hp    int
	Items []string
	block bool
}
type GameMessage struct {
	Message string
	ChatID  int64
}
type Game struct {
	turn         int
	Pl1          GamePlayer
	Pl2          GamePlayer
	BulletsOrder []string
	Items        []string
	damage       int
}
