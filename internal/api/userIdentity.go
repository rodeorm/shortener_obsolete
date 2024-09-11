package api

import (
	"context"
	"fmt"
	"log"
	"net/http"

	cookie "github.com/rodeorm/shortener/internal/api/cookie"
)

// GetUserIdentity определяет по кукам какой пользователь авторизовался, если куки некорректные, то создает новые
func (h DecoratedHandler) GetUserIdentity(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, string) {
	userKey, err := cookie.GetUserKeyFromCoockie(r)
	ctx := context.TODO()

	user, err := h.Storage.InsertUser(ctx, 0)
	if err != nil {
		log.Println("GetUserIdentity", err)
	}
	userKey = fmt.Sprint(user.Key)
	http.SetCookie(w, cookie.PutUserKeyToCookie(userKey))
	return w, userKey
}
