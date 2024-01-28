package lobby

import (
	"Buckshot_Roulette/models"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func StartGame(bot *tgbotapi.BotAPI, lobby models.Lobby, chatID int64) (bool, models.Lobby, int) {
	for i, lb := range Lobbies {
		if lb.Code == lobby.Code {
			if lb.Players[0].ChatID == chatID {
				if len(lb.Players) < 2 {
					message := fmt.Sprintf("Недостаточно игроков")
					replyMessage := tgbotapi.NewMessage(chatID, message)
					bot.Send(replyMessage)
				} else {
					message := fmt.Sprintf("Все условия для начала игры выполнены")
					replyMessage := tgbotapi.NewMessage(chatID, message)
					bot.Send(replyMessage)
					return true, lb, i
				}

			} else {
				message := fmt.Sprintf("Вы не создатель лобби")
				replyMessage := tgbotapi.NewMessage(chatID, message)
				bot.Send(replyMessage)
			}

		}
	}
	var l models.Lobby
	return false, l, 0
}
func Find(lobby models.Lobby) int {
	for i, lb := range Lobbies {
		if lb.Code == lobby.Code {
			return i
		}
	}
	return -1
}
