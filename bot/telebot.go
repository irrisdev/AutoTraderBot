package bot

import (
	"fmt"
	telebot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"math/rand/v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var sessions = make(map[int64]*UserSession)
var bot *telebot.BotAPI
var err error

//var counter = 0

type UserSession struct {
	chatID         int64
	Messages       map[int]struct{}
	Stage          string
	RequestDetails CarRequest
	CarDetails     string
	Choosing       bool
}

func (s *UserSession) deleteSession() {
	var wg sync.WaitGroup
	wg.Add(len(s.Messages))

	func(msgs map[int]struct{}) {
		for msgId := range msgs {
			go deleteMessage(s.chatID, msgId, &wg)
		}
	}(s.Messages)

	wg.Wait()
	delete(sessions, s.chatID)
}

func deleteMessage(chatID int64, messageID int, group *sync.WaitGroup) {
	deleteMsg := telebot.NewDeleteMessage(chatID, messageID)
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
	err := telebot.SetLogger(&log.Logger)
	if err != nil {
		return
	}
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

func deleteMsg(chatID int64, messageID int) {
	deleteMsg := telebot.NewDeleteMessage(chatID, messageID)
	if _, err := bot.Request(deleteMsg); err != nil {
		log.Err(err).Msg("Failed to delete message")
	}
}

func handleUpdate(bot *telebot.BotAPI, update *telebot.Update) {
	if update.Message != nil {
		messageID := update.Message.MessageID
		UserID := update.Message.From.ID
		chatID := update.Message.Chat.ID

		session, exists := sessions[UserID]
		if !exists {
			session = &UserSession{
				chatID:   chatID,
				Messages: map[int]struct{}{},
				Stage:    "start",
				RequestDetails: CarRequest{
					Postcode: Postcodes[rand.IntN(len(Postcodes))],
				},
				CarDetails: "Car Request Details\n",
			}
			sessions[UserID] = session
		}

		session.Messages[messageID] = struct{}{}

		msg := telebot.NewMessage(update.Message.Chat.ID, "")

		if update.Message.Text == "Cancel" && session.Stage == "start" {
			session.deleteSession()
		}
		switch update.Message.Command() {
		case "start":

			session.Stage = "begin"
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
			send(session, msg)

			//case "screenshot":
			//	url := update.Message.CommandArguments()
			//	argsLen := strings.Split(update.Message.Text, " ")
			//	if len(argsLen) != 2 {
			//		msg := telebot.NewMessage(chatID, "Invalid website")
			//		msg.ReplyToMessageID = update.Message.MessageID
			//		_, err := bot.Send(msg)
			//		if err != nil {
			//			log.Err(err).Msg("Error while sending message")
			//		}
			//		return
			//	}
			//	if !strings.HasPrefix(url, "https://") {
			//		url = "https://" + url
			//	}
			//
			//	resp, err := http.Get(url)
			//	if err != nil {
			//		msg := telebot.NewMessage(chatID, "Invalid website")
			//		msg.ReplyToMessageID = update.Message.MessageID
			//		_, err := bot.Send(msg)
			//		if err != nil {
			//			log.Err(err).Msg("Error while sending message")
			//		}
			//		return
			//	}
			//	defer resp.Body.Close()
			//
			//	// Check the response status code
			//	if resp.StatusCode != http.StatusOK {
			//		msg := telebot.NewMessage(chatID, "Invalid website")
			//		msg.ReplyToMessageID = update.Message.MessageID
			//		_, err := bot.Send(msg)
			//		if err != nil {
			//			log.Err(err).Msg("Error while sending message")
			//		}
			//		return
			//	}
			//
			//	page := rod.New().MustConnect().MustPage(url).MustWaitLoad()
			//
			//	page.MustWindowMaximize()
			//	img, _ := page.Screenshot(false, &proto.PageCaptureScreenshot{
			//		Format:  proto.PageCaptureScreenshotFormatJpeg,
			//		Quality: gson.Int(90),
			//	})
			//
			//	_ = utils.OutputFile(fmt.Sprintf("screenshots/%d.jpg", counter), img)
			//	counter++
			//	msg := telebot.NewPhoto(chatID, telebot.FileBytes{Bytes: img})
			//	sentMsg, err := bot.Send(msg)
			//	if err != nil {
			//		log.Err(err).Msg("Error while sending image")
			//	}
			//	session.Messages[sentMsg.MessageID] = struct{}{}
			//
		}

	} else if update.CallbackQuery != nil {

		messageID := update.CallbackQuery.Message.MessageID
		chatID := update.CallbackQuery.Message.Chat.ID
		session, exists := sessions[chatID]

		//Makes sure program doesn't error if user CallbackQuery is outside session scope
		if !exists {
			newMsg := telebot.NewMessage(chatID, "Use /start to begin")
			sentMsg, err := bot.Send(newMsg)
			if err != nil {
				log.Err(err).Msg("Error while sending Message")
			}
			go func() {
				time.Sleep(time.Second * 3)
				deleteMsg(chatID, sentMsg.MessageID)
			}()
			deleteMsg(chatID, messageID)
			return
		}

		data := update.CallbackQuery.Data

		//If callback is a make
		if _, exists := AllCarMakes[data]; exists {
			//complete := make(chan bool)

			session.RequestDetails.Make = data

			//go ScrapeModels(session, complete)

			session.CarDetails += fmt.Sprintf("\nMake : %s\n\nChoose the vehicle model", session.RequestDetails.Make)

			//<-complete
			editMessage := telebot.NewEditMessageTextAndMarkup(chatID, messageID, session.CarDetails, KeyboardMap["model"])
			sendEditMessage(editMessage)
			return
		}

		switch data {

		case "cancel":
			session.deleteSession()

		case "begin":
			editMessage := telebot.NewEditMessageTextAndMarkup(chatID, messageID, "Choose the vehicle manufacture", KeyboardMap[data])
			sendEditMessage(editMessage)

		case "popular":
			editMessage := telebot.NewEditMessageTextAndMarkup(chatID, messageID, "Choose the make of the vehicle", KeyboardMap[data])
			sendEditMessage(editMessage)

		case "nextPopular":
			editMessage := telebot.NewEditMessageTextAndMarkup(chatID, messageID, "Choose the make of the vehicle", KeyboardMap[data])
			sendEditMessage(editMessage)
		}

		return

	}

}

//func getRequestString(details scraper.CarRequest, result string) string {
//
//	result = fmt.Sprintf("Car Request Details:\n"+
//		"Make: %s\n"+
//		"Model: %s\n"+
//		"Trim: %s\n"+
//		"Year From: %s\n"+
//		"Year To: %s\n",
//		details.Make, details.Model, details.AggregatedTrim, details.YearFrom, details.YearTo)
//
//	return result
//}

func send(session *UserSession, msg telebot.MessageConfig) {

	if len(msg.Text) == 0 {
		return
	}

	sentMsg, err := bot.Send(msg)
	if err != nil {
		log.Err(err).Msg("An error occurred while sending a message")
	}
	msgID := sentMsg.MessageID
	session.Messages[msgID] = struct{}{}

}

func sendEditMessage(editMsg telebot.EditMessageTextConfig) {
	_, err := bot.Send(editMsg)
	if err != nil {
		log.Err(err).Msg("An error occurred while editing a message")

	}
}
