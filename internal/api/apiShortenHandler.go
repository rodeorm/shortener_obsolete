package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/rodeorm/shortener/internal/core"
)

// APIShortenHandler принимает в теле запроса JSON-объект {"url":"<some_url>"} и возвращает в ответ объект {"result":"<shorten_url>"}.
func (h DecoratedHandler) APIShortenHandler(w http.ResponseWriter, r *http.Request) {
	url := core.URL{}
	shortURL := core.ShortenURL{}
	ctx := context.TODO()

	w, userKey := h.GetUserIdentity(w, r)
	bodyBytes, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &url)
	if err != nil {
		log.Println("APIShortenHandler", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	shortURLKey, isDuplicated, err := h.Storage.InsertURL(ctx, url.Key, h.BaseURL, userKey)
	if err != nil {
		log.Println("APIShortenHandler", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	shortURL.Key = h.BaseURL + "/" + shortURLKey
	if isDuplicated {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	bodyBytes, err = json.Marshal(shortURL)
	if err != nil {
		log.Println("APIShortenHandler", err)
	}
	fmt.Fprint(w, string(bodyBytes))
}
