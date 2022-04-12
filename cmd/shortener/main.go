package main

import (
	"fmt"
	"os"

	"github.com/rodeorm/shortener/internal/control"
	"github.com/rodeorm/shortener/internal/repo"
)

/*
Сервис для сокращения длинных URL. Требования:
Сервер должен быть доступен по адресу: http://localhost:8080.
Сервер должен предоставлять два эндпоинта: POST / и GET /{id}.
Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400.
*/
func main() {
	handler := config()
	control.RouterStart(handler)
}

//Config выполняет первоначальную конфигурацию сервиса и возвращает - имя домена, соответствие  ключа к оригинальному URL
func config() *control.DecoratedHandler {

	//os.Setenv("SERVER_ADDRESS", "http://localhost:8080")
	//os.Setenv("BASE_URL", "http://tiny")

	//Адрес запуска HTTP-сервера с помощью переменной SERVER_ADDRESS
	sa := os.Getenv("SERVER_ADDRESS")
	if sa == "" {
		fmt.Println("Не найдена переменная среды SERVER_ADDRESS")
		sa = "http://localhost:8080"
	}
	//Базовый адрес результирующего сокращённого URL с помощью переменной BASE_URL.
	bu := os.Getenv("BASE_URL")
	if bu == "" {
		fmt.Println("Не найдена переменная среды BASE_URL")
		bu = "http://localhost:8080"
	}
	fmt.Println("Адрес запуска http сервера: ", sa, ". Базовый адрес результирующего url: ", bu)
	return &control.DecoratedHandler{ServerAddress: sa, Storage: repo.NewStorage(), BaseURL: bu}
}
