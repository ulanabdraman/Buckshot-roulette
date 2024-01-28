package main

import (
	"Buckshot_Roulette/game"
	"Buckshot_Roulette/models"
	"fmt"
	"strconv"
	"strings"
	//"database/sql"
	//"fmt"
	"Buckshot_Roulette/lobby"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	//"strconv"
	//"strings"
)

type UserState int

const (
	Idle UserState = iota
	InLobby
	GivingLevel
	GivingCode
	InGame
)

var levels = []int{1, 2, 3}

func removeKeyboard(bot *tgbotapi.BotAPI, chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	bot.Send(msg)
}

var userStates = make(map[int]UserState)

func main() {
	bot, err := tgbotapi.NewBotAPI("6422826842:AAGz359zrP2w3N8KvmB9dhYdFPSMeDz5V7I")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	playerLobby := make(map[int]models.Lobby)

	messageCh := make([]chan game.GameMessage, 10)
	for i := 0; i < 10; i++ {
		messageCh[i] = make(chan game.GameMessage)
	}
	//var lastActivity []time.Time

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
		currentState := userStates[userID]
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
		if currentState == Idle {
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
				userStates[userID] = GivingLevel
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
				userStates[userID] = GivingCode
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
		if currentState == InLobby {
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
				userStates[userID] = Idle
				lobby.DeletePlayerFromLobby(bot, playerLobby[userID], chatID)
				message := fmt.Sprintf("Вы покинули лобби")
				replyMessage := tgbotapi.NewMessage(chatID, message)
				bot.Send(replyMessage)
			}
			if messageText == "/startgame" {
				a, lb, i := lobby.StartGame(bot, playerLobby[userID], chatID)
				//log.Println(i)
				if a {
					playerLobby[userID] = lb
					userStates[userID] = InGame
					userStates[lb.Players[1].UserID] = InGame
					var g game.Game
					if lb.Level == 1 {
						go g.LevelFirst(bot, playerLobby[userID], messageCh[i])
					}
					if lb.Level == 2 {
						go g.LevelSecond(bot, playerLobby[userID], messageCh[i])
					}
					if lb.Level == 3 {

					}
				}
			}
		}
		if currentState == GivingLevel {
			level, err := strconv.Atoi(messageText)
			if err != nil {
			}
			if level == 1 || level == 2 || level == 3 {
				lb := lobby.LobbyCreate(bot, chatID, level, userID, username)
				userStates[userID] = InLobby
				message := fmt.Sprintf("Лобби создано \nКод: %s \nСложность:%d \nЧтобы начать игру введите '/startgame'", lb.Code, lb.Level)
				removeKeyboard(bot, chatID, message)
				playerLobby[userID] = lb
			} else {
				message := fmt.Sprintf("Выбрана неправильная сложность")
				replyMessage := tgbotapi.NewMessage(chatID, message)
				bot.Send(replyMessage)
			}

		}
		if currentState == GivingCode {
			if messageText == "/back" {
				userStates[userID] = Idle
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
					userStates[userID] = InLobby
					message := fmt.Sprintf("Успешное подключение к %s", messageText)
					replyMessage := tgbotapi.NewMessage(chatID, message)
					bot.Send(replyMessage)
					playerLobby[userID] = lb
				}
			}
		}
		if currentState == InGame {
			if messageText == "/leavelobby" {
				userStates[userID] = Idle
				lobby.DeletePlayerFromLobby(bot, playerLobby[userID], chatID)
				message := fmt.Sprintf("Вы покинули лобби")
				replyMessage := tgbotapi.NewMessage(chatID, message)
				bot.Send(replyMessage)

			} else {
				fmt.Println("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF\n")
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
		if messageText == "/showmystate" {
			replyMessage := tgbotapi.NewMessage(chatID, fmt.Sprintf("Текущее состояние: %v", currentState))
			bot.Send(replyMessage)
		}
	}
}
