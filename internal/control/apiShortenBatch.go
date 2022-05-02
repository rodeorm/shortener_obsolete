package control

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	repo "github.com/rodeorm/shortener/internal/repo"
)

func (h DecoratedHandler) APIShortenBatch(w http.ResponseWriter, r *http.Request) {
	w, userKey := h.GetUserIdentity(w, r)
	var urlReq []repo.UrlWithCorrelationRequest
	var urlRes []repo.UrlWithCorrelationResponse
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &urlReq)

	if err != nil {
		fmt.Println("Ошибка", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, value := range urlReq {
		shortURLKey, err, _ := h.Storage.InsertURL(value.Origin, h.BaseURL, userKey)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		urlResPart := repo.UrlWithCorrelationResponse{CorID: value.CorID, Short: h.BaseURL + "/" + shortURLKey}
		urlRes = append(urlRes, urlResPart)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	bodyBytes, err = json.Marshal(urlRes)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, string(bodyBytes))
}
