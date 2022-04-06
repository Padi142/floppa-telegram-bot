package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("5187849216:AAE0mAjGV7wLf0ER-fp2jR5gD7oz8qzp8ek")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.Text == "flop" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "FLOP FLOP!")

				bot.Send(msg)
			}
			switch update.Message.Command() {
			case "subscribe":

				file, err := ioutil.ReadFile("ids.json")
				var arr []int64
				json.Unmarshal(file, &arr)
				id := update.Message.Chat.ID
				if !contains(arr, id) {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "SUBSCRIBED TO FLOPPA PHOTOS!")
					bot.Send(msg)
					msg = tgbotapi.NewMessage(update.Message.Chat.ID, "FLOP FLOP!")
					bot.Send(msg)

					arr = append(arr, id)

					file, err = json.Marshal(arr)

					err = ioutil.WriteFile("ids.json", file, 0644)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("saving id: " + strconv.Itoa(int(update.Message.Chat.ID)))
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You can only subscribe once")
					bot.Send(msg)
				}
			case "floppinson":
				go func() {
					file, err := ioutil.ReadFile("ids.json")
					var arr []int64
					json.Unmarshal(file, &arr)
					for index := 0; index < len(arr); index++ {
						s1 := rand.NewSource(time.Now().UnixNano())
						rng := rand.New(s1)
						picture := rng.Intn(32)

						id := arr[index]
						photoBytes, err := ioutil.ReadFile("floppa/" + strconv.Itoa(picture) + ".jpg")
						if err != nil {
							panic(err)
						}
						photoFileBytes := tgbotapi.FileBytes{
							Name:  "Flopik",
							Bytes: photoBytes,
						}
						_, err = bot.Send(tgbotapi.NewPhoto(int64(id), photoFileBytes))
					}
					if err != nil {
						fmt.Println(err)
					}
				}()
			case "earrape":
				go func() {
					file, err := ioutil.ReadFile("ids.json")
					var arr []int64
					json.Unmarshal(file, &arr)
					for index := 0; index < len(arr); index++ {

						id := arr[index]
						photoBytes, err := ioutil.ReadFile("video/earrape.mp4")
						if err != nil {
							panic(err)
						}
						photoFileBytes := tgbotapi.FileBytes{
							Name:  "Flopik",
							Bytes: photoBytes,
						}
						_, err = bot.Send(tgbotapi.NewVideo(int64(id), photoFileBytes))
					}
					if err != nil {
						fmt.Println(err)

					}
				}()
			case "ids":
				go func() {
					file, err := ioutil.ReadFile("ids.json")
					var arr []int
					var IDs []string
					json.Unmarshal(file, &arr)
					for _, i := range arr {
						IDs = append(IDs, strconv.Itoa(i))
					}
					idstring := ""
					for _, id := range IDs {
						idstring = idstring + "," + id
					}
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, idstring)
					bot.Send(msg)
					if err != nil {
						fmt.Println(err)
					}
				}()
			case "announce":
				go func() {

					message := strings.Replace(update.Message.Text, "/announce ", "", 644)

					file, err := ioutil.ReadFile("ids.json")
					var arr []int64
					json.Unmarshal(file, &arr)
					for index := 0; index < len(arr); index++ {

						msg := tgbotapi.NewMessage(arr[index], message)
						bot.Send(msg)

					}
					if err != nil {
						fmt.Println(err)
					}
				}()
			case "flop":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "FLOP FLOP!")
				bot.Send(msg)
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Type /subscribe to get daily floppas!")
				bot.Send(msg)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Co")
				bot.Send(msg)
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
