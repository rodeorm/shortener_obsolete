package main

import (
	"flag"
	"os"

	"github.com/rodeorm/shortener/internal/control"
	"github.com/rodeorm/shortener/internal/repo"
)

//config выполняет первоначальную конфигурацию
func config() *control.DecoratedHandler {
	flag.Parse()
	//	os.Setenv("SERVER_ADDRESS", "localhost:8080")
	//	os.Setenv("BASE_URL", "http://tiny")
	//  os.Setenv("FILE_STORAGE_PATH", "D:/file.txt  nn")
	var sa, bu, fsp string

	// fmt.Println("flags", *a, *b, *f)
	//Адрес запуска HTTP-сервера
	if *a == "" {
		sa = os.Getenv("SERVER_ADDRESS")
		if sa == "" {
			sa = "localhost:8080"
		}
	} else {
		sa = *a
	}

	//Базовый адрес результирующего сокращённого URL
	if *b == "" {
		bu = os.Getenv("BASE_URL")
		if bu == "" {
			bu = "http://localhost:8080"
		}
	} else {
		bu = *b
	}

	//Путь до файла
	if *f == "" {
		fsp = os.Getenv("FILE_STORAGE_PATH")
	} else {
		fsp = *f
	}

	return &control.DecoratedHandler{ServerAddress: sa, Storage: repo.NewStorage(fsp), BaseURL: bu}
}
