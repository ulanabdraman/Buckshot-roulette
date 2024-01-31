package lobby

import (
	"Buckshot_Roulette/models"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
	"time"
)

var Lobbies = make([]models.Lobby, 10)
var Slots = make([]bool, 10)

func LobbyCreate(bot *tgbotapi.BotAPI, chatID int64, level int, userID int, username string) models.Lobby {
	var lobby models.Lobby
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ0123456789"
	codeLength := 6

	code := make([]byte, codeLength)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	lobby.Code = string(code)
	lobby.Level = level
	var player models.Player
	player.UserID = userID
	player.ChatID = chatID
	player.Username = username
	lobby.Players = append(lobby.Players, player)
	Lobbies[findSlots()] = lobby
	return lobby
}
func findSlots() int {
	for i, _ := range Slots {
		if Slots[i] == false {
			Slots[i] = true
			return i
		}
	}
	return -1
}
