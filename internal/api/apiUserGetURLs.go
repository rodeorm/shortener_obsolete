package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/rodeorm/shortener/internal/api/cookie"
)

/*APIUserGetURLsHandler возвращает пользователю все когда-либо сокращённые им URL в формате JSON*/
func (h DecoratedHandler) APIUserGetURLsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	log.Println("APIUserGetURLsHandler")

	w, userKey := cookie.GetUserIdentity(h.Storage, w, r)
	userID, err := strconv.Atoi(userKey)
	if err != nil {
		fmt.Println("Проблемы с получением пользователя", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	URLHistory, err := h.Storage.SelectUserURLHistory(ctx, userID)
	if err != nil {
		fmt.Println("Проблемы с получением истории пользователя", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	bodyBytes, err := json.Marshal(URLHistory)
	if err != nil {
		fmt.Println("Проблемы при маршалинге истории урл", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bodyBytes))
}
