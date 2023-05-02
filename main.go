package main

import (
	"encoding/json"
	"fmt"

	"context"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/procyon-projects/chrono"
)

const DATA_FILE = "./ids.json"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error while loading env variables: %s", err)
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatalf("Error occured while creating Telegram bot: %s", err)
	}

	bot.Debug = true

	taskScheduler := chrono.NewDefaultTaskScheduler()

	_, err = taskScheduler.ScheduleWithCron(func(ctx context.Context) {
		log.Println("Sending daily floppas...")
		go func() {
			if err = floppinson(bot); err != nil {
				log.Printf("Failed to send daily floppas: %s", err)
			}
		}()
	}, "0 0 9 * * *") // Every day at 9:00 AM
	if err != nil {
		log.Printf("Failed to schedule daily floppas: %s", err)
	}
	log.Print("Daily floppas were scheduled successfully.")

	log.Printf("Authorized on account %s", bot.Self.UserName)

	err = initCommands(bot)
	if err != nil {
		log.Fatalf("Failed to register commands: %s", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.Text == "flop" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "FLOP FLOP!")
				if _, err := bot.Send(msg); err != nil {
					log.Printf("flop: Failed to send message: %s", err)
				}
			}

			switch update.Message.Command() {
			case "subscribe":
				ids, err := getSubscriberIDs()
				if err != nil {
					log.Printf("subscribe: Failed to get subscriber ids: %s", err)
				}

				id := update.Message.Chat.ID
				if !contains(ids, id) {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "SUBSCRIBED TO FLOPPA PHOTOS!")
					if _, err = bot.Send(msg); err != nil {
						log.Printf("subscribe: Failed to send message: %s", err)
					}
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "FLOP FLOP!")
					if _, err = bot.Send(msg); err != nil {
						log.Printf("subscribe: Failed to send message: %s", err)
					}

					err = addNewSubscriber(id)
					if err != nil {
						log.Printf("subscribe: Failed to subscribe: %s", err)
					}
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You can only subscribe once")
					if _, err = bot.Send(msg); err != nil {
						log.Printf("subscribe: Failed to send message: %s", err)
					}
				}
			case "floppinson":
				go func() {
					if err = floppinson(bot); err != nil {
						log.Printf("Failed to send daily floppas: %s", err)
					}
				}()
			case "floppik":
				go func() {
					err = flopik(bot, update.FromChat().ID)
					if err != nil {
						log.Printf("flopik: failed to send flopik: %s", err)
					}
				}()
			case "earrape":
				go func() {
					ids, err := getSubscriberIDs()
					if err != nil {
						log.Printf("earrape: Failed to get subscriber ids: %s", err)
					}

					for _, id := range ids {
						photoBytes, err := os.ReadFile("video/earrape.mp4")
						if err != nil {
							fmt.Printf("earrape: Failed to open video file: %s", err)
						}
						photoFileBytes := tgbotapi.FileBytes{
							Name:  "Flopik",
							Bytes: photoBytes,
						}
						_, err = bot.Send(tgbotapi.NewVideo(int64(id), photoFileBytes))
						if err != nil {
							fmt.Printf("earrape: Failed to send video: %s", err)
						}
					}
				}()
			case "ids":
				go func() {
					ids, err := getSubscriberIDs()
					if err != nil {
						log.Printf("ids: Failed to get subscriber ids: %s", err)
					}

					var strarr []string
					for _, id := range ids {
						strarr = append(strarr, strconv.FormatInt(id, 10))
					}
					str := strings.Join(strarr, ",")

					msg := tgbotapi.NewMessage(update.Message.Chat.ID, str)
					if _, err = bot.Send(msg); err != nil {
						log.Printf("ids: Failed to send message: %s", err)
					}
				}()
			case "announce":
				go func() {
					message := strings.Replace(update.Message.Text, "/announce ", "", 644)

					ids, err := getSubscriberIDs()
					if err != nil {
						log.Printf("announce: Failed to get subscriber ids: %s", err)
					}

					for _, id := range ids {
						msg := tgbotapi.NewMessage(id, message)
						if _, err = bot.Send(msg); err != nil {
							log.Printf("announce: Failed to send message: %s", err)
						}
					}
				}()
			case "flop":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "FLOP FLOP!")
				if _, err = bot.Send(msg); err != nil {
					log.Printf("flop: Failed to send message: %s", err)
				}
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Type /subscribe to get daily floppas!")
				if _, err = bot.Send(msg); err != nil {
					log.Printf("start: Failed to send message: %s", err)
				}
			case "chat":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, strconv.FormatInt(update.Message.Chat.ID, 10))
				if _, err = bot.Send(msg); err != nil {
					log.Printf("chat: Failed to send message: %s", err)
				}
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Co")
				if _, err = bot.Send(msg); err != nil {
					log.Printf("default: Failed to send message: %s", err)
				}
			}
		}
	}
}

func initCommands(bot *tgbotapi.BotAPI) error {
	commandsEveryone := []tgbotapi.BotCommand{
		{Command: "/subscribe", Description: "Subscribe to daily floppas"},
		{Command: "/floppik", Description: "Get floppa"},
		{Command: "/flop", Description: "flop"},
		{Command: "/chat", Description: "Returns chat id"},
	}

	commandsAdmin := []tgbotapi.BotCommand{
		{Command: "/floppinson", Description: "Manually send daily floppas to all subscribers"},
		{Command: "/earrape", Description: "Send earrape floppa video to all subscribers"},
		{Command: "/ids", Description: "Get all subscriber ids"},
		{Command: "/announce", Description: "Send an announcement to all subscribers"},
	}
	commandsAdmin = append(commandsAdmin, commandsEveryone...)

	scopeEveryone := tgbotapi.NewBotCommandScopeDefault()
	adminChatId, err := strconv.ParseInt(os.Getenv("ADMIN_CHAT_ID"), 10, 64)
	if err != nil {
		return err
	}
	scopeAdmin := tgbotapi.NewBotCommandScopeChat(adminChatId)

	config := tgbotapi.NewSetMyCommandsWithScope(scopeEveryone, commandsEveryone...)
	_, err = bot.Request(config)
	if err != nil {
		return err
	}

	config = tgbotapi.NewSetMyCommandsWithScope(scopeAdmin, commandsAdmin...)
	_, err = bot.Request(config)
	return err
}

func contains(s []int64, str int64) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// floppinson sends a random floppa image to a list of users loaded from a JSON file
func floppinson(bot *tgbotapi.BotAPI) error {
	file, err := os.ReadFile(DATA_FILE)
	if err != nil {
		return err
	}

	var arr []int64
	err = json.Unmarshal(file, &arr)
	if err != nil {
		return err
	}

	// Iterates over every user in the list and sends them a random floppa image
	for _, id := range arr {
		err = flopik(bot, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func flopik(bot *tgbotapi.BotAPI, id int64) error {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	pictureId := rng.Intn(32)

	photoBytes, err := os.ReadFile(fmt.Sprintf("floppa/%d.jpg", pictureId))
	if err != nil {
		return err
	}

	photoFileBytes := tgbotapi.FileBytes{
		Name:  "Flopik",
		Bytes: photoBytes,
	}

	_, err = bot.Send(tgbotapi.NewPhoto(id, photoFileBytes))
	return err
}
