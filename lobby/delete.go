package lobby

import (
	"Buckshot_Roulette/game"
	"Buckshot_Roulette/models"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func DeletePlayerFromLobby(bot *tgbotapi.BotAPI, lobby models.Lobby, chatID int64, messageCh chan game.GameMessage) {
	for i, lb := range Lobbies {

		if lb.Code == lobby.Code {

			playerIndex := -1

			for j, player := range Lobbies[i].Players {
				if player.ChatID == chatID {
					playerIndex = j
					break
				}
			}

			if playerIndex != -1 {
				message := fmt.Sprintf("Игрок с Username %s покинул лобби", lb.Players[playerIndex].Username)
				for _, pc := range Lobbies[i].Players {
					replyMessage := tgbotapi.NewMessage(pc.ChatID, message)
					bot.Send(replyMessage)
				}
				Lobbies[i].Players = append(Lobbies[i].Players[:playerIndex], Lobbies[i].Players[playerIndex+1:]...)
			}

			if len(Lobbies[i].Players) == 0 {
				if Lobbies[i].Game {
					var message game.GameMessage
					message.Message = "/endgame"
					message.ChatID = 0
					messageCh <- message
				}
				var l models.Lobby
				Lobbies[i] = l
				Slots[i] = false
			}
			break
		}
	}
}
