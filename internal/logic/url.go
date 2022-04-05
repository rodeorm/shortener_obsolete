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

//Обработка поступившего URL (пока только делает его строчным)
func GetClearURL(s string) string {
	s = strings.ToLower(s)

	return s
}
