package logic

import "strings"

//clearURL очищает URL от http, https и т.п.
func GetClearURL(s string) string {
	//	s = strings.ReplaceAll(s, "https://", "")
	//	s = strings.ReplaceAll(s, "http://", "")
	s = strings.ToLower(s)

	return s
}

//GetClearKey извлекает ключ из короткого url
func GetClearKey(s string) string {
	s = strings.ReplaceAll(s, "https://", "")
	s = strings.ReplaceAll(s, "http://", "")
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "/", "")

	return s
}
