package api

import (
	"context"
	"log"
	"net/http"
	"strconv"

	cookie "github.com/rodeorm/shortener/internal/api/cookie"
)

// GetUserIdentity определяет по кукам какой пользователь авторизовался, если куки некорректные, то создает новые
func (h DecoratedHandler) GetUserIdentity(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, string) {
	ctx := context.TODO()
	userKey, err := cookie.GetUserKeyFromCoockie(r)
	key, _ := strconv.Atoi(userKey)
	user, err := h.Storage.InsertUser(ctx, key)
	if err != nil {
		log.Println("GetUserIdentity", err)
	}

	http.SetCookie(w, cookie.PutUserKeyToCookie(strconv.Itoa(user.Key)))
	return w, userKey
}
