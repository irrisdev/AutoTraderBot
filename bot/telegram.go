package bot

//
//import (
//	"fmt"
//	telebot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
//	"github.com/rs/zerolog/log"
//	"os"
//	"strconv"
//	"time"
//)
//
//type sessionState struct {
//	State int
//}
//
//func StartTelegramBot(done chan<- string) {
//	var numericKeyboard = telebot.NewInlineKeyboardMarkup(
//		telebot.NewInlineKeyboardRow(
//			telebot.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
//			telebot.NewInlineKeyboardButtonData("2", "2"),
//			telebot.NewInlineKeyboardButtonData("3", "3"),
//		),
//		telebot.NewInlineKeyboardRow(
//			telebot.NewInlineKeyboardButtonData("4", "4"),
//			telebot.NewInlineKeyboardButtonData("5", "5"),
//			telebot.NewInlineKeyboardButtonData("6", "6"),
//		),
//	)
//	_ = numericKeyboard
//	//Login into bot
//	bot, err := telebot.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
//	if err != nil {
//		log.Err(err).Msg("Failed to start telegram bot")
//		return
//	}
//	log.Info().Msg("Started telegram bot")
//	log.Info().Msg("Authorised on account " + bot.Self.UserName)
//
//	//Initialise sessionState map
//	sessions := make(map[int64]*sessionState)
//	adminUser, _ := strconv.ParseInt(os.Getenv("ADMIN_USER"), 10, 64)
//	allowedUsers := make(map[int64]struct{})
//	allowedUsers[adminUser] = struct{}{}
//
//	// Start polling Telegram for updates.
//	bot.Debug = true
//	updateConfig := telebot.NewUpdate(0)
//	updateConfig.Timeout = 60
//	updates := bot.GetUpdatesChan(updateConfig)
//
//	for update := range updates {
//
//		//Discard empty messages as if no updates are detected in 60 seconds telegram will respond to request with empty struct
//		if update.Message == nil {
//			continue
//		}
//
//		if _, allowed := allowedUsers[update.Message.Chat.ID]; !allowed {
//			notAuthorised := telebot.NewMessage(update.Message.Chat.ID, "Unauthorised User")
//			notAuthorised.ReplyToMessageID = update.Message.MessageID
//			_, err := bot.Send(notAuthorised)
//			if err != nil {
//				return
//			}
//			continue
//		}
//
//		session, exists := sessions[update.Message.From.ID]
//		if !exists {
//			session = &sessionState{State: 1}
//			sessions[update.Message.From.ID] = session
//		}
//
//		msg := telebot.NewMessage(update.Message.Chat.ID, "")
//
//		if update.FromChat().ID == adminUser {
//			switch update.Message.Command() {
//			case "adduser":
//				username := update.Message.CommandArguments()
//				if username == "" {
//					msg.Text = "Usage /command username"
//				} else {
//					num, _ := strconv.ParseInt(username, 10, 64)
//					allowedUsers[num] = struct{}{}
//					msg.Text = fmt.Sprintf("Added user %s to allow list", username)
//				}
//
//				send(bot, msg)
//
//			case "removeuser":
//				username := update.Message.CommandArguments()
//				if username == "" {
//					msg.Text = "Usage /command username"
//				} else {
//					num, _ := strconv.ParseInt(username, 10, 64)
//					delete(allowedUsers, num)
//					msg.Text = fmt.Sprintf("Removed user %s from allow list", username)
//				}
//				send(bot, msg)
//
//			}
//
//		}
//
//		switch session.State {
//		case 1:
//			if update.Message.IsCommand() {
//				switch update.Message.Command() {
//				case "start":
//					msg.Text, err = handleStart(session)
//					if err != nil {
//						log.Err(err).Msg("Error occurred during /start")
//					}
//				case "help":
//					msg.Text = "Use /start to begin"
//				default:
//					continue
//				}
//			}
//		case 2:
//			switch update.Message.Command() {
//
//			case "help":
//				msg.Text = "Use /reset to reset state or /state to check current state"
//			case "reset":
//				msg.Text = resetSession(session, update)
//			case "state":
//				msg.Text = fmt.Sprintf("Current state of session %d = State: %d", update.Message.From.ID, session.State)
//			default:
//				msg.Text = "Use /help to see commands available at state 2"
//			}
//
//		default:
//			msg.Text = "Unknown Command"
//		}
//
//		send(bot, msg)
//	}
//
//	time.Sleep(time.Second * 5)
//}
//
//func resetSession(session *sessionState, update telebot.Update) string {
//	originalSession := *session
//	originalSession.State = session.State
//
//	session.State = 1
//	return fmt.Sprintf("State for session %s:%d has been reset from %d to %d", update.FromChat().UserName, update.FromChat().ID, originalSession.State, session.State)
//}
//
//func handleStart(session *sessionState) (string, error) {
//
//	session.State = 2
//
//	return "Success of /start command", nil
//}
