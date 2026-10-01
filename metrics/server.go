package metrics

import (
	"log"
	"net/http"
)

func ListenAndServe() {
	go func() {
		log.Printf("[INFO] Starting HTTP Server on :8080")

		http.Handle("/metrics", metricsHandler())
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
}
