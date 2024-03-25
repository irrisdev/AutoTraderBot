package bot

import (
	"fmt"
	telebot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
	"os"
)

func TeleBot() {

	bot, err := telebot.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Err(err).Msg("Failed to start telegram bot")
		return
	}

	log.Info().Msg("Authorised on account " + bot.Self.UserName)

	bot.Debug = true
	updateConfig := telebot.NewUpdate(0)
	updateConfig.Timeout = 60

	//Subscribe to updates
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {

		if update.Message != nil {

			msg := telebot.NewMessage(update.Message.Chat.ID, "")

			if update.Message.Text == "Cancel" {
				deleteMsg(bot, update.Message.Chat.ID, update.Message.MessageID)
			}

			switch update.Message.Command() {
			case "start":
				inlineKeyboard := telebot.NewInlineKeyboardMarkup(
					telebot.NewInlineKeyboardRow(
						telebot.NewInlineKeyboardButtonData("Begin Tracking Average Price", "begin"),
					),
					telebot.NewInlineKeyboardRow(
						telebot.NewInlineKeyboardButtonData("Cancel", "cancel"),
					),
				)
				msg.Text = fmt.Sprintf("Welcome, @%s! Please select an option", update.Message.Chat.UserName)
				msg.ReplyMarkup = inlineKeyboard
			}

			send(bot, msg)
		} else if update.CallbackQuery != nil {

			//Extract callbackId and the data
			data := update.CallbackQuery.Data
			inlineKeyboardId := update.CallbackQuery.Message.MessageID
			chatId := update.CallbackQuery.Message.Chat.ID

			switch data {

			case "cancel":

				deleteMsg(bot, chatId, inlineKeyboardId)
				deleteMsg(bot, chatId, inlineKeyboardId-1)

			case "begin":
				deleteMsg(bot, chatId, inlineKeyboardId)

				newReplyKeyboard := telebot.NewReplyKeyboard(
					telebot.NewKeyboardButtonRow(
						telebot.NewKeyboardButton("Volkswagen"),
					),
					telebot.NewKeyboardButtonRow(
						telebot.NewKeyboardButton("Other"),
					),
					telebot.NewKeyboardButtonRow(
						telebot.NewKeyboardButton("Back"),
						telebot.NewKeyboardButton("Cancel"),
						telebot.NewKeyboardButton("Next"),
					),
				)

				newMsg := telebot.NewMessage(chatId, "Choose the make of the vehicle:")
				newMsg.ReplyMarkup = newReplyKeyboard

				if _, err := bot.Send(newMsg); err != nil {
					log.Err(err).Msg("Error while sending new message")
				}
				_, err := bot.Send(telebot.NewCallback(update.CallbackQuery.ID, ""))
				if err != nil {
					return
				}

			}

		}
	}

}

func deleteMsg(bot *telebot.BotAPI, chatId int64, messageId int) {
	deleteMsg := telebot.NewDeleteMessage(chatId, messageId)
	if _, err := bot.Request(deleteMsg); err != nil {
		log.Err(err).Msg("Failed to delete message")
	}
}
