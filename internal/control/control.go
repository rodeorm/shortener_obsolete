package control

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rodeorm/shortener/internal/repo"
)

/*
Сервер должен быть доступен по адресу: http://localhost:8080.
*/
func RouterStart(h *DecoratedHandler) error {

	r := mux.NewRouter()
	r.Host(h.ServerAddress)
	r.HandleFunc("/", h.RootHandler).Methods(http.MethodPost)
	r.HandleFunc("/{URL}", h.RootURLHandler).Methods(http.MethodGet)

	r.HandleFunc("/api/shorten", h.APIShortenHandler).Methods(http.MethodPost)
	r.HandleFunc("/", h.BadRequestHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         "localhost:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

	return nil
}

type DecoratedHandler struct {
	ServerAddress string
	BaseURL       string
	Storage       *repo.Storage
}
