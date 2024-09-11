package repo

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/rodeorm/shortener/internal/core"
)

type AbstractStorage interface {
	// InsertURL сохраняет соответствие между оригинальным и коротким адресом
	InsertURL(ctx context.Context, URL, baseURL, userKey string) (string, bool, error)

	// SelectOriginalURL возвращает оригинальный адрес на основании короткого; признак, что url ранее уже сокращался; признак, что url удален
	SelectOriginalURL(ctx context.Context, shortURL string) (string, bool, bool, error)

	//InsertUser сохраняет нового пользователя или возвращает уже имеющегося в наличии
	InsertUser(ctx context.Context, Key int) (*core.User, error)

	// SelectUserURLHistory возвращает перечень соответствий между оригинальным и коротким адресом для конкретного пользователя
	SelectUserURLHistory(ctx context.Context, Key int) (*[]core.UserURLPair, error)

	// Массово помечает URL как удаленные. Успешно удалить URL может только пользователь, его создавший.
	DeleteURLs(ctx context.Context, URL, userKey string) (bool, error)

	// Закрыть соединение (только для СУБД)
	CloseConnection()
}

// NewStorage определяет место для хранения данных
func NewStorage(filePath, dbConnectionString string) AbstractStorage {
	var storage AbstractStorage

	storage, err := InitPostgresStorage(dbConnectionString)
	if err == nil {
		log.Println("Данные хранятся в Postgres: ", dbConnectionString)
		return storage
	} else {
		log.Println("Ошибка при попытке подключения к Postgress: ", err, "Строка подключения: ", dbConnectionString)
	}

	if filePath != "" {
		storage, err = InitFileStorage(filePath)
		if err == nil {
			log.Println("Данные хранятся в файле: ", filePath)
			return storage
		}
	}
	log.Println("Данные хранятся в памяти")
	storage = InitMemoryStorage()
	return storage
}

// InitMemoryStorage создает хранилище данных в оперативной памяти
func InitMemoryStorage() *memoryStorage {
	ots := make(map[string]string)
	sto := make(map[string]string)
	usr := make(map[int]*core.User)
	usrURL := make(map[int]*[]core.UserURLPair)
	storage := memoryStorage{originalToShort: ots, shortToOriginal: sto, users: usr, userURLPairs: usrURL}
	return &storage
}

// InitFileStorage создает хранилище данных на файловой системе
func InitFileStorage(filePath string) (*fileStorage, error) {
	usr := make(map[int]*core.User)
	usrURL := make(map[int]*[]core.UserURLPair)
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
	err = storage.createTables(ctx)
	if err != nil {
		return nil, err
	}

	return &storage, nil
}
