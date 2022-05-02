package cookie

import (
	"fmt"
	"net/http"

	"github.com/rodeorm/shortener/internal/logic"
)


func GetUserKeyFromCoockie(r *http.Request) (string, error) {
	tokenCookie, err := r.Cookie("token")
	if err != nil {
		return "", err
	}
	if tokenCookie.Value == "" {
		return "", fmt.Errorf("не найдено актуальных cookie")
	}
	userKey, err := logic.Decrypt(tokenCookie.Value)
	if err != nil {
		return "", err
	}
	return userKey, nil
}

func PutUserKeyToCookie(Key string) *http.Cookie {
	val, _ := logic.Encrypt(Key)

	cookie := &http.Cookie{
		Name:   "token",
		Value:  val,
		MaxAge: 300,
	}
	return cookie
}
