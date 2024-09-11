package repo

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/rodeorm/shortener/internal/core"
)

type fileStorage struct {
	filePath     string
	users        map[int]*core.User
	userURLPairs map[int]*[]core.UserURLPair
}

func (s fileStorage) CheckFile(filePath string) error {
	fileInfo, err := os.Stat(filePath)

	if errors.Is(err, os.ErrNotExist) {
		newFile, err := os.Create(filePath)
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
func (s fileStorage) InsertURL(ctx context.Context, URL, baseURL, userKey string) (string, bool, error) {

	if !core.CheckURLValidity(URL) {
		return "", false, fmt.Errorf("невалидный URL: %s", URL)
	}
	URL = core.GetClearURL(URL, "")
	key, isExist := s.getShortlURLFromFile(ctx, URL)
	if isExist {
		return key, true, nil
	}
	key, _ = core.ReturnShortKey(5)

	f, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	pair := core.URLPair{Origin: URL, Short: key}
	data, err := json.Marshal(pair)
	if err != nil {
		return "", false, err
	}
	s.insertUserURLPair(ctx, userKey, baseURL+"/"+key, URL)
	data = append(data, '\n')
	_, err = f.Write(data)
	return key, false, err
}

// getShortlURLFromFile возвращает из файла сокращенный URL по оригинальному URL
func (s fileStorage) getShortlURLFromFile(ctx context.Context, URL string) (string, bool) {

	file, err := os.Open(s.filePath)
	if err != nil {
		ctx.Done()
		log.Fatal(err)
	}
	defer file.Close()

	var up core.URLPair
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
func (s fileStorage) SelectOriginalURL(ctx context.Context, shortURL string) (string, bool, bool, error) {

	file, err := os.Open(s.filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var up core.URLPair
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		json.Unmarshal(scanner.Bytes(), &up)
		if up.Short == shortURL {
			return up.Origin, true, false, nil
		}
	}

	return "", false, false, err

}

// InsertUser сохраняет нового пользователя или возвращает уже имеющегося в наличии
func (s fileStorage) InsertUser(ctx context.Context, Key int) (*core.User, error) {
	if Key == 0 {
		user := &core.User{Key: s.getNextFreeKey()}
		s.users[user.Key] = user
		return user, nil
	}
	user, isExist := s.users[Key]
	if !isExist {
		user = &core.User{Key: Key}
		s.users[Key] = user
	}
	return user, nil
}

// InsertUserURLPair cохраняет информацию о том, что пользователь сокращал URL, если такой информации ранее не было
func (s fileStorage) insertUserURLPair(ctx context.Context, userKey, shorten, origin string) error {
	userID, err := strconv.Atoi(userKey)
	if err != nil {
		ctx.Done()
		return fmt.Errorf("ошибка обработки идентификатора пользователя: %s", err)
	}

	URLPair := &core.UserURLPair{UserKey: userID, Short: shorten, Origin: origin}

	userURLPairs, isExist := s.userURLPairs[URLPair.UserKey]
	if !isExist {
		userURLPair := *URLPair
		new := make([]core.UserURLPair, 0, 10)
		new = append(new, userURLPair)
		s.userURLPairs[URLPair.UserKey] = &new
		return nil
	}

	for _, value := range *userURLPairs {
		if value.Origin == URLPair.Origin {
			return nil
		}
	}
	*s.userURLPairs[URLPair.UserKey] = append(*s.userURLPairs[URLPair.UserKey], *URLPair)

	fmt.Println("Хранится историй запросов пользователей на данный момент: ")
	for _, v := range s.userURLPairs {
		fmt.Println(*v)
	}

	return nil
}

func (s fileStorage) SelectUserByKey(Key int) (*core.User, error) {
	user, isExist := s.users[Key]
	if !isExist {
		return nil, fmt.Errorf("нет пользователя с ключом: %d", Key)
	}
	return user, nil
}

// SelectUserURL возвращает перечень соответствий между оригинальным и коротким адресом для конкретного пользователя
func (s fileStorage) SelectUserURLHistory(ctx context.Context, Key int) (*[]core.UserURLPair, error) {
	if s.userURLPairs[Key] == nil {
		return nil, fmt.Errorf("нет истории")
	}
	return s.userURLPairs[Key], nil
}

// getNextFreeKey возвращает ближайший свободный идентификатор пользователя
func (s fileStorage) getNextFreeKey() int {
	var maxNumber int
	for maxNumber = range s.users {
		break
	}
	for n := range s.users {
		if n > maxNumber {
			maxNumber = n
		}
	}
	return maxNumber + 1
}

func (s fileStorage) CloseConnection() {
	fmt.Println("Закрыто")
}

func (s fileStorage) DeleteURLs(ctx context.Context, URL, userKey string) (bool, error) {
	return true, nil
}
