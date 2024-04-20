package telegram

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type TelegramMessage struct {
	ChatID              int    `json:"chat_id"`
	Text                string `json:"text"`
	DisableNotification bool   `json:"disable_notification"`
	ParseMode           string `json:"parse_mode"`
}

func Send(chatID int, message string, DisableNotification bool) error {
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	enableTelegramNotification, err := strconv.ParseBool(os.Getenv("ENABLE_TELEGRAM_NOTIFICATION"))
	if err != nil || !enableTelegramNotification {
		log.Debug().Msg("[TELEGRAM] " + message)
		log.Debug().Msg("env 'ENABLE_TELEGRAM_NOTIFICATION' is not found or false. Skipping telegram notification")
		return nil
	}

	botMessage := TelegramMessage{
		ChatID:              chatID,
		Text:                message,
		DisableNotification: DisableNotification,
		ParseMode:           "MarkdownV2",
	}
	botMessageBytes, err := json.Marshal(botMessage)
	if err != nil {
		log.Info().Msg("ERROR" + err.Error())
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
		log.Info().Msg("ERROR" + err.Error())
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
		log.Error().Msg(err.Error())
	}

	if strings.Contains(strings.ToLower(OnUpdateMessageBody.Message.Text), "/hello") {
		Send(telegramChatID, "Hello world!", false)
		return
	}

	jsonStringByte, err := json.MarshalIndent(OnUpdateMessageBody, "", "  ")
	if err != nil {
		log.Error().Msg(err.Error())
	}

	log.Print(string(jsonStringByte))
}
