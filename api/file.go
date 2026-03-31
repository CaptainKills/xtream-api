package api

import (
	"bufio"
	"os"
)

func ReadFile(path string) (string, error) {
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

func WriteFile(path string, data string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString(data + "\n")
	return nil
}
