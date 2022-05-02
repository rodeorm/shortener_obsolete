package control

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (h DecoratedHandler) RootHandler(w http.ResponseWriter, r *http.Request) {

	w, userKey := h.GetUserIdentity(w, r)

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	shortURLKey, err, isDuplicated := h.Storage.InsertURL(bodyString, h.BaseURL, userKey)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if isDuplicated {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	fmt.Fprintf(w, h.BaseURL+"/"+shortURLKey)
}
