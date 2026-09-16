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
	"application": {
		"pdf":                "pdf",
		"postscript":         "pdf",
		"x-gzip":             "archive",
		"zip":                "archive",
		"x-rar-compressed":   "archive",
		"vnd.ms-fontobject":  "font",
		"wasm":               "binary",
		"octet-stream":       "binary",
	},
	"audio": {},
	"font":  {},
	"image": {},
	"text":  {},
	"video": {},
}

func main() {
	// Установка флагов
	dir := flag.String("dir", "Downloads", "sortable directory")
	preDraw := flag.Bool("predraw", false, "predraw scheme")
	flag.Parse()

	// Получение итогового пути до сортируемой папки
	path, err := getFinalPath(*dir)
	if err != nil {
		log.Printf("ошибка при получении домашней директории пользователя: %s", err.Error())
		return
	}

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

		oldPath := filepath.Join(path, file.Name())
		ext, err := getFileExtension(oldPath)
		if err != nil {
			log.Println(err)
			continue
		}

		category := getFileCategory(ext)

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

	// Вывод структуры папки
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

// getFileExtension возвращает расширение файла, определяя его по магическим битам файла.
func getFileExtension(oldPath string) (string, error) {
	openedFile, err := os.Open(oldPath)
	if err != nil {
		return "", fmt.Errorf("ошибка при открытии файла: %w", err)
	}

	// Чтение файла в буфер
	buffer := make([]byte, 512)
	n, err := openedFile.Read(buffer)
	if err != nil {
		openedFile.Close()
		return "", fmt.Errorf("ошибка при чтении в буфер: %w", err)
	}
	openedFile.Close()

	// Определение расширения файла
	return http.DetectContentType(buffer[:n]), nil
}

// getFinalPath возвращает путь до сортируемой папки. Ищет папку в домашней директории пользователя.
func getFinalPath(dir string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, dir), nil
}

// getFileCategory возвращает обобщенную категорию файла.
// Если расширение файла начинает с application/, то разбор идет более детальный.
func getFileCategory(ext string) string {
	parts := strings.Split(ext, "/")
	ext = parts[0]

	subcategories, ok := categories[ext]
	if !ok {
		return "others"
	} else if ext == "application" {
		subcategory, ok := subcategories[parts[1]]
		if !ok {
			return "others"
		} else {
			return subcategory
		}
	} else {
		return ext
	}
}
