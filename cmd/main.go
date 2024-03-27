package main

import (
	"github.com/irrisdev/AutoTraderBot/bot"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
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

	//status := make(chan string)

	//go bot.TeleBot(status)
	//bot.ScrapeModelTest()
	//bot.ScrapeMakes()
	makes := bot.ScrapeMakes()
	var wg sync.WaitGroup
	for i := 0; i < len(makes)-1; i += 20 {
		wg.Add(1)
		newSlice := makes[i : i+19]
		go bot.InsertAll(newSlice, &wg)

	}
	wg.Wait()
	//go bot.InsertAll(makes)
	//go bot.InsertAll(makes)
	//go bot.InsertAll(makes)
	log.Info().Msg("Fully Finishes")

}
