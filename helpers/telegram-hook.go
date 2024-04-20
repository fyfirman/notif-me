package helpers

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"
	"time"

	"notif-me/services/telegram"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

var wg sync.WaitGroup

type TelegramHook struct{}

func (t *TelegramHook) Run(
	e *zerolog.Event,
	level zerolog.Level,
	message string,
) {
	if level > zerolog.WarnLevel {
		wg.Add(1)
		go func() {
			telegramChatID, err := strconv.Atoi(os.Getenv("TELEGRAM_CHAT_ID"))
			if err != nil {
				log.Error().Msg(err.Error())
			}

			msg := "🚨 Error captured \n ```\n" + message + "\n```"
			_ = telegram.Send(telegramChatID, msg, false)
			wg.Done()
		}()
	}
}

type errWithStackTrace struct {
	Err string `json:"error"`
	// Stacktrace *raven.Stacktrace `json:"stacktrace"` -- disabled 3rd party
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type logEvent struct {
	Level     string            `json:"level"`
	Msg       string            `json:"message"`
	Err       errWithStackTrace `json:"error"`
	Time      time.Time         `json:"time"`
	Status    int               `json:"status"`
	UserAgent string            `json:"user_agent"`
	Method    string            `json:"method"`
	URL       string            `json:"url"`
	IP        string            `json:"ip"`
}

var errSkipEvent = errors.New("skip")

// unmarshal only if the level is error.
func (l *logEvent) UnmarshalJSON(data []byte) error {
	res := gjson.Get(string(data), "level")
	if !res.Exists() || res.String() != "error" {
		return errSkipEvent
	}

	type event logEvent
	return json.Unmarshal(data, (*event)(l))
}
