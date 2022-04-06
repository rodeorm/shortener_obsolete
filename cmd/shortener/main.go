package main

import (
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
	return &control.DecoratedHandler{DomainName: "http://localhost:8080", Storage: repo.NewStorage()}
}
