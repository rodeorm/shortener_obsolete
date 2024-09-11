package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

/*
RootURLHandler GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400
При запросе удалённого URL с помощью хендлера GET /{id} нужно вернуть статус 410 Gone
*/
func (h DecoratedHandler) RootURLHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("RootURLHandler")
	currentID := mux.Vars(r)["URL"]
	ctx := context.TODO()

	originalURL, isExist, isDeleted, _ := h.Storage.SelectOriginalURL(ctx, currentID)
	if isDeleted {
		w.WriteHeader(http.StatusGone)
		return
	}
	if isExist {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Println("Оригинальный url, на который будет редирект: ", originalURL)
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
		fmt.Fprintf(w, "%s", originalURL)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", "Для данного URL не найден оригинальный URL")
	}
}
