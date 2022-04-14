package control

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

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
