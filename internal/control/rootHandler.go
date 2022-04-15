package control

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rodeorm/shortener/internal/logic"
)

func (h DecoratedHandler) RootHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	if !logic.CheckURLValidity(bodyString) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	shortURLKey, _ := h.Storage.InsertShortURL(bodyString)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, h.BaseURL+"/"+shortURLKey)
}
