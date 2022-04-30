package control

import (
	"fmt"
	"net/http"

	"github.com/rodeorm/shortener/internal/repo"
)

func (h DecoratedHandler) PingDBHandler(w http.ResponseWriter, r *http.Request) {
	err := repo.ConnectToDatabase(h.DatabaseConnectionString)
	if err == nil {
		fmt.Fprintf(w, "%s", "Успешное соединение с БД")
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", "Ошибка соединения с БД")
	}
}
