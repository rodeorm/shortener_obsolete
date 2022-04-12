package repo

import (
	"fmt"

	"github.com/rodeorm/shortener/internal/logic"
)

type fileStorage struct {
	filePath        string
	originalToShort map[string]string
	shortToOriginal map[string]string
}

// InsertShortURL принимает оригинальный URL, генерирует для него ключ и сохраняет соответствие оригинального URL и ключа (либо возвращает ранее созданный ключ)
func (s fileStorage) InsertShortURL(URL string) (string, error) {

	key, isExist := s.getShortlURLFromFile(URL)
	if isExist {
		return key, nil
	}
	key, _ = logic.ReturnShortKey(5)

	return key, nil

}

// SelectOriginalURL принимает на вход короткий URL (относительный, без имени домена), извлекает из него ключ и возвращает оригинальный URL из хранилища
func (s fileStorage) SelectOriginalURL(shortURL string) (string, bool, error) {
	fmt.Println("Короткий ключ для поиска URL: ", shortURL)
	fmt.Println(s.shortToOriginal)
	value, isExist := s.shortToOriginal[shortURL]
	return value, isExist, nil
}

func (s fileStorage) getOriginalURLFromFile(ShortURL string) (string, bool) {
	return "", false
}

func (s fileStorage) getShortlURLFromFile(originalURL string) (string, bool) {
	return "", false
}

func (s fileStorage) addPairToFile(originalURL string, shortURL string) (string, error) {
	return "", nil
}
