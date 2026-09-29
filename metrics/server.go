package metrics

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handlerStats(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.Marshal(statistics)
	if err != nil {
		msg := fmt.Errorf("<b>Invalid Statistics:</b> %v\n", err)
		fmt.Fprint(w, msg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bytes))
}

func ListenAndServe() {
	go func() {
		log.Printf("[INFO] Starting HTTP Server on :8080")

		http.HandleFunc("/metrics", handlerStats)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
}
