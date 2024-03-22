package bot

import (
	telebot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
)

func main() {

	bot, err := telebot.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))

	if err != nil {
		panic(err)
	}

	bot.Debug = true

}
