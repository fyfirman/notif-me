package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func log(s ...string) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}
	currentTime := time.Now().In(loc)

	message := currentTime.Format("2006-01-02 15:04:05") + " GMT+7 " + strings.Join(s, " ")
	fmt.Println(message)
}

func checkForUpdates(url string, noChapterIdentifier string) {
	log("Checking for updates, url :", url)

	// Send a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Read the HTML content of the page
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Check for the presence of the new release
	if strings.Contains(string(body), noChapterIdentifier) {
		log("No new chapter yet")
		sendTelegramMessage("No new chapter yet", true)
	} else {
		log("New chapter released!")
		sendTelegramMessage("New chapter released. Go to link:"+url, false)
	}
}

func main() {
	godotenv.Load(".env")

	url := "https://komikcast.site/chapter/jujutsu-kaisen-chapter-221-bahasa-indonesia/"
	noChapterIdentifier := "<title>Halaman tidak di temukan - Komikcast</title>"

	sendTelegramMessage("🚀 NotifMe started with url: "+url, false)

	checkForUpdates(url, noChapterIdentifier)
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		checkForUpdates(url, noChapterIdentifier)
	}
}
