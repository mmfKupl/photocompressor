# Билд исходников

```
GOOS=windows GOARCH=386 go build .
```

# Как пользоваться

Открыть консоль в папке с exe файлом, выполнить команду
```
./photocompressor.exe
```

Вы также можете передать данные через флаги командной строки:

- `-input`: путь к входной директории
- `-output`: путь к выходной директории
- `-bunch`: размер группы

```
./photocompressor.exe -input /path/to/input -output /path/to/output -bunch 10
```

## Описание программы

Программа `photocompressor` предназначена для сжатия фотографий и видеофайлов. Она обрабатывает файлы из указанной входной директории, сжимает их и сохраняет в указанную выходную директорию. Размер группы определяет количество файлов, обрабатываемых одновременно.