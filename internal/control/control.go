package control

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rodeorm/shortener/internal/repo"
)

/*
Сервер должен быть доступен по адресу: http://localhost:8080.
*/
func RouterStart(h *DecoratedHandler) error {

	r := mux.NewRouter()
	r.HandleFunc("/", h.RootHandler).Methods(http.MethodPost)
	r.HandleFunc("/{URL}", h.RootURLHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/shorten", h.ApiShortenHandler).Methods(http.MethodPost)
	r.HandleFunc("/", h.BadRequestHandler)

	err := http.ListenAndServe(":8080", r) // Не используем имя домена, всегда запускаем локально
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
		return err
	}

	return nil
}

type DecoratedHandler struct {
	DomainName string
	Storage    *repo.Storage
}
