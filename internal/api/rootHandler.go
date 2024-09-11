package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/rodeorm/shortener/internal/api/cookie"
)

// RootHandler POST принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
func (h DecoratedHandler) RootHandler(w http.ResponseWriter, r *http.Request) {
	w, userKey := cookie.GetUserIdentity(h.Storage, w, r)

	ctx := context.TODO()
	bodyBytes, _ := io.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	shortURLKey, isDuplicated, err := h.Storage.InsertURL(ctx, bodyString, h.BaseURL, userKey)
	if err != nil {
		log.Println("RootHandler", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if isDuplicated {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, h.BaseURL+"/"+shortURLKey)
}
