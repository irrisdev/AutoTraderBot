package main

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// Initialising bots, set up scraping and schedules
func main() {

	err := loadEnvVariables()
	if err != nil {
		log.Err(err).Msg("Failed to load environment variables")
		return
	}

}

func loadEnvVariables() error {

	err := godotenv.Load()

	if err != nil {
		return err
	}

	return err
}
