package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

type telegramMessage struct {
	ChatID              int    `json:"chat_id"`
	Text                string `json:"text"`
	DisableNotification bool   `json:"disable_notification"`
}

func sendTelegramMessage(message string, DisableNotification bool) error {
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramChatID, err := strconv.Atoi(os.Getenv("TELEGRAM_CHAT_ID"))

	if err != nil {
		return err
	}

	botMessage := telegramMessage{
		ChatID:              telegramChatID,
		Text:                message,
		DisableNotification: DisableNotification,
	}
	botMessageBytes, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	telegramAPIURL := "https://api.telegram.org/bot" + telegramBotToken + "/sendMessage"
	resp, err := http.Post(telegramAPIURL, "application/json", bytes.NewBuffer(botMessageBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
