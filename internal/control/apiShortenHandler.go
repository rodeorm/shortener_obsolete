package control

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	repo "github.com/rodeorm/shortener/internal/repo"
)

// APIShortenHandler принимает в теле запроса JSON-объект {"url":"<some_url>"} и возвращает в ответ объект {"result":"<shorten_url>"}.
func (h DecoratedHandler) APIShortenHandler(w http.ResponseWriter, r *http.Request) {
	url := repo.URL{}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &url)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURLKey, err := h.Storage.InsertShortURL(url.Key)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	shortURL := repo.ShortenURL{}
	shortURL.Key = h.BaseURL + "/" + shortURLKey
	w.WriteHeader(http.StatusCreated)
	bodyBytes, err = json.Marshal(shortURL)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, string(bodyBytes))
}
