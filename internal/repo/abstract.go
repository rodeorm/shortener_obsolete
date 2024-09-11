package repo

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type AbstractStorage interface {
	// InsertURL сохраняет соответствие между оригинальным и коротким адресом
	InsertURL(URL, baseURL, userKey string) (string, bool, error)

	// SelectOriginalURL возвращает оригинальный адрес на основании короткого; признак, что url ранее уже сокращался; признак, что url удален
	SelectOriginalURL(shortURL string) (string, bool, bool, error)

	//InsertUser сохраняет нового пользователя или возвращает уже имеющегося в наличии
	InsertUser(Key int) (*User, error)

	// SelectUserURLHistory возвращает перечень соответствий между оригинальным и коротким адресом для конкретного пользователя
	SelectUserURLHistory(Key int) (*[]UserURLPair, error)

	// Массово помечает URL как удаленные. Успешно удалить URL может только пользователь, его создавший.
	DeleteURLs(URL, userKey string) (bool, error)

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

// InitMemoryStorage создает хранилище данных в оперативной памяти
func InitMemoryStorage() *memoryStorage {
	ots := make(map[string]string)
	sto := make(map[string]string)
	usr := make(map[int]*User)
	usrURL := make(map[int]*[]UserURLPair)
	storage := memoryStorage{originalToShort: ots, shortToOriginal: sto, users: usr, userURLPairs: usrURL}
	return &storage
}

// InitFileStorage создает хранилище данных на файловой системе
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

// InitPostgresStorage создает хранилище данных в БД на экземпляре Postgres
func InitPostgresStorage(connectionString string) (*postgresStorage, error) {
	db, err := sqlx.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	ctx := context.TODO()
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	delQueue := make(chan string)
	storage := postgresStorage{DB: db, ConnectionString: connectionString, deleteQueue: delQueue}
	storage.createTables(ctx)

	return &storage, nil
}
