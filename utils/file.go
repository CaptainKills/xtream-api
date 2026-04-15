package utils

import (
	"log"
	"os"
	"strings"
)

func GetImageExtension(url string) string {
	var extension string

	if strings.Contains(url, "png") || strings.Contains(url, "PNG") {
		extension = ".png"
	} else if strings.Contains(url, "jpg") || strings.Contains(url, "JPG") {
		extension = ".jpg"
	} else if strings.Contains(url, "jpeg") || strings.Contains(url, "JPEG") {
		extension = ".jpeg"
	} else if strings.Contains(url, "webp") || strings.Contains(url, "WEBP") {
		extension = ".webp"
	} else if strings.Contains(url, "avi") || strings.Contains(url, "AVI") {
		extension = ".avif"
	} else {
		extension = ".jpg"
	}

	return extension
}

func ImageExists(file string) bool {
	f, err := os.Open(file)

	if err != nil && os.IsNotExist(err) {
		return false
	} else if err != nil {
		log.Printf("[WARNING] Failed to check if image exists: %v\n", err)
		return false
	}

	defer f.Close()
	return true
}

func WriteFile(file string, data string) (int, error) {
	updated := 0

	// Try to read existing file
	f, err := os.ReadFile(file)
	if err != nil {
		return updated, err
	}

	// Check if file already exists & has the correct stream
	if string(f) != data {
		// In case file does not exist or has the incorrect data, overwrite file
		err := os.WriteFile(file, []byte(data), 0o644)
		if err != nil {
			return updated, err
		}
		updated = 1
	}

	return updated, nil
}

func WriteImage(file string, image []byte) (int, error) {
	err := os.WriteFile(file, image, 0o644)
	if err != nil {
		return 0, err
	}

	return 1, nil
}
