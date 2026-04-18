package main

import (
	"log"
	"time"

	"github.com/CaptainKills/xtream-api/api"
)

func main() {
	// Environment Variables
	url, username, password := GetEnvironmentCredentials()
	options := GetEnvironmentOptions()

	if url == "" || username == "" || password == "" {
		log.Fatalln("[ERROR] Missing Environment Variables! Exiting Program...")
	}

	// Xtream Client
	client := api.NewClient(url, username, password, options)
	_, err := client.GetAccountInfo()
	if err != nil {
		log.Fatalf("[ERROR] Authentication Failed: %v\n", err)
	}
	log.Printf("[INFO] Authentication Successful: %s\n", url)

	// Launch Time Delay
	if time.Now().Before(options.LaunchTime) {
		launch := options.LaunchTime.Format(time.DateTime)
		now := time.Now().Format(time.DateTime)

		log.Printf("[INFO] Next run scheduled at: %s (Current time: %s)\n", launch, now)
		time.Sleep(time.Until(options.LaunchTime))
	}

	for {
		start := time.Now()

		// Run Programs
		err = livestreams.Run(client)
		if err != nil {
			log.Printf("[ERROR] Unable to run program %s: %v\n", livestreams.label, err)
		}

		err = movies.Run(client)
		if err != nil {
			log.Printf("[ERROR] Unable to run program %s: %v\n", movies.label, err)
		}

		err = series.Run(client)
		if err != nil {
			log.Printf("[ERROR] Unable to run program %s: %v\n", series.label, err)
		}

		// Wait Until Next Run
		diff := time.Since(start).Round(time.Millisecond)
		next := start.Add(24 * time.Hour)

		log.Printf("[INFO] Run Duration: %s. Next run scheduled at: %s\n", diff.String(), next.Format(time.DateTime))
		time.Sleep(time.Until(next))
	}
}
