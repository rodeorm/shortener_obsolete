package logic

import (
	"net/url"
	"strings"
)

//GetClearURL делает URL строчным, убирает наименование домена
func GetClearURL(s string, d string) string {
	s = strings.ToLower(s)
	return strings.Replace(s, d, "", 1)
}

// CheckURLValidity проверяет URL на корректность
func CheckURLValidity(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}
