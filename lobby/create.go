package lobby

import (
	"Buckshot_Roulette/models"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
	"time"
)

var Lobbies []models.Lobby
var Slots [10]bool

func LobbyCreate(bot *tgbotapi.BotAPI, chatID int64, level int, userID int, username string) models.Lobby {
	var lobby models.Lobby
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
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
	Lobbies[findSlot()] = lobby
	return lobby
}
func findSlot() int {
	for i, slot := range Slots {
		if !slot {
			Slots[i] = true
			return i
		}
	}
	return -1
}
