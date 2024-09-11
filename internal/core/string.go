package core

import "strings"

func GetSliceFromString(s string) []string {
	var replacer = strings.NewReplacer(" ", "", "\"", "", "[", "", "]", "")
	return strings.Split(replacer.Replace(s), ",")
}
