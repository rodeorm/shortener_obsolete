package repo

import (
	"context"
	"database/sql"
)

type AbstractStorage interface {
	// Сохраняет соответствие между оригинальным и коротким адресом
	InsertURL(URL, baseURL, userKey string) (string, error)

	// Возвращает оригинальный адрес на основании короткого
	SelectOriginalURL(shortURL string) (string, bool, error)

	//InsertUser сохраняет нового пользователя или возвращает уже имеющегося в наличии
	InsertUser(Key int) (*User, error)

	// Возвращает перечень соответствий между оригинальным и коротким адресом для конкретного пользователя
	SelectUserURLHistory(Key int) (*[]UserURLPair, error)

	// Закрыть соединение (только для СУБД)
	CloseConnection()
}

// NewStorage определяет место для хранения данных
func NewStorage(filePath, dbConnectionString string) AbstractStorage {
	var storage AbstractStorage

	storage, err := InitPostgresStorage(dbConnectionString)
	if err == nil {
		return storage
	}

	if filePath != "" {
		storage, err = InitFileStorage(filePath)
		if err == nil {
			return storage
		}
	}
	storage = InitMemoryStorage()
	return storage
}

//InitMemoryStorage создает хранилище данных в оперативной памяти
func InitMemoryStorage() *memoryStorage {
	ots := make(map[string]string)
	sto := make(map[string]string)
	usr := make(map[int]*User)
	usrURL := make(map[int]*[]UserURLPair)
	storage := memoryStorage{originalToShort: ots, shortToOriginal: sto, users: usr, userURLPairs: usrURL}
	return &storage
}

//InitFileStorage создает хранилище данных на файловой системе
func InitFileStorage(filePath string) (*fileStorage, error) {
	usr := make(map[int]*User)
	usrURL := make(map[int]*[]UserURLPair)
	storage := fileStorage{filePath: filePath, users: usr, userURLPairs: usrURL}
	err := storage.CheckFile(filePath)
	if err != nil {
		return nil, err
	}
	return &storage, nil
}

//InitPostgresStorage создает хранилище данных в БД на экземпляре Postgres
func InitPostgresStorage(connectionString string) (*postgresStorage, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	ctx := context.TODO()
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	storage := postgresStorage{DB: db, ConnectionString: connectionString}
	storage.createTables(ctx)

	return &storage, nil
}

