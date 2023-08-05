package telegram

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type TelegramMessage struct {
	ChatID              int    `json:"chat_id"`
	Text                string `json:"text"`
	DisableNotification bool   `json:"disable_notification"`
}

func Send(chatID int, message string, DisableNotification bool) error {
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	enableTelegramNotification, err := strconv.ParseBool(os.Getenv("ENABLE_TELEGRAM_NOTIFICATION"))
	if err != nil || !enableTelegramNotification {
		log.Println("[TELEGRAM] " + message)
		log.Println("env 'ENABLE_TELEGRAM_NOTIFICATION' is not found or false. Skipping telegram notification")
		return nil
	}

	botMessage := TelegramMessage{
		ChatID:              chatID,
		Text:                message,
		DisableNotification: DisableNotification,
	}
	botMessageBytes, err := json.Marshal(botMessage)
	if err != nil {
		log.Println("ERROR" + err.Error())
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

type SetWebhookBody struct {
	URL string `json:"url"`
}

func SetWebhook(url string) error {
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	telegramAPIURL := "https://api.telegram.org/bot" + telegramBotToken + "/setWebhook"

	body := SetWebhookBody{
		URL: url,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Println("ERROR" + err.Error())
		return err
	}

	resp, err := http.Post(telegramAPIURL, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func OnUpdateMessage(OnUpdateMessageBody OnUpdateMessageBody) {
	telegramChatID, err := strconv.Atoi(os.Getenv("TELEGRAM_CHAT_ID"))
	if err != nil {
		log.Println("ERROR " + err.Error())
	}

	if strings.Contains(strings.ToLower(OnUpdateMessageBody.Message.Text), "/hello") {
		Send(telegramChatID, "Hello world!", false)
		return
	}

	jsonStringByte, err := json.MarshalIndent(OnUpdateMessageBody, "", "  ")
	if err != nil {
		log.Println(err)
	}

	log.Print(string(jsonStringByte))
}
