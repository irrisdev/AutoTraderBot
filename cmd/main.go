package main

import (
	"github.com/irrisdev/AutoTraderBot/bot"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {

	//Setup logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	//Load env variables ~ Telegram Bot API Token
	err := godotenv.Load()
	if err != nil {
		log.Err(err).Msg("Failed to load environment variables")
		return
	}
	log.Info().Msg("Loaded environment variables")

	status := make(chan string)

	go bot.TeleBot(status)
	//bot.ScrapeModelTest()
	//bot.ScrapeMakes()
	log.Info().Msg(<-status)

}
