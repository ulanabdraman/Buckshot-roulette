package lobby

import (
	"Buckshot_Roulette/models"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func DeletePlayerFromLobby(bot *tgbotapi.BotAPI, lobby models.Lobby, chatID int64) {
	for i, lb := range Lobbies {
		if lb.Code == lobby.Code {
			playerIndex := -1
			for j, player := range lobby.Players {
				if player.ChatID == chatID {
					playerIndex = j
					break
				}
			}

			if playerIndex != -1 {
				message := fmt.Sprintf("Игрок с Username %s покинул лобби", lb.Players[0].Username)
				for _, pc := range Lobbies[i].Players {
					replyMessage := tgbotapi.NewMessage(pc.ChatID, message)
					bot.Send(replyMessage)
				}
				Lobbies[i].Players = append(Lobbies[i].Players[:playerIndex], Lobbies[i].Players[playerIndex+1:]...)
			}

			if len(Lobbies[i].Players) == 0 {
				Lobbies = append(Lobbies[:i], Lobbies[i+1:]...)
			}
			break
		}
	}
}
