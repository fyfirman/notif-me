package services

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

func checkForUpdates(url string, noChapterIdentifier string) {
	log.Println("Checking for updates, url :", url)

	// Send a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	// Read the HTML content of the page
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	// Check for the presence of the new release
	if strings.Contains(string(body), noChapterIdentifier) {
		log.Println("No new chapter yet")
		telegram.Send("No new chapter yet", true)
	} else {
		log.Println("New chapter released!")
		telegram.Send("New chapter released. Go to link: "+url, false)
	}
}

func Start(env *env.Env) {
	s := gocron.NewScheduler(time.UTC)

	s.Every(15).Minutes().Do(func() {
		rawUrl := "https://komikcast.site/chapter/jujutsu-kaisen-chapter-*-bahasa-indonesia/"
		noChapterIdentifier := "<title>Halaman tidak di temukan - Komikcast</title>"

		url := helpers.ReplaceWildcard(rawUrl, 2, 221)
		checkForUpdates(url, noChapterIdentifier)
	})
	s.StartAsync()
}
