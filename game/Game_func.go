package game

import (
	"Buckshot_Roulette/models"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
	"strings"
	"time"
)

func (g *Game) sendMessageToAll(bot *tgbotapi.BotAPI, allmessage []models.Player, message string) {
	for _, pc := range allmessage {
		replyMessage := tgbotapi.NewMessage(pc.ChatID, message)
		bot.Send(replyMessage)
	}
}

func (g *Game) sendMessage(bot *tgbotapi.BotAPI, player int64, message string) {
	replyMessage := tgbotapi.NewMessage(player, message)
	bot.Send(replyMessage)
}

func (g *Game) showTable(bot *tgbotapi.BotAPI, chatID int64) {
	var table string
	table += fmt.Sprintf("Хп 1 игрока: %d\n", g.Pl1.Hp)
	table += fmt.Sprintf("Хп 2 игрока: %d\n", g.Pl2.Hp)
	table += fmt.Sprintf("На столе у 1 игрока: ")
	message := strings.Join(g.Pl1.Items, ", ")
	message += "\n"
	table += message
	table += fmt.Sprintf("На столе у 2 игрока: ")
	message = strings.Join(g.Pl2.Items, ", ")
	message += "\n"
	table += message
	g.sendMessage(bot, chatID, table)
}

func (g *Game) distributeItems() {
	availableItems := []string{
		"Пиво",
		"Сигареты",
		"Нож",
		"Лупа",
		"Наручник",
	}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 2; i++ {
		index1 := rand.Intn(len(availableItems))
		index2 := rand.Intn(len(availableItems))

		if len(g.Pl1.Items) < 8 {
			g.Pl1.Items = append(g.Pl1.Items, availableItems[index1])
		}
		if len(g.Pl2.Items) < 8 {
			g.Pl2.Items = append(g.Pl2.Items, availableItems[index2])
		}
	}
}

func (g *Game) choice(bot *tgbotapi.BotAPI, chatID int64, messageCh chan GameMessage, enemy int64, allmessage []models.Player) bool {
	message := "Выберите в кого стрелять:"
	replyMessage := tgbotapi.NewMessage(chatID, message)

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("В себя"),
			tgbotapi.NewKeyboardButton("В противника"),
			tgbotapi.NewKeyboardButton("Показать стол"),
		),
	)

	replyMessage.ReplyMarkup = keyboard

	bot.Send(replyMessage)

	for {
		select {
		case chat := <-messageCh:
			if chat.Message == "/endgame" {
				fmt.Println("Игра закончена")

				return true
			}
			if chat.ChatID == chatID {
				fmt.Printf("Получено сообщение из канала: %s\n", chat.Message)

				if chat.Message == "В себя" {
					message := fmt.Sprintf("Ваш противник направил дробовик в себя")
					g.sendMessage(bot, enemy, message)
					if g.BulletsOrder[0] == "Холостым" {
						message := fmt.Sprintf("Выстрел оказался %s", g.BulletsOrder[0])
						g.sendMessageToAll(bot, allmessage, message)
						message = fmt.Sprintf("Вам повезло, следующий ход ваш")
						g.sendMessage(bot, chatID, message)
						message = fmt.Sprintf("Вашему противнику повезло,вы пропускаете свой ход")
						g.sendMessage(bot, enemy, message)
						g.turn++
					} else {
						message := fmt.Sprintf("Выстрел оказался %s", g.BulletsOrder[0])
						g.sendMessageToAll(bot, allmessage, message)
						message = fmt.Sprintf("Вам не повезло, -1хп у вас")
						g.sendMessage(bot, chatID, message)
						message = fmt.Sprintf("Вам повезло, -1хп у врага")
						g.sendMessage(bot, enemy, message)
						if g.turn%2 == 0 {
							g.Pl2.Hp--
							if g.Pl2.Hp <= 0 {
								message = fmt.Sprintf("Второй игрок проиграл")
								g.sendMessageToAll(bot, allmessage, message)
								return true
							}
						}
						if g.turn%2 == 1 {
							g.Pl1.Hp--
							if g.Pl1.Hp <= 0 {
								message = fmt.Sprintf("Первый игрок проиграл")
								g.sendMessageToAll(bot, allmessage, message)
								return true
							}
						}
					}
					g.drop()
				}
				if chat.Message == "В противника" {
					message := fmt.Sprintf("Ваш противник направил дробовик в вас")
					g.sendMessage(bot, enemy, message)
					if g.BulletsOrder[0] == "Холостым" {
						message := fmt.Sprintf("Выстрел оказался %s", g.BulletsOrder[0])
						g.sendMessageToAll(bot, allmessage, message)
						message = fmt.Sprintf("Вам не повезло")
						g.sendMessage(bot, chatID, message)
						message = fmt.Sprintf("Вам повезло")
						g.sendMessage(bot, enemy, message)
					} else {
						message := fmt.Sprintf("Выстрел оказался %s", g.BulletsOrder[0])
						g.sendMessageToAll(bot, allmessage, message)
						message = fmt.Sprintf("Вам повезло -1хп у врага")
						g.sendMessageToAll(bot, allmessage, message)
						message = fmt.Sprintf("Вам не повезло, -1хп у вас")
						g.sendMessage(bot, enemy, message)
						if g.turn%2 == 0 {
							g.Pl1.Hp--
							if g.Pl1.Hp <= 0 {
								message = fmt.Sprintf("Второй игрок проиграл")
								g.sendMessageToAll(bot, allmessage, message)
								return true
							}
						}
						if g.turn%2 == 1 {
							g.Pl2.Hp--
							if g.Pl2.Hp <= 0 {
								message = fmt.Sprintf("Первый игрок проиграл")
								g.sendMessageToAll(bot, allmessage, message)
								return true
							}
						}
					}
					g.drop()
				}
				if chat.Message == "Показать стол" {
					g.showTable(bot, chatID)
					continue
				}
				return false
			}

		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func (g *Game) drop() {
	if len(g.BulletsOrder) > 0 {
		g.BulletsOrder = g.BulletsOrder[1:]
		fmt.Println("Первый патрон удален")
	} else {
		clear := make([]string, 0)
		g.BulletsOrder = clear
		fmt.Println("Срез пуст, удаление невозможно")
	}
}

func (g *Game) choiceItems(bot *tgbotapi.BotAPI, chatID int64, messageCh chan GameMessage, enemy int64, allmessage []models.Player) bool {
	message := "Выберите в кого стрелять:"
	replyMessage := tgbotapi.NewMessage(chatID, message)

	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("В себя"),
			tgbotapi.NewKeyboardButton("В противника"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Выбрать предмет"),
			tgbotapi.NewKeyboardButton("Показать стол"),
		),
	)

	replyMessage.ReplyMarkup = keyboard
	bot.Send(replyMessage)
	for {
		select {
		case chat := <-messageCh:
			if chat.Message == "/endgame" {
				fmt.Println("Игра закончена")

				return true
			}
			if chat.ChatID == chatID {
				fmt.Printf("Получено сообщение из канала: %s\n", chat.Message)

				if chat.Message == "В себя" {
					message := fmt.Sprintf("Ваш противник направил дробовик в себя")
					g.sendMessage(bot, enemy, message)
					if g.BulletsOrder[0] == "Холостым" {
						message := fmt.Sprintf("Выстрел оказался %s", g.BulletsOrder[0])
						g.sendMessageToAll(bot, allmessage, message)
						message = fmt.Sprintf("Вам повезло, следующий ход ваш")
						g.sendMessage(bot, chatID, message)
						message = fmt.Sprintf("Вашему противнику повезло,вы пропускаете свой ход")
						g.sendMessage(bot, enemy, message)
						g.turn++
					} else {
						message := fmt.Sprintf("Выстрел оказался %s", g.BulletsOrder[0])
						g.sendMessageToAll(bot, allmessage, message)
						message = fmt.Sprintf("Вам не повезло, -%dхп у вас", g.damage)
						g.sendMessage(bot, chatID, message)
						message = fmt.Sprintf("Вам повезло, -%dхп у врага", g.damage)
						g.sendMessage(bot, enemy, message)
						if g.turn%2 == 0 {
							g.Pl2.Hp -= g.damage
							if g.Pl2.Hp < 1 {
								message = fmt.Sprintf("Второй игрок проиграл")
								g.sendMessageToAll(bot, allmessage, message)
								return true
							}
							message = fmt.Sprintf("Хп 2 игрока: %d", g.Pl2.Hp)
							g.sendMessageToAll(bot, allmessage, message)
						}
						if g.turn%2 == 1 {
							g.Pl1.Hp -= g.damage
							if g.Pl1.Hp < 1 {
								message = fmt.Sprintf("Первый игрок проиграл")
								g.sendMessageToAll(bot, allmessage, message)
								return true
							}
							message = fmt.Sprintf("Хп 1 игрока: %d", g.Pl1.Hp)
							g.sendMessageToAll(bot, allmessage, message)
						}
					}
					g.drop()
					g.damage = 1
				}
				if chat.Message == "В противника" {
					message := fmt.Sprintf("Ваш противник направил дробовик в вас")
					g.sendMessage(bot, enemy, message)
					if g.BulletsOrder[0] == "Холостым" {
						message := fmt.Sprintf("Выстрел оказался %s", g.BulletsOrder[0])
						g.sendMessageToAll(bot, allmessage, message)
						message = fmt.Sprintf("Вам не повезло")
						g.sendMessage(bot, chatID, message)
						message = fmt.Sprintf("Вам повезло")
						g.sendMessage(bot, enemy, message)
					} else {
						message := fmt.Sprintf("Выстрел оказался %s", g.BulletsOrder[0])
						g.sendMessageToAll(bot, allmessage, message)
						message = fmt.Sprintf("Вам повезло -%dхп у врага", g.damage)
						g.sendMessage(bot, chatID, message)
						message = fmt.Sprintf("Вам не повезло, -%dхп у вас", g.damage)
						g.sendMessage(bot, enemy, message)
						if g.turn%2 == 0 {
							g.Pl1.Hp -= g.damage
							if g.Pl1.Hp < 1 {
								message = fmt.Sprintf("Первый игрок проиграл")
								g.sendMessageToAll(bot, allmessage, message)
								return true
							}
							message = fmt.Sprintf("Хп 1 игрока: %d", g.Pl1.Hp)
							g.sendMessageToAll(bot, allmessage, message)
						}
						if g.turn%2 == 1 {
							g.Pl2.Hp -= g.damage
							if g.Pl2.Hp < 1 {
								message = fmt.Sprintf("Второй игрок проиграл")
								g.sendMessageToAll(bot, allmessage, message)
								return true
							}
							message = fmt.Sprintf("Хп 2 игрока: %d", g.Pl2.Hp)
							g.sendMessageToAll(bot, allmessage, message)
						}
					}
					g.drop()
					g.damage = 1
				}
				if chat.Message == "Выбрать предмет" {
					g.useItem(bot, messageCh, chatID, allmessage)
					message := "Выберите в кого стрелять:"
					replyMessage := tgbotapi.NewMessage(chatID, message)

					keyboard := tgbotapi.NewReplyKeyboard(
						tgbotapi.NewKeyboardButtonRow(
							tgbotapi.NewKeyboardButton("В себя"),
							tgbotapi.NewKeyboardButton("В противника"),
						),
						tgbotapi.NewKeyboardButtonRow(
							tgbotapi.NewKeyboardButton("Выбрать предмет"),
							tgbotapi.NewKeyboardButton("Показать стол"),
						),
					)

					replyMessage.ReplyMarkup = keyboard
					bot.Send(replyMessage)
					continue
				}
				if chat.Message == "Показать стол" {
					g.showTable(bot, chatID)
					continue
				}

				return false

			}

		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func (g *Game) useItem(bot *tgbotapi.BotAPI, messageCh chan GameMessage, chatID int64, allmessage []models.Player) {
	message := "Ваш инвентарь:\n"
	if g.turn%2 == 0 {
		if len(g.Pl2.Items) == 0 {
			message += "Инвентарь пуст."
		} else {
			for i, item := range g.Pl2.Items {
				message += fmt.Sprintf("%d. %s\n", i+1, item)
			}
		}
	} else {
		if len(g.Pl1.Items) == 0 {
			message += "Инвентарь пуст."
		} else {
			for i, item := range g.Pl1.Items {
				message += fmt.Sprintf("%d. %s\n", i+1, item)
			}
		}
	}
	replyMessage := tgbotapi.NewMessage(chatID, message)
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Пиво"),
			tgbotapi.NewKeyboardButton("Сигареты"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Нож"),
			tgbotapi.NewKeyboardButton("Лупа"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Наручник"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Выход"),
		),
	)

	replyMessage.ReplyMarkup = keyboard
	bot.Send(replyMessage)
	for {
		select {
		case chat := <-messageCh:
			if chat.ChatID == chatID {
				g.itemLogics(bot, chat.Message, chatID, allmessage)
				return
			}

		default:
			time.Sleep(1 * time.Second)
		}
	}
}

func (g *Game) itemLogics(bot *tgbotapi.BotAPI, message string, chatID int64, allmessage []models.Player) {
	var currentPlayer GamePlayer
	if g.turn%2 == 0 {
		currentPlayer = g.Pl2
	} else {
		currentPlayer = g.Pl1
	}
	for i, item := range currentPlayer.Items {
		if item == message {

			switch message {
			case "Сигареты":
				if g.turn%2 == 0 && g.Pl2.Hp < 4 {
					g.Pl2.Hp++
					str := fmt.Sprintf("2 игрок восстановил себе 1хп")
					g.sendMessageToAll(bot, allmessage, str)
				}
				if g.turn%2 == 1 && g.Pl1.Hp < 4 {
					g.Pl1.Hp++
					str := fmt.Sprintf("1 игрок восстановил себе 1хп")
					g.sendMessageToAll(bot, allmessage, str)
				}
			case "Лупа":
				str := fmt.Sprintf("%d игрок использвал лупу", 2-g.turn%2)
				g.sendMessageToAll(bot, allmessage, str)
				str = fmt.Sprintf("Следующий патрон стрельнет: %s", g.BulletsOrder[0])
				g.sendMessage(bot, chatID, str)
			case "Пиво":
				str := fmt.Sprintf("%d игрок выпил пивка\n", 2-g.turn%2)
				g.sendMessageToAll(bot, allmessage, str)
				str += fmt.Sprintf("Патрон оказался: %s", g.BulletsOrder[0])
				g.sendMessageToAll(bot, allmessage, str)
				g.drop()
			case "Нож":
				g.damage = 2
				str := fmt.Sprintf("%d игрок использвал нож", 2-g.turn%2)
				g.sendMessageToAll(bot, allmessage, str)

			case "Наручник":
				if g.turn%2 == 0 {
					if g.Pl1.block == false {
						g.Pl1.block = true
						str := fmt.Sprintf("1 игрок был прикован в наручники")
						g.sendMessageToAll(bot, allmessage, str)
					} else {
						str := fmt.Sprintf("Вы уже использвали наручники")
						g.sendMessage(bot, chatID, str)
					}
				}
				if g.turn%2 == 1 {
					if g.Pl2.block == false {
						g.Pl2.block = true
						str := fmt.Sprintf("2 игрок был прикован в наручники")
						g.sendMessageToAll(bot, allmessage, str)
					} else {
						str := fmt.Sprintf("Вы уже использвали наручники")
						g.sendMessage(bot, chatID, str)
					}
				}
			case "Выход":
				{
					str := fmt.Sprintf("Вы покинули инвентарь")
					g.sendMessage(bot, chatID, str)
				}

			}

			if g.turn%2 == 0 {
				g.Pl2.Items = append(g.Pl2.Items[:i], g.Pl2.Items[i+1:]...)
				return
			}
			if g.turn%2 == 1 {
				g.Pl1.Items = append(g.Pl1.Items[:i], g.Pl1.Items[i+1:]...)
				return
			}

		}
	}
	str := fmt.Sprintf("У вас нету %s", message)
	g.sendMessage(bot, chatID, str)
}
