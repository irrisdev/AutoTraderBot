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
	time.Sleep(time.Second * 5)
	done <- "Completed telegram"
}
