package cron

import (
	"io/ioutil"
	"log"
	"net/http"
	"notif-me/env"
	"notif-me/helpers"
	"notif-me/services/telegram"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
)

func checkForUpdates(url string, noChapterIdentifier string) (bool, error) {
	log.Println("Checking for updates, url :", url)

	// Send a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("No new chapter yet")
		return false, nil
	}

	// Read the HTML content of the page
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return false, err
	}

	// Check for the presence of the new release
	if strings.Contains(string(body), noChapterIdentifier) {
		log.Println("No new chapter yet")
		return false, nil
	}

	log.Println("New chapter released!")

	return true, nil
}

func Start(env *env.Env) {
	s := gocron.NewScheduler(time.UTC)

	s.Every(15).Minutes().Do(func() {
		log.Println("Cron every 15 minutes starting...")

		res, err := GetAll(env)

		if err != nil {
			log.Println("Error : " + err.Error())
			return
		}

		for _, mangaUpdate := range res {
			url := helpers.ReplaceWildcard(mangaUpdate.RawURL, 2, mangaUpdate.LastChapter)

			hasNewChapter, err := checkForUpdates(url, mangaUpdate.NegativeIdentifier)

			if hasNewChapter {
				telegram.Send(mangaUpdate.ChatID, "New chapter released. Go to link: "+url, false)
			}

			if err != nil {
				log.Println("ERROR: " + err.Error())
				continue
			}

			var payload map[string]interface{}

			if !hasNewChapter {
				payload = map[string]interface{}{
					"id":              mangaUpdate.ID,
					"last_checked_at": time.Now(),
					"updated_at":      time.Now(),
				}
				err = UpdateById(env, mangaUpdate.ID, payload)
				if err != nil {
					log.Println("ERROR: " + err.Error())
				}
				continue
			}

			payload = map[string]interface{}{
				"id":              mangaUpdate.ID,
				"last_chapter":    mangaUpdate.LastChapter + 1,
				"last_checked_at": time.Now(),
				"updated_at":      time.Now(),
			}
			err = UpdateById(env, mangaUpdate.ID, payload)

			if err != nil {
				log.Println("ERROR: " + err.Error())
				continue
			}
		}
	})

	s.StartAsync()
}
