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

func writeFile(path string, data []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(data)
	return nil
}

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

func WriteStream(dir string, file string, url string) (int, error) {
	updated := 0

	// Create Subdirectory
	err := os.Mkdir(dir, 0o750)
	if err != nil && !os.IsExist(err) {
		return updated, err
	}

	// Try to read existing .strm file
	data, err := readFile(file)
	if err != nil {
		return updated, err
	}

	// Check if .strm already exists & has the correct stream
	if data != url {
		// In case .strm does not exist or has the incorrect stream, overwrite file
		err := writeFile(file, []byte(url+"\n"))
		if err != nil {
			return updated, err
		}
		updated = 1
	}

	return updated, nil
}

func WriteImage(dir string, file string, url string, enabled bool) (int, error) {
	updated := 0

	// If string is invalid http(s) link, do not update image
	if !strings.HasPrefix(url, "http") || !enabled {
		return 0, nil
	}

	// Open existing Image File, if it exists
	f, err := os.Open(file)
	f.Close()

	if err != nil && !os.IsNotExist(err) {
		return updated, err
	} else if os.IsNotExist(err) {
		// In case Image does not exist, download & create file
		image, err := SendRequest(url)
		if err != nil {
			// If image fetch fails, skip image creation, without error
			return updated, nil
		}

		err = writeFile(file, image)
		if err != nil {
			return updated, err
		}
		updated = 1
	}

	return updated, nil
}

func WriteNfo(dir string, file string, nfo string) (int, error) {
	updated := 0

	// Create Subdirectory
	err := os.Mkdir(dir, 0o750)
	if err != nil && !os.IsExist(err) {
		return updated, err
	}

	// Try to read existing .nfo file
	data, err := readFile(file)
	if err != nil {
		return updated, err
	}

	// Check if .nfo already exists & has the correct stream
	if data != nfo {
		// In case .nfo does not exist or has the incorrect stream, overwrite file
		err := writeFile(file, []byte(nfo))
		if err != nil {
			return updated, err
		}
		updated = 1
	}

	return updated, nil
}
