package repo

import (
	"fmt"

	"github.com/rodeorm/shortener/internal/logic"
)

type memoryStorage struct {
	originalToShort map[string]string
	shortToOriginal map[string]string
}

// InsertShortURL принимает оригинальный URL, генерирует для него ключ и сохраняет соответствие оригинального URL и ключа (либо возвращает ранее созданный ключ)
func (s memoryStorage) InsertShortURL(URL string) (string, error) {
	if !logic.CheckURLValidity(URL) {
		return "", fmt.Errorf("невалидный URL: %s", URL)
	}
	key, isExist := s.originalToShort[URL]
	if isExist {
		return key, nil
	}
	key, _ = logic.ReturnShortKey(5)

	s.originalToShort[URL] = key
	s.shortToOriginal[key] = URL
	return key, nil
}

// SelectOriginalURL принимает на вход короткий URL (относительный, без имени домена), извлекает из него ключ и возвращает оригинальный URL из хранилища
func (s memoryStorage) SelectOriginalURL(shortURL string) (string, bool, error) {
	originalURL, isExist := s.shortToOriginal[shortURL]
	return originalURL, isExist, nil
}
