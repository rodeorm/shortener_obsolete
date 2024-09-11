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

func (h DecoratedHandler) APIShortenBatch(w http.ResponseWriter, r *http.Request) {
	var urlReq []core.URLWithCorrelationRequest
	var urlRes []core.URLWithCorrelationResponse

	w, userKey := h.GetUserIdentity(w, r)
	bodyBytes, _ := io.ReadAll(r.Body)
	err := json.Unmarshal(bodyBytes, &urlReq)
	ctx := context.TODO()

	if err != nil {
		log.Println("APIShortenBatch", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, value := range urlReq {
		shortURLKey, _, err := h.Storage.InsertURL(ctx, value.Origin, h.BaseURL, userKey)
		if err != nil {
			log.Println("APIShortenBatch", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		urlResPart := core.URLWithCorrelationResponse{CorID: value.CorID, Short: h.BaseURL + "/" + shortURLKey}
		urlRes = append(urlRes, urlResPart)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	bodyBytes, err = json.Marshal(urlRes)
	if err != nil {
		log.Println("APIShortenBatch", err)
	}
	fmt.Fprint(w, string(bodyBytes))
}
