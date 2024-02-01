package lobby

import (
	"Buckshot_Roulette/models"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func JoinLobby(chatID int64, bot *tgbotapi.BotAPI, code string, userID int, username string) (bool, models.Lobby) {
	for i, lobby := range Lobbies {
		if lobby.Code == code {
			message := fmt.Sprintf("Лобби найден %s", code)
			replyMessage := tgbotapi.NewMessage(chatID, message)
			bot.Send(replyMessage)
			var player models.Player
			player.UserID = userID
			player.ChatID = chatID
			player.Username = username
			playerExists := false
			for _, play := range lobby.Players {
				if play.ChatID == chatID {
					playerExists = true
					break
				}
			}
			if len(Lobbies[i].Players) > 1 {
				message := fmt.Sprintf("Лобби полное")
				replyMessage := tgbotapi.NewMessage(chatID, message)
				bot.Send(replyMessage)
				var lb models.Lobby
				return false, lb
			}
			if playerExists {
				message := fmt.Sprintf("Игрок с UserID %d уже присутствует в лобби %s", userID, code)
				replyMessage := tgbotapi.NewMessage(chatID, message)
				bot.Send(replyMessage)
			} else {
				Lobbies[i].Players = append(Lobbies[i].Players, player)

				message := fmt.Sprintf("Игрок с Username %s добавлен в лобби %s", username, code)
				replyMessage := tgbotapi.NewMessage(chatID, message)
				bot.Send(replyMessage)
				replyMessage = tgbotapi.NewMessage(Lobbies[i].Players[0].ChatID, message)
				bot.Send(replyMessage)
			}
			return true, Lobbies[i]
		}
	}
	message := fmt.Sprintf("Лобби не найден %s", code)
	replyMessage := tgbotapi.NewMessage(chatID, message)
	bot.Send(replyMessage)
	var l models.Lobby
	return false, l
}
