package control

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rodeorm/shortener/internal/repo"
)

type DecoratedHandler struct {
	DomainName string
	Storage    *repo.Storage
}

/* Сервер должен предоставлять два эндпоинта: POST / и GET /{id}.
Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400
*/
func (h DecoratedHandler) returnURLHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		bodyString := string(bodyBytes)
		fmt.Println("Тело POST запроса. Оригинальный URL: ", bodyString)
		shortURLKey, _ := h.Storage.InsertShortURL(bodyString)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, h.DomainName+"/"+shortURLKey)
	case "GET":
		currentID := r.URL.String()
		originalURL, isExist, _ := h.Storage.SelectOriginalURL(currentID)
		if isExist {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			//originalURL = "http://" + originalURL
			fmt.Println("Оригинальный url, на который будет редирект: ", originalURL)
			w.Header().Set("Location", originalURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
			fmt.Fprintf(w, "%s", originalURL)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "%s", "Для данного URL не найден оригинальный URL")
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
