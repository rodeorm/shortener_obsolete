package middleware

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		bodyBytes, _ := ioutil.ReadAll(r.Body)
		if IsGzip(r.Header) {
			bodyBytes, _ = DecompressGzip(bodyBytes)
		}

		r.Body = ioutil.NopCloser(strings.NewReader(string(bodyBytes)))

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func DecompressGzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("ошибка при декомпрессии данных из gzip: %v", err)
	}
	defer r.Close()

	var b bytes.Buffer
	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("ошибка при декомпрессии данных из gzip: %v", err)
	}

	return b.Bytes(), nil
}

func IsGzip(headers map[string][]string) bool {
	for _, value := range headers["Content-Encoding"] {
		if value == "application/gzip" || value == "gzip" {
			return true
		}
	}
	return false
}
