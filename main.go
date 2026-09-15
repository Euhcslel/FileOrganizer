package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Обобщенные расширения файлов
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
	// Установка флагов
	dir := flag.String("dir", "Downloads", "sortable directory")
	preDraw := flag.Bool("predraw", false, "predraw scheme")
	flag.Parse()

	// Определение домашней директории
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("ошибка при получении домашней директории пользователя: %s", err.Error())
		return
	}
	path := filepath.Join(homeDir, *dir)

	// Чтение папки
	files, err := os.ReadDir(path)
	if err != nil {
		log.Printf("ошибка при чтении директории: %s", err.Error())
		return
	}

	// Переменная для хранения новой структуры папки
	var structure = map[string][]string{}

	// Цикл по файлам в папке
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Открытие файла (для чтения)
		oldPath := filepath.Join(path, file.Name())
		openedFile, err := os.Open(oldPath)
		if err != nil {
			log.Printf("ошибка при открытии файла: %s", err.Error())
			continue
		}
		
		// Чтение файла в буфер
		buffer := make([]byte, 512)
		n, err := openedFile.Read(buffer)
		if err != nil {
			log.Printf("ошибка при чтении в буфер: %s", err.Error())
			openedFile.Close()
			continue
		}
		openedFile.Close()

		// Определение расширения файла
		ext := http.DetectContentType(buffer[:n])
		parts := strings.Split(ext, "/")
		ext = parts[0]

		// Определение категории в зависимости от расширения
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

		// Заполнение итоговой структуры папок
		if *preDraw {
			structure[category] = append(structure[category], file.Name())
			continue
		}

		// Создание папки для хранения файлов определенной категории
		categoryPath := filepath.Join(path, category)
		if err := os.MkdirAll(categoryPath, 0755); err != nil {
			log.Printf("ошибка при создании новой папки: %s", err.Error())
		}

		// Перемещение файла в указанную папку
		newPath := filepath.Join(categoryPath, file.Name())
		if err := os.Rename(oldPath, newPath); err != nil {
			log.Printf("ошибка при перемещении файла %s: %s", file.Name(), err.Error())
		}
	}

	if *preDraw {
		for category, files := range structure {
			fmt.Printf("\n%s/", category)
	
			for _, file := range files {
				fmt.Printf("\n\t%s", file)
			}
		}
		fmt.Println()
	}
}