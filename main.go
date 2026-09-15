package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var categories = map[string]map[string]string{
	"image": {},
	"application": {
		"pdf": "pdf",
		"zip": "archive",
		"octet-stream": "binary",
	},
	"text": {},
	"video": {},
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("ошибка при получении домашней директории пользователя: %s", err.Error())
		return
	}
	defaultDir := "Downloads"

	path := filepath.Join(homeDir, defaultDir)
	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("ошибка при чтении директории: %s", err.Error())
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		oldPath := filepath.Join(path, file.Name())
		openedFile, err := os.Open(oldPath)
		if err != nil {
			log.Printf("ошибка при открытии файла: %s", err.Error())
			continue
		}
		
		
		buffer := make([]byte, 512)
		n, err := openedFile.Read(buffer)
		if err != nil {
			log.Printf("ошибка при чтении в буфер: %s", err.Error())
			openedFile.Close()
			continue
		}
		openedFile.Close()

		ext := http.DetectContentType(buffer[:n])
		parts := strings.Split(ext, "/")
		ext = parts[0]

		var category string
		values, ok := categories[ext]
		if !ok {
			category = "others"
		} else if ext == "application" {
			if category, ok = values[parts[1]]; !ok {
				category = "others"
			}
		} else {
			category = ext
		}

		categoryPath := filepath.Join(path, category)
		if err := os.MkdirAll(categoryPath, 0755); err != nil {
			log.Printf("ошибка при создании новой папки: %s", err.Error())
		}

		newPath := filepath.Join(categoryPath, file.Name())
		if err := os.Rename(oldPath, newPath); err != nil {
			log.Printf("ошибка при перемещении файла %s: %s", file.Name(), err.Error())
		}
	}
}