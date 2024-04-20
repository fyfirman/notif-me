package main

import (
	"notif-me/env"
	"notif-me/helpers"
	cronService "notif-me/services/cron"
	"notif-me/services/telegram"
	"os"
	"strconv"

	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	appVersion := "0.3.0"

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	log.Logger = log.With().Str("app_version", appVersion).Logger()
	log.Logger = log.Hook(&helpers.TelegramHook{})

	log.Info().Msg("🚀 Starting NotifMe v" + appVersion)

	godotenv.Load(".env")

	db, err := ConnectDB()

	if err != nil {
		panic(err.Error())
	}

	telegramChatID, err := strconv.Atoi(os.Getenv("TELEGRAM_CHAT_ID"))

	if err != nil {
		log.Error().Msg(err.Error())
	}

	telegram.Send(telegramChatID, "🚀 NotifMe v"+appVersion+" has started", false)

	env := &env.Env{Db: db}

	cronService.Start(env)

	router := mux.NewRouter()

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	router.HandleFunc("/api/telegram", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			URL string `json:"url"`
		}

		err := json.NewDecoder(r.Body).Decode(&body)

		if err != nil {
			log.Error().Msg(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]bool{"ok": false})
			return
		}

		telegram.SetWebhook(body.URL)

		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	router.HandleFunc("/api/telegram/receive", func(w http.ResponseWriter, r *http.Request) {
		var onUpdateMessageBody telegram.OnUpdateMessageBody
		err := json.NewDecoder(r.Body).Decode(&onUpdateMessageBody)

		if err != nil {
			log.Error().Msg(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]bool{"ok": false})
			return
		}

		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	http.Handle("/", router)
	srv := &http.Server{
		Handler: router,
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Info().Msg("Application served at :8080")
	log.Error().Msg(srv.ListenAndServe().Error())
}
