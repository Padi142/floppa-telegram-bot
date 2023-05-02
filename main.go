package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/procyon-projects/chrono"
	"log"
	"os"
)

type telegramBot struct {
	tgbot *tgbotapi.BotAPI
}

const DATA_FILE = "./ids.json"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error while loading env variables: %s", err)
	}

	tgbot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatalf("Error occured while creating Telegram bot: %s", err)
	}
	bot := telegramBot{tgbot: tgbot}

	bot.tgbot.Debug = true

	taskScheduler := chrono.NewDefaultTaskScheduler()

	_, err = taskScheduler.ScheduleWithCron(func(ctx context.Context) {
		log.Println("Sending daily floppas...")
		go func() {
			if err = bot.floppinson(); err != nil {
				log.Printf("Failed to send daily floppas: %s", err)
			}
		}()
	}, "0 0 9 * * *") // Every day at 9:00 AM
	if err != nil {
		log.Printf("Failed to schedule daily floppas: %s", err)
	}
	log.Print("Daily floppas were scheduled successfully.")

	log.Printf("Authorized on account %s", bot.tgbot.Self.UserName)

	err = bot.initCommands()
	if err != nil {
		log.Fatalf("Failed to register commands: %s", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.tgbot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.Text == "flop" {
				bot.flop(update)
			}

			switch update.Message.Command() {
			case "subscribe":
				bot.subscribe(update)
			case "floppinson":
				go func() {
					if err = bot.floppinson(); err != nil {
						log.Printf("Failed to send daily floppas: %s", err)
					}
				}()
			case "floppik":
				go func() {
					if err = bot.flopik(update.FromChat().ID); err != nil {
						log.Printf("flopik: failed to send flopik: %s", err)
					}
				}()
			case "earrape":
				go func() {
					bot.earrape()
				}()
			case "ids":
				go func() {
					bot.ids(update)
				}()
			case "announce":
				go func() {
					bot.announce(update)
				}()
			case "flop":
				bot.flop(update)
			case "start":
				bot.start(update)
			case "chat":
				bot.chat(update)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Co")
				if _, err = bot.tgbot.Send(msg); err != nil {
					log.Printf("default: Failed to send message: %s", err)
				}
			}
		}
	}
}

func contains(s []int64, str int64) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
