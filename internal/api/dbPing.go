package api

import (
	"fmt"
	"net/http"

	repo "github.com/rodeorm/shortener/internal/repo"
)

func (h DecoratedHandler) PingDBHandler(w http.ResponseWriter, r *http.Request) {
	_, err := repo.InitPostgresStorage(h.DatabaseConnectionString)
	if err != nil {
		fmt.Fprintf(w, "%s", "Успешное соединение с БД")
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", "Ошибка соединения с БД")
	}
}
