package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var carMakes = []string{
	"Audi",
	"BMW",
	"Ford",
	"Jaguar",
	"Land Rover",
	"Mercedes",
	"Nissan",
	"Porsche",
	"Toyota",
	"Vauxhall",
	"Volkswagen",
	"Volvo",
}

var keyboardMap = make(map[string]tgbotapi.InlineKeyboardMarkup)

func getInlineKeyboard(keyboard string) tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("", ""),
		),
	)
}
