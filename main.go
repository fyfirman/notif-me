package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
		sendTelegramMessage("No new chapter yet", true)
	} else {
		log.Println("New chapter released!")
		sendTelegramMessage("New chapter released. Go to link:"+url, false)
	}
}

func main() {
	godotenv.Load(".env")

	url := "https://komikcast.site/chapter/jujutsu-kaisen-chapter-221-bahasa-indonesia/"
	noChapterIdentifier := "<title>Halaman tidak di temukan - Komikcast</title>"

	sendTelegramMessage("🚀 NotifMe v0.0.1 started with url: "+url, false)

	s := gocron.NewScheduler(time.UTC)
	s.Every(15).Minutes().Do(func() {
		checkForUpdates(url, noChapterIdentifier)
	})
	s.StartAsync()

	router := mux.NewRouter()

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
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

	log.Println("Application served at :8080")
	log.Fatal(srv.ListenAndServe())
}
