package control

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (h DecoratedHandler) RootHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	if bodyString == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("Тело POST запроса. Оригинальный URL: ", bodyString)
	shortURLKey, _ := h.Storage.InsertShortURL(bodyString)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, h.DomainName+"/"+shortURLKey)
}
