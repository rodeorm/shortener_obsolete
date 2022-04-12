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
	if !logic.CheckURLValidity(url.Key) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	url.Key = logic.GetClearURL(url.Key, h.BaseURL)
	shortURLKey, _ := h.Storage.InsertShortURL(url.Key)
	w.Header().Set("Content-Type", "application/json")
	shortURL := logic.ShortenURL{}
	shortURL.Key = h.BaseURL + "/" + shortURLKey
	w.WriteHeader(http.StatusCreated)
	bodyBytes, err = json.Marshal(shortURL)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, string(bodyBytes))
}
