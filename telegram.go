package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const telegramBotToken = "xxxxx"
const telegramChatID = -81321321

type telegramMessage struct {
	ChatID              int64  `json:"chat_id"`
	Text                string `json:"text"`
	DisableNotification bool   `json:"disable_notification"`
}

func sendTelegramMessage(message string, DisableNotification bool) error {

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
