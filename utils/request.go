package utils

import (
	"context"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

var (
	limiter *rate.Limiter
	client  http.Client
)

func InitRequest(limit time.Duration, timeout time.Duration) {
	// Rate Limiter
	r := rate.Every(60 * time.Second / limit)
	limiter = rate.NewLimiter(r, 1)

	// HTTP Client
	client = http.Client{Timeout: timeout * time.Second}
}

func SendRequest(query string) ([]byte, error) {
	err := limiter.Wait(context.Background())
	if err != nil {
		return []byte{}, err
	}

	req, err := http.NewRequest(http.MethodGet, query, nil)
	if err != nil {
		return []byte{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
