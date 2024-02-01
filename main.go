package main

import (
	"Buckshot_Roulette/game"
	"Buckshot_Roulette/models"
	"fmt"
	"strconv"
	"strings"
	"time"

	//"time"
	//"database/sql"
	//"fmt"
	"Buckshot_Roulette/lobby"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	//"strconv"
	//"strings"
)

var levels = []int{1, 2, 3}

func removeKeyboard(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	msg.ParseMode = "Markdown"
	bot.Send(msg)
}

var UserStates = make(map[int]models.UserState)
var playerLobby = make(map[int]models.Lobby)
var messageCh = make([]chan game.GameMessage, 10)

func main() {
	bot, err := tgbotapi.NewBotAPI("6422826842:AAGz359zrP2w3N8KvmB9dhYdFPSMeDz5V7I")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	for i := 0; i < 10; i++ {
		messageCh[i] = make(chan game.GameMessage)
	}

	messageAFK := make(chan models.UserInfo)
	go AfkChecker(bot, messageAFK)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	logger := NewLogger()
	go logger.Run()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		userID := update.Message.From.ID
		username := update.Message.From.UserName
		//firstName := update.Message.From.FirstName
		//lastName := update.Message.From.LastName
		messageText := update.Message.Text
		currentState := UserStates[userID]
		if update.Message == nil {
			continue
		}
		if update.Message.From.ID == bot.Self.ID {
			message := fmt.Sprintf("Отправлено сообщение от бота: %s", messageText)
			logger.Log(message)
		}
		if update.Message.Text != "" {
			message := fmt.Sprintf("Получено сообщение от пользователя %d: %s", userID, messageText)
			logger.Log(message)
		}
		if userID == 5791061854 {
			message := fmt.Sprintf("Я люблю тебя")
			replyMessage := tgbotapi.NewMessage(chatID, message)
			bot.Send(replyMessage)
		}
		if currentState == models.Idle {
			if messageText == "/start" {
				message := fmt.Sprintf("Вы в меню %s: \n Вы можете создать лобби '/createlobby' \n Вы можете присоединиться к лобби '/joinlobby' \n Вы можете выйти с лобби '/leavelobby'", username)

				var replyButtons []tgbotapi.KeyboardButton
				var options = []string{"/createlobby", "/joinlobby"}
				for _, option := range options {
					replyButtons = append(replyButtons, tgbotapi.NewKeyboardButton(option))
				}
				replyKeyboard := tgbotapi.NewReplyKeyboard(replyButtons)

				msg := tgbotapi.NewMessage(chatID, message)
				msg.ReplyMarkup = replyKeyboard

				bot.Send(msg)
			}
			if messageText == "/createlobby" {
				var replyButtons []tgbotapi.KeyboardButton
				for _, option := range levels {
					replyButtons = append(replyButtons, tgbotapi.NewKeyboardButton(strconv.Itoa(option)))
				}
				replyKeyboard := tgbotapi.NewReplyKeyboard(replyButtons)

				msg := tgbotapi.NewMessage(chatID, "Выберите сложность:")
				msg.ReplyMarkup = replyKeyboard
				bot.Send(msg)

				UserStates[userID] = models.GivingLevel
			}
			if messageText == "/joinlobby" {
				msg := fmt.Sprintf("Отправьте код лобби \nЕсли хотите вернуться введите '/back'")
				replyMessage := tgbotapi.NewMessage(chatID, msg)
				keyboard := tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("/back"),
					),
				)

				replyMessage.ReplyMarkup = keyboard
				bot.Send(replyMessage)
				UserStates[userID] = models.GivingCode
			}
			if messageText == "/leavelobby" {
				message := fmt.Sprintf("Вы не в лобби")
				var replyButtons []tgbotapi.KeyboardButton
				var options = []string{"/createlobby", "/joinlobby"}
				for _, option := range options {
					replyButtons = append(replyButtons, tgbotapi.NewKeyboardButton(option))
				}
				replyKeyboard := tgbotapi.NewReplyKeyboard(replyButtons)

				msg := tgbotapi.NewMessage(chatID, message)
				msg.ReplyMarkup = replyKeyboard
				bot.Send(msg)
			}
		}
		if currentState == models.InLobby {
			if strings.HasPrefix(messageText, "/createlobby") {
				message := fmt.Sprintf("Вы уже в лобби")
				replyMessage := tgbotapi.NewMessage(chatID, message)
				bot.Send(replyMessage)
			}
			if strings.HasPrefix(messageText, "/joinlobby") {
				message := fmt.Sprintf("Вы уже в лобби")
				replyMessage := tgbotapi.NewMessage(chatID, message)
				bot.Send(replyMessage)
			}
			if messageText == "/leavelobby" {
				UserStates[userID] = models.Idle
				lobby.DeletePlayerFromLobby(bot, playerLobby[userID], chatID, messageCh[lobby.Find(playerLobby[userID])])
				message := fmt.Sprintf("Вы покинули лобби")
				replyMessage := tgbotapi.NewMessage(chatID, message)
				bot.Send(replyMessage)
			}
			if messageText == "/startgame" {
				a, lb, i := lobby.StartGame(bot, playerLobby[userID], chatID)
				//log.Println(i)
				if a {
					playerLobby[userID] = lb
					UserStates[userID] = models.InGame
					UserStates[lb.Players[1].UserID] = models.InGame
					var g game.Game
					if lb.Level == 1 {
						go g.LevelFirst(bot, playerLobby[userID], messageCh[i], &UserStates)
					}
					if lb.Level == 2 {
						go g.LevelSecond(bot, playerLobby[userID], messageCh[i], &UserStates)
					}
					if lb.Level == 3 {
						go g.LevelThird(bot, playerLobby[userID], messageCh[i], &UserStates)
					}
				}
			}
		}
		if currentState == models.GivingLevel {
			level, err := strconv.Atoi(messageText)
			if err != nil {
			}
			if level == 1 || level == 2 || level == 3 {
				lb := lobby.LobbyCreate(bot, chatID, level, userID, username)
				UserStates[userID] = models.InLobby
				message := fmt.Sprintf("Лобби создано\nКод: `%s`\nСложность: %d\nКод можно скопировать нажатием и нужно скинуть другу\nЧтобы начать игру введите '/startgame'", lb.Code, lb.Level)
				removeKeyboard(bot, chatID, message)
				playerLobby[userID] = lb
			} else {
				message := fmt.Sprintf("Выбрана неправильная сложность")
				replyMessage := tgbotapi.NewMessage(chatID, message)
				bot.Send(replyMessage)
			}

		}
		if currentState == models.GivingCode {
			if messageText == "/back" {
				UserStates[userID] = models.Idle
				message := fmt.Sprintf("Вы в меню %s: \n Вы можете создать лобби '/createlobby' \n Вы можете присоединиться к лобби '/joinlobby' \n Вы можете выйти с лобби '/leavelobby'", username)

				var replyButtons []tgbotapi.KeyboardButton
				var options = []string{"/createlobby", "/joinlobby"}
				for _, option := range options {
					replyButtons = append(replyButtons, tgbotapi.NewKeyboardButton(option))
				}
				replyKeyboard := tgbotapi.NewReplyKeyboard(replyButtons)

				msg := tgbotapi.NewMessage(chatID, message)
				msg.ReplyMarkup = replyKeyboard

				bot.Send(msg)
			} else {
				a, lb := lobby.JoinLobby(chatID, bot, messageText, userID, username)
				if a {
					UserStates[userID] = models.InLobby
					message := fmt.Sprintf("Успешное подключение к %s", messageText)
					replyMessage := tgbotapi.NewMessage(chatID, message)
					bot.Send(replyMessage)
					playerLobby[userID] = lb
				}
			}
		}
		if currentState == models.InGame {
			if messageText == "/leavelobby" {
				fmt.Println("Покидаем лобби...")
				lobby.DeletePlayerFromLobby(bot, playerLobby[userID], chatID, messageCh[lobby.Find(playerLobby[userID])])
				UserStates[userID] = models.Idle
				message := fmt.Sprintf("Вы покинули лобби")
				replyMessage := tgbotapi.NewMessage(chatID, message)
				bot.Send(replyMessage)
			} else {
				fmt.Println("Отправляем сообщение в игру...")
				var message game.GameMessage
				message.Message = messageText
				message.ChatID = chatID
				messageCh[lobby.Find(playerLobby[userID])] <- message
			}
		}

		if messageText == "/showlobbies" {
			var message string
			var lb models.Lobby
			message += fmt.Sprintf("Список лобби \n")
			for _, lb = range lobby.Lobbies {
				message += fmt.Sprintf("%d %s %d\n", lb.Level, lb.Code, len(lb.Players))
			}
			replyMessage := tgbotapi.NewMessage(chatID, message)
			bot.Send(replyMessage)
		}
		if messageText == "/help" {
			helpText := `
1. **Начало игры и создание лобби:**
   - Используйте команду /start, чтобы открыть меню и начать новую игру.
   - Выберите команду /createlobby, чтобы создать новое лобби.
   - Введите уровень сложности (1, 2 или 3), чтобы создать лобби с выбранным уровнем.
   - Если хотите присоединиться к существующему лобби, используйте команду /joinlobby.

2. **Правила игры:**
   - Цель игры - выжить и убить противника, можно использовать предметы.
   - Перед каждым раундом генерируется случайное количество патронов и конкретное (в зависимости от уровня) количество предметов.
   - Вам дается выбор стрельбы в себя или противника. В случае, если вы стреляете в себя, и патрон оказывается холостым, ваш ход сохраняется.
   - Игра заканчивается, когда кто-то из игроков умирает.

3. **Предметы:**
   - 🍺 Пиво - пропускает текущий патрон без выстрела, сохраняя вашу очередь.
   - 🚬 Сигарета - восстанавливает 1 здоровье (хп).
   - 🔪 Нож - следующий выстрел будет с двойным уроном.
   - 🔍 Лупа - позволяет узнать, какой патрон в дробовике на данный момент.
   - 🔗 Наручники - блокирует 1 ход противника.

4. **Уровни:**
   - **Уровень 1:** Простой уровень без дополнительных предметов.
   - **Уровень 2:** Игроки начинают с 4 хп, генерируются 2 предмета перед каждым раундом.
   - **Уровень 3:** Игроки начинают с 5 хп, генерируются 4 предмета перед каждым раундом.

5. **Команды и объяснение:**
   - /start - открыть меню и начать новую игру.
   - /createlobby - создать новое лобби.
   - /joinlobby - присоединиться к существующему лобби.
   - /leavelobby - покинуть текущее лобби.
   - /startgame - начать игру в текущем лобби.
   - /help - отобразить это сообщение справки.


    `
			replyMessage := tgbotapi.NewMessage(chatID, helpText)
			replyMessage.ParseMode = "Markdown"
			bot.Send(replyMessage)
		}

		if messageText == "/showmystate" {
			replyMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("Текущее состояние: %v", currentState))
			bot.Send(replyMessage)
		}
		var userMsg models.UserInfo
		userMsg.LastActivity = time.Now()
		userMsg.UserID = userID
		userMsg.ChatID = chatID
		messageAFK <- userMsg
	}
}
