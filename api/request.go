package api

import (
	"io"
	"net/http"
)

func SendRequest(query string) ([]byte, error) {
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
