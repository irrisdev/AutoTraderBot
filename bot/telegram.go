package bot

import (
	telebot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func StartTelegramBot(done chan<- string) {

	bot, err := telebot.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))

	if err != nil {

		log.Err(err).Msg("Failed to start telegram bot")
		return
	}

	log.Info().Msg("Started telegram bot")

	bot.Debug = true
	log.Info().Msg("Authorized on account " + bot.Self.UserName)

	updateConfig := telebot.NewUpdate(0)
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		//Discard empty messages as if no updates are detected in 30 seconds telegram will respond to request with empty struct
		if update.Message == nil {
			continue
		}
		if !update.Message.IsCommand() {
			continue
		}

		//fmt.Println(update.Message.Text)

		//msg := telebot.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID

		msg := telebot.NewMessage(update.Message.Chat.ID, "")

		switch update.Message.Command() {
		case "help":
			msg.Text = "I understand /start and /stop"
		case "start":
			msg.Text = "Make sure this is a channel"
		case "stop":
			msg.Text = "Stopping update in {channelname}"
		case "stopserver":
			done <- "Completed telegram"
		default:
			msg.Text = "Not a command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Err(err).Msg("An error occured while sending a message")
		}

	}

	time.Sleep(time.Second * 5)
}
