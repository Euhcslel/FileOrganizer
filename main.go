package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var categories = map[string]string{
	"image/jpeg": "image",
	"image/png": "image",
	"application/octet-stream": "application",
	"application/pdf": "pdf",
}

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {

	}
	defaultDir := "Downloads"

	path := filepath.Join(homeDir, defaultDir)
	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("ошибка при чтении директории: %s", err.Error())
		return
	}
	for _, file := range files {
		fileName := file.Name()
		oldPath := filepath.Join(path, file.Name())
		file, err := os.Open(oldPath)
		if err != nil {
			log.Printf("ошибка при открытии файла: %s", err.Error())
			return
		}
		
		defer file.Close()

		buffer := make([]byte, 512)
		file.Read(buffer)

		ext := http.DetectContentType(buffer)

		categoryPath := filepath.Join(path, categories[ext])
		os.Mkdir(categoryPath, 0755)

		newPath := filepath.Join(categoryPath, fileName)
		log.Println(newPath, fileName)
		if err := os.Rename(oldPath, newPath); err != nil {
			log.Printf("ошибка при перемещении файла %s: %s", file.Name(), err.Error())
		}
	}
}