package middleware

import (
	"log"
	"net/http"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.Body, r.Header, r.PostForm, r.MultipartForm)
		next.ServeHTTP(w, r)
	})
}
