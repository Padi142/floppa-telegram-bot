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
				if err = bot.flop(update); err != nil {
					log.Printf("flop: Error: %s", err)
				}
			}

			switch update.Message.Command() {
			case "subscribe":
				if err = bot.subscribe(update); err != nil {
					log.Printf("subscribe: Failed to subscribe: %s", err)
				}
			case "floppinson":
				go func() {
					if err = bot.floppinson(); err != nil {
						log.Printf("floppinson: Failed to send daily floppas: %s", err)
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
					if err = bot.earrape(); err != nil {
						log.Printf("earrape: Failed to send earrape: %s", err)
					}
				}()
			case "ids":
				go func() {
					if err = bot.ids(update); err != nil {
						log.Printf("ids: Failed to get subsriber ids: %s", err)
					}
				}()
			case "announce":
				go func() {
					if err = bot.announce(update); err != nil {
						log.Printf("announce: Failed to send announcement to subscribers: %s", err)
					}
				}()
			case "flop":
				if err = bot.flop(update); err != nil {
					log.Printf("flop: Error: %s", err)
				}
			case "start":
				if err = bot.start(update); err != nil {
					log.Printf("start: Error: %s", err)
				}
			case "chat":
				if err = bot.chat(update); err != nil {
					log.Printf("chat: Error: %s", err)
				}
			case "unsubscribe":
				if err = bot.unsubscribe(update); err != nil {
					log.Printf("chat: Error: %s", err)
				}
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
