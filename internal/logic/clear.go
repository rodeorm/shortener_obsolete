package logic

import "strings"

//Обработка поступившего URL (пока только делает его строчным)
func GetClearURL(s string) string {
	s = strings.ToLower(s)

	return s
}
