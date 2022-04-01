package control

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

/*
Сервер должен быть доступен по адресу: http://localhost:8080.
*/
func RouterStart(handler *DecoratedHandler) error {

	r := mux.NewRouter()
	r.Methods("GET", "POST").Handler(http.HandlerFunc(handler.returnURLHandler))

	err := http.ListenAndServe(":8080", r) // Не используем имя домена, всегда запускаем локально
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return err
	}

	return nil
}
