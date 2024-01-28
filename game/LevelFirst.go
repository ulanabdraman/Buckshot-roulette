package game

import (
	"Buckshot_Roulette/models"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
	"time"
)

func (g *Game) LevelFirst(bot *tgbotapi.BotAPI, lobby models.Lobby, messageCh chan GameMessage) {
	player1 := lobby.Players[0].Username
	player2 := lobby.Players[1].Username
	message1 := lobby.Players[0].ChatID
	message2 := lobby.Players[1].ChatID
	allmessage := lobby.Players
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
				if len(g.BulletsOrder) == 1 {
					message = fmt.Sprintf("Очередь второго игрока '%s'", player2)
					g.sendMessageToAll(bot, allmessage, message)
					end := g.choice(bot, message2, messageCh, message1, allmessage)
					if end {
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
