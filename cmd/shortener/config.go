package main

import (
	"flag"

	env "github.com/caarlos0/env/v6"

	"github.com/rodeorm/shortener/internal/api"
	"github.com/rodeorm/shortener/internal/repo"
)

type Config struct {
	serverAddress   string `env:"SERVER_ADDRESS"`    //Адрес запуска HTTP-сервера
	baseURL         string `env:"BASE_URL"`          //Базовый адрес результирующего сокращённого URL
	fileStoragePath string `env:"FILE_STORAGE_PATH"` //Путь до файла
	databaseDSN     string `env:"DATABASE_DSN"`      //Строка подключения к БД
}

/*
Сonfig выполняет первоначальную конфигурацию.

Приоритет параметров сервера должен быть таким:
Если указана переменная окружения, то используется она.
Если нет переменной окружения, но есть аргумент командной строки (флаг), то используется он.
Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.
*/
//config выполняет первоначальную конфигурацию
func config() *api.DecoratedHandler {

	/*
		os.Setenv("SERVER_ADDRESS", "localhost:8080")
		os.Setenv("BASE_URL", "http://tiny")
		os.Setenv("FILE_STORAGE_PATH", "D:/file.txt")
		os.Setenv("DATABASE_DSN", "postgres://app:qqqQQQ123@localhost:5433/shortener?sslmode=disable")
	*/
	var cfg Config

	flag.Parse()
	env.Parse(&cfg)

	if *a != "" {
		cfg.serverAddress = *a
	}

	if *b != "" {
		cfg.baseURL = *b
	}

	if *f != "" {
		cfg.fileStoragePath = *f
	}

	if *d != "" {
		cfg.databaseDSN = *d
	}

	return &api.DecoratedHandler{ServerAddress: cfg.serverAddress, Storage: repo.NewStorage(cfg.fileStoragePath, cfg.databaseDSN), BaseURL: cfg.baseURL}
}
