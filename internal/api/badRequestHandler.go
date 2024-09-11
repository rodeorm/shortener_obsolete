package api

import (
	"log"
	"net/http"
)

func (h DecoratedHandler) BadRequestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("BadRequestHandler")
	w.WriteHeader(http.StatusBadRequest)
}
