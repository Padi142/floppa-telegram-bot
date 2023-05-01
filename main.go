package main

import (
	"encoding/json"
	"fmt"

	"context"
	"io/ioutil"
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
				file, err := os.ReadFile(DATA_FILE)
				if err != nil {
					log.Printf("subscribe: Failed to open data file: %s", err)
				}

				var arr []int64
				err = json.Unmarshal(file, &arr)
				if err != nil {
					log.Printf("subscribe: Failed to unmarshal JSON file: %s", err)
				}

				id := update.Message.Chat.ID
				if !contains(arr, id) {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "SUBSCRIBED TO FLOPPA PHOTOS!")
					if _, err = bot.Send(msg); err != nil {
						log.Printf("subscribe: Failed to send message: %s", err)
					}
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "FLOP FLOP!")
					if _, err = bot.Send(msg); err != nil {
						log.Printf("subscribe: Failed to send message: %s", err)
					}

					arr = append(arr, id)
					file, err = json.Marshal(arr)
					if err != nil {
						log.Printf("subscribe: Failed to marshal JSON file: %s", err)
					}

					err = os.WriteFile(DATA_FILE, file, 0644)
					if err != nil {
						fmt.Printf("subscribe: Failed to write JSON into data file: %s", err)
					}
					fmt.Println("saving id: " + strconv.Itoa(int(update.Message.Chat.ID)))
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
					s1 := rand.NewSource(time.Now().UnixNano())
					rng := rand.New(s1)
					picture := rng.Intn(32)

					photoBytes, err := ioutil.ReadFile("floppa/" + strconv.Itoa(picture) + ".jpg")
					if err != nil {
						panic(err)
					}
					photoFileBytes := tgbotapi.FileBytes{
						Name:  "Flopik",
						Bytes: photoBytes,
					}
					_, err = bot.Send(tgbotapi.NewPhoto(update.FromChat().ID, photoFileBytes))
				}()
			case "earrape":
				go func() {
					file, err := os.ReadFile(DATA_FILE)
					if err != nil {
						log.Printf("subscribe: Failed to open data file: %s", err)
					}

					var arr []int64
					err = json.Unmarshal(file, &arr)
					if err != nil {
						log.Printf("subscribe: Failed to unmarshal JSON file: %s", err)
					}

					for index := 0; index < len(arr); index++ {

						id := arr[index]
						photoBytes, err := ioutil.ReadFile("video/earrape.mp4")
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
					file, err := os.ReadFile(DATA_FILE)
					if err != nil {
						log.Printf("subscribe: Failed to open data file: %s", err)
					}
					var arr []int
					var IDs []string
					err = json.Unmarshal(file, &arr)
					if err != nil {
						log.Printf("subscribe: Failed to unmarshal JSON file: %s", err)
					}
					for _, i := range arr {
						IDs = append(IDs, strconv.Itoa(i))
					}
					idstring := ""
					for _, id := range IDs {
						idstring = idstring + "," + id
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, idstring)
					if _, err = bot.Send(msg); err != nil {
						log.Printf("ids: Failed to send message: %s", err)
					}
				}()
			case "announce":
				go func() {
					message := strings.Replace(update.Message.Text, "/announce ", "", 644)

					file, err := os.ReadFile(DATA_FILE)
					if err != nil {
						log.Printf("subscribe: Failed to open data file: %s", err)
					}
					var arr []int64
					err = json.Unmarshal(file, &arr)
					for index := 0; index < len(arr); index++ {
						msg := tgbotapi.NewMessage(arr[index], message)
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
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Co")
				if _, err = bot.Send(msg); err != nil {
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
		if err != nil {
			return err
		}
	}

	return nil
}
