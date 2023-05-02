package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func (b *telegramBot) initCommands() error {
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
	_, err = b.tgbot.Request(config)
	if err != nil {
		return err
	}

	config = tgbotapi.NewSetMyCommandsWithScope(scopeAdmin, commandsAdmin...)
	_, err = b.tgbot.Request(config)
	return err
}

func (b *telegramBot) chat(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, strconv.FormatInt(update.Message.Chat.ID, 10))
	if _, err := b.tgbot.Send(msg); err != nil {
		log.Printf("chat: Failed to send message: %s", err)
	}
}

func (b *telegramBot) start(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Type /subscribe to get daily floppas!")
	if _, err := b.tgbot.Send(msg); err != nil {
		log.Printf("start: Failed to send message: %s", err)
	}
}

func (b *telegramBot) announce(update tgbotapi.Update) {
	message := strings.Replace(update.Message.Text, "/announce ", "", 644)

	ids, err := getSubscriberIDs()
	if err != nil {
		log.Printf("announce: Failed to get subscriber ids: %s", err)
	}

	for _, id := range ids {
		msg := tgbotapi.NewMessage(id, message)
		if _, err = b.tgbot.Send(msg); err != nil {
			log.Printf("announce: Failed to send message: %s", err)
		}
	}
}

func (b *telegramBot) ids(update tgbotapi.Update) {
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
	if _, err = b.tgbot.Send(msg); err != nil {
		log.Printf("ids: Failed to send message: %s", err)
	}
}

func (b *telegramBot) earrape() {
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
		_, err = b.tgbot.Send(tgbotapi.NewVideo(int64(id), photoFileBytes))
		if err != nil {
			fmt.Printf("earrape: Failed to send video: %s", err)
		}
	}
}

func (b *telegramBot) subscribe(update tgbotapi.Update) {
	ids, err := getSubscriberIDs()
	if err != nil {
		log.Printf("subscribe: Failed to get subscriber ids: %s", err)
	}

	id := update.Message.Chat.ID
	if !contains(ids, id) {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "SUBSCRIBED TO FLOPPA PHOTOS!")
		if _, err = b.tgbot.Send(msg); err != nil {
			log.Printf("subscribe: Failed to send message: %s", err)
		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "FLOP FLOP!")
		if _, err = b.tgbot.Send(msg); err != nil {
			log.Printf("subscribe: Failed to send message: %s", err)
		}

		err = addNewSubscriber(id)
		if err != nil {
			log.Printf("subscribe: Failed to subscribe: %s", err)
		}
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You can only subscribe once")
		if _, err = b.tgbot.Send(msg); err != nil {
			log.Printf("subscribe: Failed to send message: %s", err)
		}
	}
}

func (b *telegramBot) flop(update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "FLOP FLOP!")
	if _, err := b.tgbot.Send(msg); err != nil {
		log.Printf("flop: Failed to send message: %s", err)
	}
}

// floppinson sends a random floppa image to a list of users loaded from a JSON file
func (b *telegramBot) floppinson() error {
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
		err = b.flopik(id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *telegramBot) flopik(id int64) error {
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

	_, err = b.tgbot.Send(tgbotapi.NewPhoto(id, photoFileBytes))
	return err
}
