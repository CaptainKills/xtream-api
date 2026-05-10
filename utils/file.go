package utils

import (
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

func WriteFile(file string, data []byte) (int, error) {
	err := os.WriteFile(file, data, 0o644)
	if err != nil {
		return 0, err
	}

	return 1, nil
}
