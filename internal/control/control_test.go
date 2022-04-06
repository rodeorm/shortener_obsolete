package control

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/rodeorm/shortener/internal/repo"
)

func TestDecoratedHandler_RootURLHandler(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
	}
	tests := []struct {
		name    string
		handler DecoratedHandler
		method  string
		want    want
		request string
		body    []byte
	}{
		{
			//Эндпоинт GET /{id} принимает в качестве URL-параметра идентификатор сокращённого URL и возвращает ответ с кодом 307 и оригинальным URL в HTTP-заголовке Location.
			name:    "Проверка обработки корректных GET запросов (отсутствуют данные короткого url)",
			handler: DecoratedHandler{DomainName: "http://localhost:8080", Storage: repo.NewStorage()},
			method:  "GET",
			request: "http://localhost:8080/10",
			want:    want{statusCode: 400},
		},
		{
			//Эндпоинт POST / принимает в теле запроса строку URL для сокращения и возвращает ответ с кодом 201 и сокращённым URL в виде текстовой строки в теле.
			name:    "Проверка обработки некорректных POST запросов",
			handler: DecoratedHandler{DomainName: "http://localhost:8080", Storage: repo.NewStorage()},
			method:  "POST",
			request: "http://localhost:8080",
			want:    want{statusCode: 400},
		},
		{
			//Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400 (любые кроме GET и POST)
			name:    "Проверка обработки некорректных запросов: PUT",
			handler: DecoratedHandler{DomainName: "http://localhost:8080", Storage: repo.NewStorage()},
			method:  "PUT",
			request: "http://localhost:8080",
			want:    want{statusCode: 400},
		},
		{
			//Нужно учесть некорректные запросы и возвращать для них ответ с кодом 400 (любые кроме GET и POST)
			name:    "Проверка обработки некорректных запросов: DELETE",
			handler: DecoratedHandler{DomainName: "http://localhost:8080", Storage: repo.NewStorage()},
			method:  "DELETE",
			request: "http://localhost:8080",
			want:    want{statusCode: 400},
		},
		{
			//Нужно принимать и возвращать JSON
			name:    "Проверка обработки некорректных запросов: DELETE",
			handler: DecoratedHandler{DomainName: "http://localhost:8080", Storage: repo.NewStorage()},
			method:  "POST (JSON)",
			request: "http://localhost:8080/api/shorten",
			want:    want{statusCode: 400, contentType: "json"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var request *http.Request
			switch tt.method {
			case "POST":
				request = httptest.NewRequest(http.MethodPost, tt.request, nil)
			case "GET":
				request = httptest.NewRequest(http.MethodGet, tt.request, nil)
			case "PUT":
				request = httptest.NewRequest(http.MethodPut, tt.request, nil)
			case "DELETE":
				request = httptest.NewRequest(http.MethodDelete, tt.request, nil)
			}
			w := httptest.NewRecorder()
			h := http.HandlerFunc(tt.handler.RootHandler)
			h.ServeHTTP(w, request)
			result := w.Result()
			err := result.Body.Close()
			require.NoError(t, err)
			assert.Equal(t, tt.want.statusCode, result.StatusCode)

		})
	}
}
