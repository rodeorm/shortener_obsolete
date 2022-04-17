package control

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (h DecoratedHandler) RootHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	shortURLKey, err := h.Storage.InsertShortURL(bodyString)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, h.BaseURL+"/"+shortURLKey)
}
