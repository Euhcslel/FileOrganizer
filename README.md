# FileOrganizer
Программа для сортировки файлов в папке по отдельным директориям.

Получает файлы исходной директории и определяет их расширение по магическим битам. Затем перемещает их в папки, созданные по категориям.

## Особенности

По умолчанию сортируется директория: `/home/user/Downloads`.

Файлы с неопределенным расширением отправляются в папку `others`.

## Установка

```bash
git clone https://github.com/Euhcslel/file-organizer
cd file-organizer
go build -o organizer
```

## Флаги

Указать собственную директорию можно с помощью ключа `-dir` (программа принимает путь только относительно домашней директории!):

```bash
./organizer -dir my_directory
```

Для того, чтобы увидеть как будет выглядеть дерево директории после изменений, можно использовать флаг `-predraw` (по умолчанию флаг имеет значение `false`):

```bash
./organizer -predraw
```
Пример вывода:
```
archive/
        archive.zip
binary/
        app.deb
pdf/
        pdf_file1.pdf
        pdf_file2.pdf
image/
        image1.jpg
        image2.jpeg
        image3.png
```

## Категории

| Папка | Что попадает |
|-------|--------------|
| `image/`   | Картинки (png, jpeg, gif, webp, bmp, ico) |
| `video/`   | Видео (mp4, webm, avi) |
| `audio/`   | Аудио (mp3, aiff, midi, wav, ogg) |
| `font/`    | Шрифты (ttf, otf, woff, woff2) |
| `text/`    | Текст (txt, html, xml) |
| `pdf/`     | PDF и PostScript |
| `archive/` | zip, gzip, rar |
| `binary/`  | AppImage, exe, wasm, неизвестные бинарники |
| `others/`  | Всё, что не подошло |

## Требования

Версия Go 1.25+