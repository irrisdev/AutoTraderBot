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

			session.Stage = 1
			inlineKeyboard := telebot.NewInlineKeyboardMarkup(
				telebot.NewInlineKeyboardRow(
					telebot.NewInlineKeyboardButtonData("Start Tracking Average Price", "begin"),
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

		//Makes sure program doesn't error if user CallbackQuery is outside session scope
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

		if session.Stage > 2 {
			_, exists := AllCarMakes[data]
			if exists {
				fmt.Println(data)
			}
		}

		switch data {

		case "cancel":
			session.deleteSession()

		case "begin":
			editMessage := telebot.NewEditMessageTextAndMarkup(ChatID, MessageID, "Choose the vehicle manufacture", KeyboardMap[data])
			sendEditMessage(editMessage)
			session.Stage = 2

		case "popular":
			editMessage := telebot.NewEditMessageTextAndMarkup(ChatID, MessageID, "Choose the make of the vehicle", KeyboardMap[data])
			sendEditMessage(editMessage)
			session.Stage = 3

		case "nextPopular":
			editMessage := telebot.NewEditMessageTextAndMarkup(ChatID, MessageID, "Choose the make of the vehicle", KeyboardMap[data])
			sendEditMessage(editMessage)
			session.Stage = 3

		}

		return

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
