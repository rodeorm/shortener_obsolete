package main

import (
	"flag"
	"os"

	"github.com/rodeorm/shortener/internal/control"
	"github.com/rodeorm/shortener/internal/repo"
)

/*
Сonfig выполняет первоначальную конфигурацию.

Приоритет параметров сервера должен быть таким:
Если указана переменная окружения, то используется она.
Если нет переменной окружения, но есть аргумент командной строки (флаг), то используется он.
Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.
*/
func config() *control.DecoratedHandler {
	flag.Parse()

	/*
		serverAddress := "localhost:8080"                                                               //Адрес запуска HTTP-сервера
		baseURL := "http://localhost:8080"                                                              //Базовый URL
		fileStoragePath := "D:/file.txt"                                                                //Путь до файла
		databaseConnectionString := "postgres://app:qqqQQQ123@localhost:5433/shortener?sslmode=disable" //Строка подключения к БД
	*/
	serverAddress := "localhost:8080"  //Адрес запуска HTTP-сервера
	baseURL := "http://localhost:8080" //Базовый URL
	fileStoragePath := ""              //Путь до файла
	databaseConnectionString := ""     //Строка подключения к БД

	if *a != "" {
		serverAddress = *a
	} else if os.Getenv("SERVER_ADDRESS") == "" {
		serverAddress = os.Getenv("SERVER_ADDRESS")
	}

	if *b != "" {
		baseURL = *b
	} else if os.Getenv("BASE_URL") == "" {
		baseURL = os.Getenv("BASE_URL")
	}

	if *f != "" {
		fileStoragePath = *f
	} else if os.Getenv("FILE_STORAGE_PATH") == "" {
		fileStoragePath = os.Getenv("FILE_STORAGE_PATH")
	}

	if *d != "" {
		databaseConnectionString = os.Getenv("DATABASE_DSN")
	} else if os.Getenv("DATABASE_DSN") != "" {
		databaseConnectionString = *d
	}

	/*
		if os.Getenv("SERVER_ADDRESS") != "" {
			serverAddress = os.Getenv("SERVER_ADDRESS")
		} else if *a != "" {
			serverAddress = *a
		}

		if os.Getenv("BASE_URL") != "" {
			baseURL = os.Getenv("BASE_URL")
		} else if *b != "" {
			baseURL = *b
		}

		if os.Getenv("FILE_STORAGE_PATH") != "" {
			fileStoragePath = os.Getenv("FILE_STORAGE_PATH")
		} else if *f != "" {
			fileStoragePath = *f
		}

		if os.Getenv("DATABASE_DSN") != "" {
			databaseConnectionString = os.Getenv("DATABASE_DSN")
		} else if *d != "" {
			databaseConnectionString = *d
		}
	*/
	return &control.DecoratedHandler{ServerAddress: serverAddress, Storage: repo.NewStorage(fileStoragePath, databaseConnectionString), BaseURL: baseURL}
}
