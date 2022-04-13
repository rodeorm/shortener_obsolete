package main

import "flag"

var (
	a *string
	b *string
	f *string
)

func init() {
	//флаг -a, отвечающий за адрес запуска HTTP-сервера (переменная SERVER_ADDRESS)
	a = flag.String("a", "", "SERVER_ADDRESS")
	//флаг -b, отвечающий за базовый адрес результирующего сокращённого URL (переменная BASE_URL)
	b = flag.String("b", "", "BASE_URL")
	//флаг -f, отвечающий за путь до файла с сокращёнными URL (переменная FILE_STORAGE_PATH)
	f = flag.String("f", "", "FILE_STORAGE_PATH")

}
