package logic

import (
	"strings"
)

type URL struct {
	Key string `json:"url,omitempty"`
}

type ShortenURL struct {
	Key string `json:"result,omitempty"`
}

type URLPair struct {
	Origin string `json:"origin,omitempty"`
	Short  string `json:"short,omitempty"`
}

//GetClearURL делает URL строчным, убирает наименование домена
func GetClearURL(s string, d string) string {
	s = strings.ToLower(s)
	return strings.Replace(s, d, "", 1)
}
