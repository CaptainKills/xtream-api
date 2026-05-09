package main

import (
	"log"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/CaptainKills/xtream-api/api"
)

var (
	Version  = "v2.1.0"
	Commit   = "unknown"
	Branch   = "unknown"
	Date     = "unknown"
	Modified = "false"
)

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				Commit = setting.Value
			case "vcs.time":
				Date = setting.Value
			case "vcs.modified":
				Modified = setting.Value
			}
		}
	}
}

func main() {
	log.Printf("Starting xtream-api %s", Version)
	log.Printf("  Branch:     %s", Branch)
	log.Printf("  Commit:     %s", Commit)
	log.Printf("  Build Time: %s", Date)
	log.Printf("  Modified:   %s", Modified)
	log.Printf("  Go Version: %s", runtime.Version())
	log.Printf("  OS/Arch:    %s/%s", runtime.GOOS, runtime.GOARCH)

	// Environment Variables
	url, username, password := GetCredentials()
	config := GetApplicationConfig()
	options := GetXtreamOptions()

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
	if time.Now().Before(config.LaunchTime) {
		launch := config.LaunchTime.Format(time.DateTime)
		now := time.Now().Format(time.DateTime)

		log.Printf("[INFO] Next run scheduled at: %s (Current time: %s)\n", launch, now)
		time.Sleep(time.Until(config.LaunchTime))
	}

	for {
		start := time.Now()

		// Run Programs
		err = livestreams.Run(client, config)
		if err != nil {
			log.Printf("[ERROR] (%s) Unable to run program: %v\n", livestreams.label, err)
		}

		err = movies.Run(client, config)
		if err != nil {
			log.Printf("[ERROR] (%s) Unable to run program: %v\n", movies.label, err)
		}

		err = series.Run(client, config)
		if err != nil {
			log.Printf("[ERROR] (%s) Unable to run program: %v\n", series.label, err)
		}

		// Wait Until Next Run
		diff := time.Since(start).Round(time.Millisecond)
		next := start.Add(24 * time.Hour)

		log.Printf("[INFO] Run Duration: %s. Next run scheduled at: %s\n", diff.String(), next.Format(time.DateTime))
		time.Sleep(time.Until(next))
	}
}
