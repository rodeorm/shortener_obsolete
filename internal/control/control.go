package control

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	middleware "github.com/rodeorm/shortener/internal/control/middleware"
	repo "github.com/rodeorm/shortener/internal/repo"
)

//RouterStart запускает веб-сервер
func RouterStart(h *DecoratedHandler) error {

	r := mux.NewRouter()
	r.HandleFunc("/", h.RootHandler).Methods(http.MethodPost)
	r.HandleFunc("/{URL}", h.RootURLHandler).Methods(http.MethodGet)

	r.HandleFunc("/api/shorten", h.APIShortenHandler).Methods(http.MethodPost)
	r.HandleFunc("/", h.BadRequestHandler)
	r.Use(middleware.GzipMiddleware)
	srv := &http.Server{
		Handler:      r,
		Addr:         h.ServerAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
	return nil
}

type DecoratedHandler struct {
	ServerAddress string
	BaseURL       string
	Storage       repo.AbstractStorage
}
