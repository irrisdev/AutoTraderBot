package main

import (
	"github.com/irrisdev/AutoTraderBot/bot"
	"github.com/irrisdev/AutoTraderBot/scraper"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

// Initialising bots, set up scraping and schedules
func main() {
	//Setup logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	err := loadEnvVariables()
	if err != nil {
		log.Err(err).Msg("Failed to load environment variables")
		return
	}
	log.Info().Msg("Loaded environment variables")

	telegramStatus := make(chan string)

	go bot.StartTelegramBot(telegramStatus)

	go scraper.Scrape()

	log.Info().Msg(<-telegramStatus)

}

func loadEnvVariables() error {

	err := godotenv.Load()

	if err != nil {
		return err
	}

	return err
}
