package api

import (
	"bufio"
	"os"
	"strings"
)

func readFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		return "", err
	}

	var data string
	for scanner.Scan() {
		data = data + scanner.Text()
	}

	return data, nil
}

func writeFile(path string, data string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString(data + "\n")
	return nil
}

func WriteStream(dir string, file string, url string) (int, error) {
	updated := 0

	// Create Subdirectory
	err := os.Mkdir(dir, 0o750)
	if err != nil && !os.IsExist(err) {
		return 0, err
	}

	// Check if .strm already exists & has the correct stream
	data, err := readFile(file)
	if err != nil {
		return 0, err
	}

	// In case .strm does not exist or has the incorrect stream, overwrite file
	if data != url {
		err := writeFile(file, url)
		if err != nil {
			return 0, err
		}
		updated = 1
	}

	// Download Cover Image
	// if strings.HasPrefix(icon, "http") {
	// 	image, err := c.sendRequest(icon)
	// 	if err != nil {
	// 		// log.Printf("[WARNING] Unable to fetch cover image: %v\n", err)
	// 		skipped_images++
	// 	}
	// 	WriteImage(pathDirectory, icon, image)
	// }

	return updated, nil
}

func WriteImage(path string, url string, data []byte) error {
	var filename string

	if strings.Contains(url, "png") || strings.Contains(url, "PNG") {
		filename = "cover.png"
	} else if strings.Contains(url, "jpg") || strings.Contains(url, "JPG") {
		filename = "cover.jpg"
	} else if strings.Contains(url, "jpeg") || strings.Contains(url, "JPEG") {
		filename = "cover.jpeg"
	} else if strings.Contains(url, "webp") || strings.Contains(url, "WEBP") {
		filename = "cover.webp"
	} else if strings.Contains(url, "avi") || strings.Contains(url, "AVI") {
		filename = "cover.avif"
	} else {
		filename = "cover.jpg"
	}

	file, err := os.Create(path + "/" + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(data)
	return nil
}
