package main

import (
	"flag"
	"os"

	"github.com/rodeorm/shortener/internal/api"
	"github.com/rodeorm/shortener/internal/repo"
)

/*
Сonfig выполняет первоначальную конфигурацию.

Приоритет параметров сервера должен быть таким:
Если указана переменная окружения, то используется она.
Если нет переменной окружения, но есть аргумент командной строки (флаг), то используется он.
Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.
*/
//config выполняет первоначальную конфигурацию
func config() *api.DecoratedHandler {

	flag.Parse()

	/*
		os.Setenv("SERVER_ADDRESS", "localhost:8080")
		os.Setenv("BASE_URL", "http://tiny")
		os.Setenv("FILE_STORAGE_PATH", "D:/file.txt")
		os.Setenv("DATABASE_DSN", "postgres://app:qqqQQQ123@localhost:5433/shortener?sslmode=disable")
	*/

	var serverAddress, baseURL, fileStoragePath, databaseConnectionString string

	//Адрес запуска HTTP-сервера
	if *a == "" {
		serverAddress = os.Getenv("SERVER_ADDRESS")
		if serverAddress == "" {
			serverAddress = "localhost:8080"
		}
	} else {
		serverAddress = *a
	}

	//Базовый адрес результирующего сокращённого URL
	if *b == "" {
		baseURL = os.Getenv("BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:8080"
		}
	} else {
		baseURL = *b
	}

	//Путь до файла
	if *f == "" {
		fileStoragePath = os.Getenv("FILE_STORAGE_PATH")
	} else {
		fileStoragePath = *f
	}

	//Строка подключения к БД
	if *d == "" {
		databaseConnectionString = os.Getenv("DATABASE_DSN")
	} else {
		databaseConnectionString = *d
	}

	return &api.DecoratedHandler{ServerAddress: serverAddress, Storage: repo.NewStorage(fileStoragePath, databaseConnectionString), BaseURL: baseURL}

}
