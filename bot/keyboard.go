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

var keyboardMap = map[string]tgbotapi.InlineKeyboardMarkup{
	"popular": tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Audi", "Audi"),
			tgbotapi.NewInlineKeyboardButtonData("BMW", "BMW"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ford", "Ford"),
			tgbotapi.NewInlineKeyboardButtonData("Jaguar", "Jaguar"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Land Rover", "Land_Rover"),
			tgbotapi.NewInlineKeyboardButtonData("Mercedes", "Mercedes"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", "begin"),
			tgbotapi.NewInlineKeyboardButtonData("Cancel", "cancel"),
			tgbotapi.NewInlineKeyboardButtonData("Next", "nextPopular"),
		),
	),
	"nextPopular": tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Nissan", "Nissan"),
			tgbotapi.NewInlineKeyboardButtonData("Porsche", "Porsche"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Toyota", "Toyota"),
			tgbotapi.NewInlineKeyboardButtonData("Vauxhall", "Vauxhall"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Volkswagen", "Volkswagen"),
			tgbotapi.NewInlineKeyboardButtonData("Volvo", "Volvo"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Back", "popular"),
			tgbotapi.NewInlineKeyboardButtonData("Cancel", "cancel"),
		),
	),
	"begin": tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Popular Makes", "popular"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("!German", "german"),
			tgbotapi.NewInlineKeyboardButtonData("!British", "british"),
			tgbotapi.NewInlineKeyboardButtonData("!Japanese", "japanese"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("!Other", "other"),
			tgbotapi.NewInlineKeyboardButtonData("Cancel", "cancel"),
		),
	),
}

func getInlineKeyboard(keyboard string) tgbotapi.InlineKeyboardMarkup {

	return keyboardMap[keyboard]
}
