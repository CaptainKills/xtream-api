package api

import (
	"context"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// limiter ensures that there are at most 1000 requests per 60 seconds.
var limiter = rate.NewLimiter(rate.Every(60*time.Second/1000), 1)

func SendRequest(query string) ([]byte, error) {
	err := limiter.Wait(context.Background())
	if err != nil {
		return []byte{}, err
	}

	req, err := http.NewRequest(http.MethodGet, query, nil)
	if err != nil {
		return []byte{}, err
	}

	resp, err := http.DefaultClient.Do(req)
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
