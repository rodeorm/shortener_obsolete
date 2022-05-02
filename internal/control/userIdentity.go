package control

import (
	"fmt"
	"net/http"

	cookie "github.com/rodeorm/shortener/internal/control/cookie"
)

//GetUserIdentity определяет по кукам какой пользователь авторизовался, если куки некорректные, то создает новые
func (h DecoratedHandler) GetUserIdentity(w http.ResponseWriter, r *http.Request) (http.ResponseWriter, string) {
	userKey, err := cookie.GetUserKeyFromCoockie(r)

	if err != nil {

		user, _ := h.Storage.InsertUser(0)
		userKey = fmt.Sprint(user.Key)
		http.SetCookie(w, cookie.PutUserKeyToCookie(userKey))
		return w, userKey
	}

	http.SetCookie(w, cookie.PutUserKeyToCookie(userKey))
	return w, userKey
}
