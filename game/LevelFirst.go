package game

import (
	"Buckshot_Roulette/models"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"math/rand"
	"time"
)

func (g *Game) LevelFirst(bot *tgbotapi.BotAPI, lb models.Lobby, messageCh chan GameMessage, userStates *map[int]models.UserState) {
	player1 := lb.Players[0].Username
	player2 := lb.Players[1].Username
	message1 := lb.Players[0].ChatID
	message2 := lb.Players[1].ChatID
	allmessage := lb.Players
	g.Pl1.Hp = 2
	g.Pl2.Hp = 2
	rand.Seed(time.Now().UnixNano())
	g.turn = rand.Intn(2)
	message := fmt.Sprintf("Игра началась!")
	g.sendMessageToAll(bot, allmessage, message)

	for {
		totalBullets := rand.Intn(6) + 3

		blankBullets := rand.Intn(totalBullets-1) + 1

		if blankBullets == totalBullets {
			blankBullets--
		}

		combatBullets := totalBullets - blankBullets

		message = fmt.Sprintf("Холостых: %d \nБоевых: %d", blankBullets, combatBullets)
		replyMessage := tgbotapi.NewMessage(message1, message)
		sentMsg1, err := bot.Send(replyMessage)
		replyMessage = tgbotapi.NewMessage(message2, message)
		sentMsg2, err := bot.Send(replyMessage)
		time.Sleep(5 * time.Second)
		newMsg1 := tgbotapi.NewEditMessageText(message1, sentMsg1.MessageID, "Патроны больше не доступны")
		_, err = bot.Send(newMsg1)
		if err != nil {
			log.Panic(err)
		}
		newMsg2 := tgbotapi.NewEditMessageText(message2, sentMsg2.MessageID, "Патроны больше не доступны")
		_, err = bot.Send(newMsg2)
		if err != nil {
			log.Panic(err)
		}

		indices := rand.Perm(totalBullets)
		OrderBullet := make([]string, totalBullets)
		for i, index := range indices {
			if i < blankBullets {
				OrderBullet[index] = "Холостым"
			} else {
				OrderBullet[index] = "Боевым"
			}
		}
		g.BulletsOrder = OrderBullet

		for {
			g.turn++
			if g.turn%2 == 0 {
				if len(g.BulletsOrder) == 1 {
					message = fmt.Sprintf("Очередь второго игрока '%s'", player2)
					g.sendMessageToAll(bot, allmessage, message)
					end := g.choice(bot, message2, messageCh, message1, allmessage)
					if end {
						msg := "Вы можете ещё раз начать игру написав /startgame"
						g.sendMessageToAll(bot, allmessage, msg)
						if (*userStates)[lb.Players[0].UserID] == models.InGame {
							(*userStates)[lb.Players[0].UserID] = models.InLobby
							(*userStates)[lb.Players[1].UserID] = models.InLobby
						}
						return
					}
					break
				} else {
					message = fmt.Sprintf("Очередь второго игрока '%s'", player2)
					g.sendMessageToAll(bot, allmessage, message)
					end := g.choice(bot, message2, messageCh, message1, allmessage)
					if end {
						return
					}
				}

			} else {
				if len(g.BulletsOrder) == 1 {
					message = fmt.Sprintf("Очередь первого игрока '%s'", player1)
					g.sendMessageToAll(bot, allmessage, message)
					end := g.choice(bot, message1, messageCh, message2, allmessage)
					if end {

						if (*userStates)[lb.Players[0].UserID] == models.InGame {
							msg := "Вы можете ещё раз начать игру написав /startgame"
							g.sendMessageToAll(bot, allmessage, msg)
							(*userStates)[lb.Players[0].UserID] = models.InLobby
							(*userStates)[lb.Players[1].UserID] = models.InLobby
						}
						return
					}
					break
				} else {
					message = fmt.Sprintf("Очередь первого игрока '%s'", player1)
					g.sendMessageToAll(bot, allmessage, message)
					end := g.choice(bot, message1, messageCh, message2, allmessage)
					if end {
						return
					}
				}
			}
		}
	}
}
