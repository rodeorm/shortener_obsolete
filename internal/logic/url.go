package logic

import (
	"strings"
)

type URL struct {
	Value string `json:"url,omitempty"`
}

type ShortenURL struct {
	Value string `json:"result,omitempty"`
}

//GetClearURL делает URL строчным, убирает наименование домена
func GetClearURL(s string, d string) string {
	s = strings.ToLower(s)
	return strings.Replace(s, d, "", 1)
}
