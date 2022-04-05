package control

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rodeorm/shortener/internal/logic"
)

// APIShortenHandler принимает в теле запроса JSON-объект {"url":"<some_url>"} и возвращает в ответ объект {"result":"<shorten_url>"}.
func (h DecoratedHandler) APIShortenHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	url := logic.URL{}
	err := json.Unmarshal(bodyBytes, &url)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("Тело POST запроса. Оригинальный URL:", url.Value)
	shortURLKey, _ := h.Storage.InsertShortURL(url.Value)
	w.Header().Set("Content-Type", "JSON")
	shortURL := logic.ShortenURL{}
	shortURL.Value = h.DomainName + "/" + shortURLKey
	w.WriteHeader(http.StatusCreated)
	bodyBytes, err = json.Marshal(shortURL)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, string(bodyBytes))
}
