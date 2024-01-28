package game

import (
	"Buckshot_Roulette/models"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
	"time"
)

func (g *Game) LevelSecond(bot *tgbotapi.BotAPI, lobby models.Lobby, messageCh chan GameMessage) {
	player1 := lobby.Players[0].Username
	player2 := lobby.Players[1].Username
	message1 := lobby.Players[0].ChatID
	message2 := lobby.Players[1].ChatID
	allmessage := lobby.Players
	g.Pl1.Hp = 4
	g.Pl2.Hp = 4
	rand.Seed(time.Now().UnixNano())
	g.turn = rand.Intn(2)
	g.Pl1.block = false
	g.Pl2.block = false

	message := fmt.Sprintf("Игра началась!")
	g.sendMessageToAll(bot, allmessage, message)

	for {
		g.damage = 1
		g.distributeItems()
		totalBullets := rand.Intn(6) + 3
		blankBullets := rand.Intn(totalBullets-1) + 1

		if blankBullets == totalBullets {
			blankBullets--
		}

		combatBullets := totalBullets - blankBullets

		message = fmt.Sprintf("Холостых: %d \nБоевых: %d", blankBullets, combatBullets)
		g.sendMessageToAll(bot, allmessage, message)

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

				if g.Pl2.block {
					g.Pl2.block = false
					msg := fmt.Sprintf("Ваш ход заблокирован")
					g.sendMessage(bot, message2, msg)
					msg = fmt.Sprintf("Ход вашего противника заблокирован")
					g.sendMessage(bot, message1, msg)
					continue
				}
				message = fmt.Sprintf("Очередь второго игрока '%s'", player2)
				g.sendMessageToAll(bot, allmessage, message)
				end := g.choiceItems(bot, message2, messageCh, message1, allmessage)
				if end {
					return
				}
				if len(g.BulletsOrder) == 0 {
					break
				}

			} else {

				if g.Pl1.block {
					g.Pl1.block = false
					msg := fmt.Sprintf("Ваш ход заблокирован")
					g.sendMessage(bot, message1, msg)
					msg = fmt.Sprintf("Ход вашего противника заблокирован")
					g.sendMessage(bot, message2, msg)
					continue
				}
				message = fmt.Sprintf("Очередь первого игрока '%s'", player1)
				g.sendMessageToAll(bot, allmessage, message)
				end := g.choiceItems(bot, message1, messageCh, message2, allmessage)
				if end {
					return
				}
				if len(g.BulletsOrder) == 0 {
					break
				}

			}
		}
	}
}
