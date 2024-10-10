package telegram

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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

type TelegramResponse struct {
	OK          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
}

type SetWebhookBody struct {
	URL string `json:"url"`
}

func escapeMarkdownV2(text string) string {
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, "\\"+char)
	}
	return text
}

func Send(chatID int, message string, DisableNotification bool) error {
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	enableTelegramNotification, err := strconv.ParseBool(os.Getenv("ENABLE_TELEGRAM_NOTIFICATION"))
	if err != nil || !enableTelegramNotification {
		log.Debug().Msg("[TELEGRAM] " + message)
		log.Debug().Msg("env 'ENABLE_TELEGRAM_NOTIFICATION' is not found or false. Skipping telegram notification")
		return nil
	}

	escapedMessage := escapeMarkdownV2(message)

	botMessage := TelegramMessage{
		ChatID:              chatID,
		Text:                escapedMessage,
		DisableNotification: DisableNotification,
		ParseMode:           "MarkdownV2",
	}
	botMessageBytes, err := json.Marshal(botMessage)
	if err != nil {
		log.Error().Err(err).Msg("Error marshaling message")
		return err
	}
	telegramAPIURL := "https://api.telegram.org/bot" + telegramBotToken + "/sendMessage"

	resp, err := http.Post(telegramAPIURL, "application/json", bytes.NewBuffer(botMessageBytes))
	if err != nil {
		log.Error().Err(err).Msg("Error sending POST request")
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading response body")
		return err
	}

	var telegramResp TelegramResponse
	if err := json.Unmarshal(body, &telegramResp); err != nil {
		log.Error().Err(err).Msg("Error unmarshaling response")
		return err
	}

	if !telegramResp.OK {
		log.Error().
			Int("error_code", telegramResp.ErrorCode).
			Str("description", telegramResp.Description).
			Msg("Telegram API error")
		return err
	}

	log.Info().Msg("Message sent successfully to Telegram")
	return nil
}

func SetWebhook(url string) error {
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")

	telegramAPIURL := "https://api.telegram.org/bot" + telegramBotToken + "/setWebhook"

	body := SetWebhookBody{
		URL: url,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Error().Err(err).Msg("Error marshaling webhook body")
		return err
	}

	resp, err := http.Post(telegramAPIURL, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Error().Err(err).Msg("Error sending POST request")
		return err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("Error reading response body")
		return err
	}

	var telegramResp TelegramResponse
	if err := json.Unmarshal(responseBody, &telegramResp); err != nil {
		log.Error().Err(err).Msg("Error unmarshaling response")
		return err
	}

	if !telegramResp.OK {
		log.Error().
			Int("error_code", telegramResp.ErrorCode).
			Str("description", telegramResp.Description).
			Msg("Telegram API error")
		return err
	}

	log.Info().Msg("Webhook set successfully")
	return nil
}

func OnUpdateMessage(OnUpdateMessageBody OnUpdateMessageBody) {
	telegramChatID, err := strconv.Atoi(os.Getenv("TELEGRAM_CHAT_ID"))
	if err != nil {
		log.Error().Err(err).Msg("Error converting TELEGRAM_CHAT_ID to integer")
	}

	if strings.Contains(strings.ToLower(OnUpdateMessageBody.Message.Text), "/hello") {
		err := Send(telegramChatID, "Hello world!", false)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send message")
		}
		return
	}

	jsonStringByte, err := json.MarshalIndent(OnUpdateMessageBody, "", "  ")
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal OnUpdateMessageBody")
	}

	log.Print(string(jsonStringByte))
}
