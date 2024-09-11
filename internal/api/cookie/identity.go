package cookie

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/rodeorm/shortener/internal/cookie"
	"github.com/rodeorm/shortener/internal/repo"
)

// GetUserIdentity определяет по кукам какой пользователь авторизовался, если куки некорректные, то создает новые и пишет данные о пользователе в хранилище
func GetUserIdentity(storage repo.AbstractStorage, w http.ResponseWriter, r *http.Request) (http.ResponseWriter, string) {
	ctx := context.TODO()
	strKey, err := cookie.GetUserKeyFromCoockie(r)
	if err != nil {
		log.Println("GetUserIdentity. Ошибка при получении куков", err)
	}

	intKey, err := strconv.Atoi(strKey)
	if err != nil {
		intKey = 0
	}

	user, err := storage.InsertUser(ctx, intKey)
	if err != nil {
		log.Println("GetUserIdentity. Ошибка при вставке пользователя", err)
	}

	http.SetCookie(w, cookie.PutUserKeyToCookie(strconv.Itoa(user.Key)))
	return w, strconv.Itoa(user.Key)
}
