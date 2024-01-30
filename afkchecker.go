package main

import (
	"Buckshot_Roulette/lobby"
	"Buckshot_Roulette/models"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
)

var usersActivity []models.UserInfo

func AfkChecker(bot *tgbotapi.BotAPI, message chan models.UserInfo) {
	go every1minute(bot)
	for {
		select {
		case chat := <-message:
			for i, user := range usersActivity {
				if user.UserID == chat.UserID {
					usersActivity[i].LastActivity = chat.LastActivity
				}
			}
			usersActivity = append(usersActivity, chat)
		default:
			time.Sleep(1 * time.Second)
		}
	}

}
func every1minute(bot *tgbotapi.BotAPI) {
	for {
		time.Sleep(1 * time.Minute)
		for i := len(usersActivity) - 1; i >= 0; i-- {
			user := usersActivity[i]
			if time.Since(user.LastActivity) > 1*time.Minute {
				UserStates[user.UserID] = models.Idle
				usersActivity = append(usersActivity[:i], usersActivity[i+1:]...)
				index := lobby.Find(playerLobby[user.UserID])
				if index != -1 {
					fmt.Println("Удалили игрока")
					lobby.DeletePlayerFromLobby(bot, playerLobby[user.UserID], user.ChatID, messageCh[lobby.Find(playerLobby[user.UserID])])
				}
			}
		}
	}
}
