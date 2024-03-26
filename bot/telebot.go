package bot

import (
	"fmt"
	telebot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var sessions = make(map[int64]*UserSession)
var bot *telebot.BotAPI
var err error

type UserSession struct {
	ChatID       int64
	Messages     map[int]struct{}
	UserMessages map[int]struct{}
	BotMessages  map[int]struct{}
	Stage        int
}

func (s *UserSession) deleteSession() {
	var wg sync.WaitGroup
	wg.Add(len(s.Messages))

	func(msgs map[int]struct{}) {
		for msgId := range msgs {
			go deleteMessage(s.ChatID, msgId, &wg)
		}
	}(s.Messages)

	wg.Wait()
	delete(sessions, s.ChatID)
}

func deleteMessage(ChatID int64, MessageId int, group *sync.WaitGroup) {
	deleteMsg := telebot.NewDeleteMessage(ChatID, MessageId)
	_, err := bot.Request(deleteMsg)
	_ = err
	group.Done()
}

func TeleBot(status chan<- string) {

	// Create a channel to receive OS signals
	sig := make(chan os.Signal, 1)

	// Register for interrupt signals (e.g., Ctrl+C)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// Block until a signal is received
	go waitForExit(sig)

	bot, err = telebot.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Err(err).Msg("Failed to start telegram bot")
		status <- "Failed"
		return
	}

	log.Info().Msg("Authorised on account " + bot.Self.UserName)

	bot.Debug = true
	updateConfig := telebot.NewUpdate(0)
	updateConfig.Timeout = 60

	//Subscribe to updates
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		go handleUpdate(bot, &update)

	}

}

func waitForExit(sig chan os.Signal) {

	<-sig

	log.Debug().Msg("Terminating Telegram Bot")
	log.Debug().Msg("Purging all UserSessions")

	for _, usrSession := range sessions {
		usrSession.deleteSession()
	}

	os.Exit(0)

}

func deleteMsg(chatId int64, messageId int) {
	deleteMsg := telebot.NewDeleteMessage(chatId, messageId)
	if _, err := bot.Request(deleteMsg); err != nil {
		log.Err(err).Msg("Failed to delete message")
	}
}

func handleUpdate(bot *telebot.BotAPI, update *telebot.Update) {
	if update.Message != nil {
		MessageID := update.Message.MessageID
		UserID := update.Message.From.ID
		ChatID := update.Message.Chat.ID

		session, exists := sessions[UserID]
		if !exists {
			session = &UserSession{
				ChatID:       ChatID,
				Messages:     map[int]struct{}{},
				UserMessages: map[int]struct{}{},
				BotMessages:  map[int]struct{}{},
				Stage:        0,
			}
			sessions[UserID] = session
		}

		session.UserMessages[MessageID] = struct{}{}
		session.Messages[MessageID] = struct{}{}

		msg := telebot.NewMessage(update.Message.Chat.ID, "")

		if update.Message.Text == "Cancel" && session.Stage > 0 {
			session.deleteSession()
		}

		switch update.Message.Command() {
		case "start":
			if session.Stage > 0 {
				deleteMsg(ChatID, MessageID)
				return
			}
			session.Stage = 1
			inlineKeyboard := telebot.NewInlineKeyboardMarkup(
				telebot.NewInlineKeyboardRow(
					telebot.NewInlineKeyboardButtonData("Start Tracking Average Price", "beginBetter"),
				),
				telebot.NewInlineKeyboardRow(
					telebot.NewInlineKeyboardButtonData("Cancel", "cancel"),
				),
			)
			msg.Text = fmt.Sprintf("Welcome, @%s! Please select an option", update.Message.Chat.UserName)
			msg.ReplyMarkup = inlineKeyboard
		}

		send(session, msg)

	} else if update.CallbackQuery != nil {
		MessageID := update.CallbackQuery.Message.MessageID
		ChatID := update.CallbackQuery.Message.Chat.ID

		session, exists := sessions[ChatID]
		if !exists {
			newMsg := telebot.NewMessage(ChatID, "Use /start to begin")
			sentMsg, err := bot.Send(newMsg)
			if err != nil {
				log.Err(err).Msg("Error while sending Message")
			}
			go func() {
				time.Sleep(time.Second * 3)
				deleteMsg(ChatID, sentMsg.MessageID)
			}()
			deleteMsg(ChatID, MessageID)
			return
		}

		data := update.CallbackQuery.Data

		switch data {

		case "cancel":
			session.deleteSession()

		case "begin":
			newReplyKeyboard := telebot.NewReplyKeyboard(
				telebot.NewKeyboardButtonRow(
					telebot.NewKeyboardButton("Audi"),
					telebot.NewKeyboardButton("BMW"),
					telebot.NewKeyboardButton("Ford"),
				),
				telebot.NewKeyboardButtonRow(
					telebot.NewKeyboardButton("Jaguar"),
					telebot.NewKeyboardButton("Land Rover"),
					telebot.NewKeyboardButton("Mercedes"),
				),
				telebot.NewKeyboardButtonRow(
					telebot.NewKeyboardButton("Nissan"),
					telebot.NewKeyboardButton("Porsche"),
					telebot.NewKeyboardButton("Toyota"),
				),
				telebot.NewKeyboardButtonRow(
					telebot.NewKeyboardButton("Vauxhall"),
					telebot.NewKeyboardButton("Volkswagen"),
					telebot.NewKeyboardButton("Volvo"),
				),
				telebot.NewKeyboardButtonRow(
					telebot.NewKeyboardButton("Other"),
				),
				telebot.NewKeyboardButtonRow(
					telebot.NewKeyboardButton("Back"),
					telebot.NewKeyboardButton("Cancel"),
					telebot.NewKeyboardButton("Next"),
				),
			)
			newMsg := telebot.NewMessage(session.ChatID, "ok")
			newMsg.ReplyMarkup = newReplyKeyboard

			send(session, newMsg)

		case "beginBetter":

			topLevel := telebot.NewInlineKeyboardMarkup(
				telebot.NewInlineKeyboardRow(
					telebot.NewInlineKeyboardButtonData("Popular Makes", "popular"),
				),
				telebot.NewInlineKeyboardRow(
					telebot.NewInlineKeyboardButtonData("German", "german"),
					telebot.NewInlineKeyboardButtonData("British", "british"),
					telebot.NewInlineKeyboardButtonData("Japanese", "japanese"),
				),
				telebot.NewInlineKeyboardRow(
					telebot.NewInlineKeyboardButtonData("Other", "other"),
					telebot.NewInlineKeyboardButtonData("Cancel", "cancel"),
				),
			)

			editMessage := telebot.NewEditMessageTextAndMarkup(ChatID, MessageID, "Choose the vehicle manufacture", topLevel)

			sendEditMessage(editMessage)

			//germanCarsRow := telebot.NewInlineKeyboardRow(
			//	telebot.NewInlineKeyboardButtonData("Audi", "Audi"),
			//	telebot.NewInlineKeyboardButtonData("BMW", "BMW"),
			//	telebot.NewInlineKeyboardButtonData("Mercedes", "Mercedes"),
			//	telebot.NewInlineKeyboardButtonData("Porsche", "Porsche"),
			//)
			//britishCarsRow := telebot.NewInlineKeyboardRow(
			//	telebot.NewInlineKeyboardButtonData("Jaguar", "Jaguar"),
			//	telebot.NewInlineKeyboardButtonData("Land Rover", "LandRover"),
			//	telebot.NewInlineKeyboardButtonData("Vauxhall", "Vauxhall"),
			//	telebot.NewInlineKeyboardButtonData("Volkswagen", "Volkswagen"),
			//)
			//japaneseCarsRow := telebot.NewInlineKeyboardRow(
			//	telebot.NewInlineKeyboardButtonData("Nissan", "Nissan"),
			//	telebot.NewInlineKeyboardButtonData("Toyota", "Toyota"),
			//	telebot.NewInlineKeyboardButtonData("Volvo", "Volvo"),
			//)
		case "popular":
			popularKeyboard := telebot.NewInlineKeyboardMarkup(
				telebot.NewInlineKeyboardRow(
					telebot.NewInlineKeyboardButtonData("Audi", "Audi"),
					telebot.NewInlineKeyboardButtonData("BMW", "BMW"),
				),
				telebot.NewInlineKeyboardRow(
					telebot.NewInlineKeyboardButtonData("Ford", "Ford"),
					telebot.NewInlineKeyboardButtonData("Jaguar", "Jaguar"),
				),
				telebot.NewInlineKeyboardRow(
					telebot.NewInlineKeyboardButtonData("Land Rover", "Land_Rover"),
					telebot.NewInlineKeyboardButtonData("Mercedes", "Mercedes"),
				),
				telebot.NewInlineKeyboardRow(
					telebot.NewInlineKeyboardButtonData("Back", "beginBetter"),
					telebot.NewInlineKeyboardButtonData("Cancel", "cancel"),
					telebot.NewInlineKeyboardButtonData("Next", "nextPopular"),
				),
			)
			editMessage := telebot.NewEditMessageTextAndMarkup(ChatID, MessageID, "Choose the make of the vehicle", popularKeyboard)
			sendEditMessage(editMessage)
		case "nextPopular":
			nextPopular := telebot.NewInlineKeyboardMarkup(
				telebot.NewInlineKeyboardRow(
					telebot.NewInlineKeyboardButtonData("Nissan", "Nissan"),
					telebot.NewInlineKeyboardButtonData("Porsche", "Porsche"),
				),
				telebot.NewInlineKeyboardRow(
					telebot.NewInlineKeyboardButtonData("Toyota", "Toyota"),
					telebot.NewInlineKeyboardButtonData("Vauxhall", "Vauxhall"),
				),
				telebot.NewInlineKeyboardRow(
					telebot.NewInlineKeyboardButtonData("Volkswagen", "Volkswagen"),
					telebot.NewInlineKeyboardButtonData("Volvo", "Volvo"),
				),
				telebot.NewInlineKeyboardRow(
					telebot.NewInlineKeyboardButtonData("Back", "popular"),
					telebot.NewInlineKeyboardButtonData("Cancel", "cancel"),
				),
			)
			editMessage := telebot.NewEditMessageTextAndMarkup(ChatID, MessageID, "Choose the make of the vehicle", nextPopular)
			sendEditMessage(editMessage)

		}
		_, err := bot.Send(telebot.NewCallback(update.CallbackQuery.ID, ""))
		_ = err

	}

}

func send(session *UserSession, msg telebot.MessageConfig) {

	if len(msg.Text) == 0 {
		return
	}

	sentMsg, err := bot.Send(msg)
	if err != nil {
		log.Err(err).Msg("An error occurred while sending a message")
	}
	msgID := sentMsg.MessageID
	session.BotMessages[msgID] = struct{}{}
	session.Messages[msgID] = struct{}{}

}

func sendEditMessage(editMsg telebot.EditMessageTextConfig) {
	_, err := bot.Send(editMsg)
	if err != nil {
		log.Err(err).Msg("An error occurred while editing a message")

	}
}
