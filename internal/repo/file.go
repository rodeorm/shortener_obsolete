package repo

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/rodeorm/shortener/internal/logic"
)

type fileStorage struct {
	filePath string
}

func (s fileStorage) CheckFile() error {
	fileInfo, err := os.Stat(s.filePath)

	if errors.Is(err, os.ErrNotExist) {
		newFile, err := os.Create(s.filePath)
		if err != nil {
			log.Fatal(err)
			return err
		}
		newFile.Close()
		fmt.Println("Создан файл: ", newFile.Name())
		return nil
	}
	fmt.Println("Файл уже есть: ", fileInfo.Name())
	return nil
}

// InsertShortURL принимает оригинальный URL, генерирует для него ключ и сохраняет соответствие оригинального URL и ключа (либо возвращает ранее созданный ключ)
func (s fileStorage) InsertShortURL(URL string) (string, error) {

	if !logic.CheckURLValidity(URL) {
		return "", fmt.Errorf("невалидный URL: %s", URL)
	}
	URL = logic.GetClearURL(URL, "")
	key, isExist := s.getShortlURLFromFile(URL)
	if isExist {
		return key, nil
	}
	key, _ = logic.ReturnShortKey(5)

	f, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	pair := logic.URLPair{Origin: URL, Short: key}
	data, err := json.Marshal(pair)
	if err != nil {
		return "", err
	}
	data = append(data, '\n')
	_, err = f.Write(data)
	return key, err
}

//getShortlURLFromFile возвращает из файла сокращенный URL по оригинальному URL
func (s fileStorage) getShortlURLFromFile(URL string) (string, bool) {

	file, err := os.Open(s.filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var up logic.URLPair
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		json.Unmarshal(scanner.Bytes(), &up)
		if up.Origin == URL {
			return up.Short, true
		}
	}

	return "", false
}

// SelectOriginalURL принимает на вход короткий URL (относительный, без имени домена), извлекает из него ключ и возвращает оригинальный URL из хранилища
func (s fileStorage) SelectOriginalURL(shortURL string) (string, bool, error) {

	file, err := os.Open(s.filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var up logic.URLPair
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		json.Unmarshal(scanner.Bytes(), &up)
		if up.Short == shortURL {
			return up.Origin, true, nil
		}
	}

	return "", false, err

}
