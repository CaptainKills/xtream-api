package main

import (
	"log"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/CaptainKills/xtream-api/api"
	"github.com/CaptainKills/xtream-api/metrics"
)

var (
	Version  = "v2.2.0"
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

	metrics.UpdateApp(metrics.AppStatistics{Version: Version, Commit: Commit, Branch: Branch, Date: Date, Modified: Modified})
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
	options := GetXtreamOptions()

	if url == "" || username == "" || password == "" {
		log.Fatalln("[ERROR] Missing Environment Variables! Exiting Program...")
	}

	metrics.UpdateOptions(options)
	metrics.ListenAndServe()

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
			log.Printf("[ERROR] (%s) Unable to run program: %v\n", livestreams.label, err)
		}

		err = movies.Run(client)
		if err != nil {
			log.Printf("[ERROR] (%s) Unable to run program: %v\n", movies.label, err)
		}

		err = series.Run(client)
		if err != nil {
			log.Printf("[ERROR] (%s) Unable to run program: %v\n", series.label, err)
		}

		// Wait Until Next Run
		end := time.Now()
		diff := time.Since(start).Round(time.Millisecond)
		next := start.Add(options.LaunchPeriod * time.Hour)
		metrics.UpdateTime(metrics.TimeStatistics{Start: start, End: end, Duration: diff, Next: next})

		log.Printf("[INFO] Run Duration: %s. Next run scheduled at: %s\n", diff.String(), next.Format(time.DateTime))
		time.Sleep(time.Until(next))
	}
}
