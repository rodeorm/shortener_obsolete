package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rodeorm/shortener/internal/control"
	"github.com/rodeorm/shortener/internal/repo"
)

//Config выполняет первоначальную конфигурацию
func config() *control.DecoratedHandler {
	flag.Parse()
	//	os.Setenv("SERVER_ADDRESS", "localhost:8080")
	//	os.Setenv("BASE_URL", "http://tiny")
	//  os.Setenv("FILE_STORAGE_PATH", "D:/file.txt  nn")

	var sa, bu, fsp string

	//Адрес запуска HTTP-сервера с помощью переменной SERVER_ADDRESS
	if *a == "" {
		sa = os.Getenv("SERVER_ADDRESS")
		if sa == "" {
			//		fmt.Println("Не найдена переменная среды SERVER_ADDRESS")
			sa = "localhost:8080"
		}
	}

	//Базовый адрес результирующего сокращённого URL с помощью переменной BASE_URL.
	if *b == "" {
		bu = os.Getenv("BASE_URL")
		if bu == "" {
			//	fmt.Println("Не найдена переменная среды BASE_URL")
			bu = "http://localhost:8080"
		}
	}

	if *f == "" {
		fsp = os.Getenv("FILE_STORAGE_PATH")
	}
	//Путь до файла должен передаваться в переменной окружения FILE_STORAGE_PATH.

	fmt.Println("Адрес запуска http сервера: ", sa)
	fmt.Println("Базовый адрес результирующего сокращённого URL: ", bu)
	fmt.Println("Путь до файла: ", fsp)

	return &control.DecoratedHandler{ServerAddress: sa, Storage: repo.NewStorage(fsp), BaseURL: bu}
}
