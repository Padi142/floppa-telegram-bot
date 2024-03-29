package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *telegramBot) initCommands() error {
	commandsEveryone := []tgbotapi.BotCommand{
		{Command: "/subscribe", Description: "Subscribe to daily floppas"},
		{Command: "/floppik", Description: "Get floppa"},
		{Command: "/flop", Description: "flop"},
		{Command: "/chat", Description: "Returns chat id"},
		{Command: "/unsubscribe", Description: "Unsubscribes you from daily floppas"},
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

func (b *telegramBot) chat(update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, strconv.FormatInt(update.Message.Chat.ID, 10))
	_, err := b.tgbot.Send(msg)
	return err
}

func (b *telegramBot) start(update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Type /subscribe to get daily floppas!")
	_, err := b.tgbot.Send(msg)
	return err
}

func (b *telegramBot) announce(update tgbotapi.Update) error {
	message := strings.Replace(update.Message.Text, "/announce ", "", 644)

	ids, err := getSubscriberIDs()
	if err != nil {
		return err
	}

	for _, id := range ids {
		msg := tgbotapi.NewMessage(id, message)
		if _, err = b.tgbot.Send(msg); err != nil {
			log.Printf("announce: Failed to send message: %s", err)
		}
	}
	return nil
}

func (b *telegramBot) ids(update tgbotapi.Update) error {
	ids, err := getSubscriberIDs()
	if err != nil {
		return err
	}

	var strarr []string
	for _, id := range ids {
		strarr = append(strarr, strconv.FormatInt(id, 10))
	}
	str := strings.Join(strarr, ",")

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, str)
	_, err = b.tgbot.Send(msg)
	return err
}

func (b *telegramBot) earrape() error {
	ids, err := getSubscriberIDs()
	if err != nil {
		return err
	}

	for _, id := range ids {
		photoBytes, err := os.ReadFile("video/earrape.mp4")
		if err != nil {
			return err
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
	return nil
}

func (b *telegramBot) subscribe(update tgbotapi.Update) error {
	ids, err := getSubscriberIDs()
	if err != nil {
		return err
	}

	id := update.Message.Chat.ID
	if !contains(ids, id) {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "SUBSCRIBED TO FLOPPA PHOTOS!")
		if _, err = b.tgbot.Send(msg); err != nil {
			return err
		}
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "FLOP FLOP!")
		if _, err = b.tgbot.Send(msg); err != nil {
			return err
		}

		return addNewSubscriber(id)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "You can only subscribe once")
	_, err = b.tgbot.Send(msg)
	return err
}

func (b *telegramBot) unsubscribe(update tgbotapi.Update) error {
	// msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Type /subscribe to get daily floppas!")

	sender := update.Message.From
	if(sender.UserName=="LufyCZ"){
		_, err := b.tgbot.Request(tgbotapi.KickChatMemberConfig{
			ChatMemberConfig: tgbotapi.ChatMemberConfig{
				ChatID: update.Message.Chat.ID,
				UserID: sender.ID,
			},
			UntilDate: 1,
		})
		if(err != nil){
			log.Print(err)
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, sender.UserName+" was kicked from this chat")
		_, err = b.tgbot.Send(msg)
	}

	 err := b.sprcha(update.Message.Chat.ID)

	// _, err := b.tgbot.Send(msg)
	return err
}

func (b *telegramBot) flop(update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "FLOP FLOP!")
	_, err := b.tgbot.Send(msg)
	return err
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
			log.Printf("floppik: Failed to send floppik: %s", err)
		}
	}

	return nil
}

func (b *telegramBot) flopik(id int64) error {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	pictureId := rng.Intn(41)

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

func (b *telegramBot) sprcha(id int64) error {

	photoBytes, err := os.ReadFile(fmt.Sprintf("floppa/sprcha.png"))
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
