package repo

import (
	"fmt"

	"github.com/rodeorm/shortener/internal/logic"
)

type Storage struct {
	originalToShort map[string]string
	shortToOriginal map[string]string
}

func NewStorage() *Storage {
	ots := make(map[string]string)
	sto := make(map[string]string)
	storage := Storage{originalToShort: ots, shortToOriginal: sto}
	return &storage
}

// ReturnShortURL принимает оригинальный URL, генерирует для него ключ и сохраняет соответствие оригинального URL и ключа (либо возвращает ранее созданный ключ)
func (s Storage) InsertShortURL(URL string) (string, error) {
	clearURL := logic.GetClearURL(URL)
	fmt.Println("Очищенный оригинальный URL", clearURL)
	value, isExist := s.originalToShort[clearURL]
	if isExist {
		return value, nil
	}
	key, _ := logic.ReturnShortKey(5)

	s.originalToShort[clearURL] = key
	s.shortToOriginal[key] = clearURL
	return key, nil

}

// ReturnOriginalURL принимает на вход короткий URL (относительный, без имени домена), извлекает из него ключ и возвращает оригинальный URL из хранилища
func (s Storage) SelectOriginalURL(shortURL string) (string, bool, error) {
	fmt.Println("Короткий ключ для поиска URL: ", shortURL)
	fmt.Println(s.shortToOriginal)
	value, isExist := s.shortToOriginal[shortURL]
	return value, isExist, nil
}
