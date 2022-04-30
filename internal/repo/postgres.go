package repo

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

/*
type postgresStorage struct {
	originalToShort  map[string]string
	shortToOriginal  map[string]string
	users            map[int]*User
	userURLPairs     map[int]*[]UserURLPair
	DB               *sql.DB // Драйвер подключения к СУБД
	DBName           string  // Имя БД из конфиг.файла
	ConnectionString string  // Строка подключения из конфиг.файла
}

//InsertShortURL принимает оригинальный URL, генерирует для него ключ и сохраняет соответствие оригинального URL и ключа (либо возвращает ранее созданный ключ)
func (s postgresStorage) InsertURL(URL, baseURL, userKey string) (string, error) {
	if !logic.CheckURLValidity(URL) {
		return "", fmt.Errorf("невалидный URL: %s", URL)
	}
	key, isExist := s.originalToShort[URL]
	if isExist {
		s.insertUserURLPair(userKey, baseURL+"/"+key, URL)
		return key, nil
	}
	key, _ = logic.ReturnShortKey(5)

	s.originalToShort[URL] = key
	s.shortToOriginal[key] = URL

	s.insertUserURLPair(userKey, baseURL+"/"+key, URL)

	return key, nil
}

//SelectOriginalURL принимает на вход короткий URL (относительный, без имени домена), извлекает из него ключ и возвращает оригинальный URL из хранилища
func (s postgresStorage) SelectOriginalURL(shortURL string) (string, bool, error) {
	originalURL, isExist := s.shortToOriginal[shortURL]
	return originalURL, isExist, nil
}

//InsertUser сохраняет нового пользователя или возвращает уже имеющегося в наличии
func (s postgresStorage) InsertUser(Key int) (*User, error) {
	if Key == 0 {
		user := &User{Key: s.getNextFreeKey()}
		s.users[user.Key] = user
		return user, nil
	}
	user, isExist := s.users[Key]
	if !isExist {
		user = &User{Key: Key}
		s.users[Key] = user
	}
	return user, nil
}

//InsertUserURLPair cохраняет информацию о том, что пользователь сокращал URL, если такой информации ранее не было
func (s postgresStorage) insertUserURLPair(userKey, shorten, origin string) error {
	userID, err := strconv.Atoi(userKey)
	if err != nil {
		return fmt.Errorf("ошибка обработки идентификатора пользователя: %s", err)
	}

	URLPair := &UserURLPair{UserKey: userID, Short: shorten, Origin: origin}

	userURLPairs, isExist := s.userURLPairs[URLPair.UserKey]
	if !isExist {
		userURLPair := *URLPair
		new := make([]UserURLPair, 0, 10)
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

func (s postgresStorage) SelectUserByKey(Key int) (*User, error) {
	user, isExist := s.users[Key]
	if !isExist {
		return nil, fmt.Errorf("нет пользователя с ключом: %d", Key)
	}
	return user, nil
}

//SelectUserURL возвращает перечень соответствий между оригинальным и коротким адресом для конкретного пользователя
func (s postgresStorage) SelectUserURLHistory(Key int) (*[]UserURLPair, error) {
	if s.userURLPairs[Key] == nil {
		return nil, fmt.Errorf("нет истории")
	}
	return s.userURLPairs[Key], nil
}

//getNextFreeKey возвращает ближайший свободный идентификатор пользователя
func (s postgresStorage) getNextFreeKey() int {
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

//CloseConnect закрывает соединение
func (s postgresStorage) CloseConnect() error {
	var err error
	fmt.Println("Закрыто соединение с сервером баз данных")
	err = s.DB.Close()
	if err != nil {
		fmt.Println("Ошибка при попытке закрыть соединение с базой данных: ", err.Error())
	}
	return err
}
*/

//ConnectToDatabase соединяет непосредственно с экземпляром СУБД
func ConnectToDatabase(connectionString string) error {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Println("Не подключается к СУБД: ", err.Error())
		return err
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("Не пингуется СУБД: ", err.Error())
		return err
	}
	fmt.Println("Соединение с СУБД установлено", connectionString)
	return err
}
