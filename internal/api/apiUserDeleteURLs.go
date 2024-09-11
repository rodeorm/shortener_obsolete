package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

/*
APIUserDeleteURLsHandler принимает список идентификаторов сокращённых URL для удаления в формате: [ "a", "b", "c", "d", ...].
В случае успешного приёма запроса хендлер должен возвращать HTTP-статус 202 Accepted.
*/
func (h DecoratedHandler) APIUserDeleteURLsHandler(w http.ResponseWriter, r *http.Request) {
	w, userKey := h.GetUserIdentity(w, r)

	ctx := context.TODO()

	_, err := strconv.Atoi(userKey)
	if err != nil {
		log.Println("APIUserDeleteURLsHandler", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("APIUserDeleteURLsHandler", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	go h.Storage.DeleteURLs(ctx, string(bodyBytes), userKey)
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, string(bodyBytes))
}
