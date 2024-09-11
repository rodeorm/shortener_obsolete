package api

import (
	"net/http"
)

func (h DecoratedHandler) BadRequestHandler(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusBadRequest)
}
